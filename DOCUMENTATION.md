# Guia da API REST em Go (Task Manager)

Este projeto é uma API RESTful minimalista construída em Go para gerenciar uma lista de tarefas (TODO list). Ele utiliza PostgreSQL para persistência de dados e Docker para orquestração do ambiente.

## 🛠 Tecnologias Utilizadas

*   **Go (Golang) 1.22**: Linguagem de programação eficiente e tipada.
*   **Chi Router**: Um roteador leve e idiomático para Go.
*   **pgx**: Driver de PostgreSQL de alta performance.
*   **PostgreSQL**: Banco de dados relacional.
*   **Docker & Docker Compose**: Para containerização e gerenciamento de serviços.

---

## 🏗 Estrutura do Projeto

*   `main.go`: Contém toda a lógica da aplicação, roteamento e interação com o banco de dados.
*   `Dockerfile`: Define como construir a imagem Docker da aplicação Go.
*   `docker-compose.yml`: Define os serviços (App e Banco de Dados) e como eles se conectam.
*   `go.mod`: Gerencia as dependências do projeto.

---

## 🚀 Funcionamento Geral

1.  **Inicialização**: O programa começa carregando as variáveis de ambiente (como a URL do banco).
2.  **Conexão com Banco de Dados**: Utiliza um "Pool" de conexões (`pgxpool`). Existe uma lógica de *retry* que tenta conectar 5 vezes antes de desistir, garantindo que a aplicação espere o banco subir no Docker.
3.  **Migração Automática**: Ao iniciar, a API executa um comando SQL para criar a tabela `tasks` caso ela não exista.
4.  **Servidor HTTP**: O roteador Chi define os caminhos (endpoints) e vincula cada um a uma função específica (Handler).

---

## 📖 Funções e Endpoints

### 1. `main()`
Ponto de entrada do programa. Configura a conexão com o banco, cria a tabela, define as rotas e inicia o servidor na porta 8080.

### 2. `listTasks(w, r)` -> `GET /api/v1/tasks`
Busca todas as tarefas no banco de dados ordenadas pela data de criação. Retorna um array JSON.

### 3. `getTask(w, r)` -> `GET /api/v1/tasks/{id}`
Busca uma única tarefa pelo seu ID. Se não encontrar, retorna um erro `404 Not Found`.

### 4. `createTask(w, r)` -> `POST /api/v1/tasks`
Recebe um JSON com `title` e `completed`, insere no banco e retorna a tarefa criada com seu novo `id` e `created_at`.

### 5. `updateTask(w, r)` -> `PATCH /api/v1/tasks/{id}`
Atualiza os campos de uma tarefa existente. Utiliza o comando `UPDATE ... RETURNING` para devolver os dados atualizados em uma única operação.

### 6. `deleteTask(w, r)` -> `DELETE /api/v1/tasks/{id}`
Remove o registro do banco. Retorna um status `204 No Content` se a exclusão for bem-sucedida.

---

## 💡 Para entender APIs em Go, você precisa saber:

Se você está começando agora, foque nestes 4 pilares:

### 1. Structs e Tags JSON
Em Go, usamos `structs` para modelar os dados. As "tags" (como `json:"id"`) dizem ao Go como transformar a struct em JSON e vice-versa.
```go
type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}
```

### 2. Handlers HTTP (`w` e `r`)
Toda função que lida com uma requisição recebe dois argumentos:
*   `w http.ResponseWriter`: Onde você escreve a resposta (status code, corpo JSON).
*   `r *http.Request`: Onde você lê os dados que o usuário enviou (parâmetros de URL, corpo da requisição).

### 3. Contexto (`context`)
O Go usa `context` para gerenciar o ciclo de vida de uma requisição. Se um usuário cancela a chamada, o `context` avisa o banco de dados para parar a query e economizar recursos.

### 4. Tratamento de Erros
Em Go, erros são valores. Você verá muito o padrão:
```go
if err != nil {
    // Trata o erro aqui
}
```
Isso torna o código muito explícito sobre o que pode falhar.

---

## 📝 Primeiros Passos para Estudo
1.  **Tour of Go**: Faça o tutorial oficial no site [go.dev](https://go.dev/tour/).
2.  **Standard Library**: Explore o pacote `net/http`. É a base de tudo.
3.  **JSON Marshalling**: Entenda como o pacote `encoding/json` funciona.
4.  **Interfaces**: Entenda como o Go usa interfaces para abstrair comportamentos (como o banco de dados).
