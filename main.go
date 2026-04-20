package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
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
		r.Get("/{id}", getTask)
		r.Post("/", createTask)
		r.Patch("/{id}", updateTask)
		r.Delete("/{id}", deleteTask)
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

func getTask(w http.ResponseWriter, r *http.Request) {
    // 1. Captura o ID da URL usando o Chi
    idStr := chi.URLParam(r, "id")
    
    // 2. Converte para inteiro
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    var t Task
    // 3. Executa a query
    err = dbPool.QueryRow(context.Background(),
        "SELECT id, title, completed, created_at FROM tasks WHERE id = $1", id).
        Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)

    if err != nil {
        // 4. Trata especificamente o caso de não encontrar nada (404)
        if err == pgx.ErrNoRows {
            http.Error(w, "Task not found", http.StatusNotFound)
            return
        }
        // Se for outro erro de banco, devolve 500
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(t)
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

func updateTask(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    // Faz o parse do payload (JSON) para pegar os novos dados
    var req Task
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var t Task
    // Atualiza e já retorna os dados atualizados
    err = dbPool.QueryRow(context.Background(),
        "UPDATE tasks SET title = $1, completed = $2 WHERE id = $3 RETURNING id, title, completed, created_at",
        req.Title, req.Completed, id).
        Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)

    if err != nil {
        if err == pgx.ErrNoRows {
            http.Error(w, "Task not found", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(t)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    result, err := dbPool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if result.RowsAffected() == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
