# Go REST API - Task Manager 🚀

Uma API RESTful performática e minimalista construída em **Go**, projetada para gerenciar tarefas (TODO list) com persistência em **PostgreSQL** e orquestração via **Docker**.

## 📋 Sobre o Projeto

Este projeto serve como uma base sólida para entender a construção de APIs em Go, utilizando as melhores práticas da linguagem, como roteamento leve com Chi e drivers de banco de dados de alta performance.

### ✨ Funcionalidades
- 📝 CRUD Completo de tarefas.
- 🗄️ Persistência robusta com PostgreSQL.
- 🐳 Ambiente totalmente containerizado.
- 🔄 Migrações automáticas de banco de dados na inicialização.
- ⏳ Resiliência com lógica de retry na conexão com o banco.

---

## 🛠 Tecnologias

- **Linguagem:** Go 1.22
- **Roteador:** [Chi Router](https://github.com/go-chi/chi)
- **Banco de Dados:** PostgreSQL 15
- **Driver DB:** [pgx](https://github.com/jackc/pgx)
- **Containerização:** Docker & Docker Compose

---

## 🚀 Como Executar

Certifique-se de ter o **Docker** e o **Docker Compose** instalados em sua máquina.

1. **Clone o repositório:**
   ```bash
   git clone <URL_DO_REPO>
   cd go-rest-api
   ```

2. **Suba os containers:**
   ```bash
   docker-compose up --build
   ```

A API estará disponível em `http://localhost:8080`.

---

## 📡 Endpoints da API

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `GET` | `/api/v1/tasks` | Lista todas as tarefas |
| `GET` | `/api/v1/tasks/{id}` | Obtém uma tarefa específica |
| `POST` | `/api/v1/tasks` | Cria uma nova tarefa |
| `PATCH` | `/api/v1/tasks/{id}` | Atualiza uma tarefa existente |
| `DELETE` | `/api/v1/tasks/{id}` | Remove uma tarefa |

---

## 📖 Documentação Adicional

Para uma explicação detalhada sobre a arquitetura do código, conceitos de Go aplicados e guia de estudo, consulte o arquivo [DOCUMENTATION.md](./DOCUMENTATION.md).

---

## 🔨 Desenvolvimento

Se desejar rodar localmente sem Docker (necessita de um Postgres rodando):

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
go run main.go
```

---
Desenvolvido com ❤️ e ☕ por Dev Leonzera
