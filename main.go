package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var dbPool *pgxpool.Pool

func main() {
	// 1. Connection string from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}

	// 2. Connect to Database (with retry)
	var err error
	for i := 0; i < 5; i++ {
		dbPool, err = pgxpool.New(context.Background(), dbURL)
		if err == nil {
			err = dbPool.Ping(context.Background())
			if err == nil {
				break
			}
		}
		log.Printf("Connecting to DB... attempt %d", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	log.Println("Successfully connected to database")

	// 3. Simple Migration (Create table if not exists)
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = dbPool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	// 4. Router Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", listTasks)
		r.Post("/", createTask)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := dbPool.Query(context.Background(), "SELECT id, title, completed, created_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := dbPool.QueryRow(context.Background(),
		"INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id, created_at",
		t.Title, t.Completed).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}
