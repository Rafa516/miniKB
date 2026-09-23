# Arquitetura — branch `melhorias`

Este documento explica como o Mini Kanban funciona nesta branch (`melhorias`): a versão
completa do projeto, com login/cadastro, arrastar e soltar, busca, paginação e testes
automatizados. Para a versão base (sem esses extras), ver
[ARQUITETURA-main.md](https://github.com/Rafa516/miniKB/blob/main/ARQUITETURA-main.md), na
branch `main`.

## Visão geral

```text
Navegador (React)  ──HTTP/JSON + cookie──▶  Backend (Go)  ──SQL──▶  SQLite (kanban.db)
   localhost:5173                              localhost:8080
```

A diferença em relação à `main` é a camada de autenticação: toda chamada às rotas de tarefas
carrega um cookie de sessão, e o backend confere esse cookie antes de responder. Ver detalhes
de tudo que foi adicionado nesta branch em [MELHORIAS.md](MELHORIAS.md) — este documento foca
em explicar a arquitetura em si, não o histórico do que mudou.

## Backend (`backend/`)

```text
backend/
├── cmd/server/main.go                 # rotas HTTP + CORS (com credenciais)
└── internal/
    ├── database/database.go           # conexão SQLite, tabelas, migração, usuário admin
    ├── models/
    │   ├── task.go                    # struct Task
    │   └── user.go                    # struct User (exposta pela API, sem a senha)
    ├── repository/
    │   ├── task_repository.go         # SQL de tarefas (+ paginação)
    │   ├── task_repository_test.go
    │   ├── auth_repository.go         # SQL de usuários e sessões
    ├── handlers/
    │   ├── task_handler.go            # HTTP de tarefas (+ paginação)
    │   ├── task_handler_test.go
    │   ├── auth_handler.go            # HTTP de login/cadastro/sessão
    │   └── auth_handler_test.go
```

### Rotas da API

| Método | Rota          | Autenticação | O que faz                                   |
|--------|---------------|:---:|-----------------------------------------------------|
| POST   | `/register`   | não | Cria um usuário e já retorna logado                 |
| POST   | `/login`      | não | Confere usuário/senha e cria uma sessão              |
| POST   | `/logout`     | não | Encerra a sessão atual                               |
| GET    | `/me`         | não* | Diz quem está logado (401 se não houver sessão)     |
| GET    | `/tasks`      | **sim** | Lista tarefas (aceita `?page=` e `?pageSize=`)   |
| POST   | `/tasks`      | **sim** | Cria uma tarefa                                  |
| PUT    | `/tasks/{id}` | **sim** | Atualiza uma tarefa                              |
| DELETE | `/tasks/{id}` | **sim** | Exclui uma tarefa                                |

\* `/me` não exige sessão para ser *chamada*, mas só responde `200` se houver uma válida —
é assim que o frontend descobre se alguém já está logado.

### Autenticação, passo a passo

```text
1. POST /register ou /login
   └─ AuthHandler.Register / Login
        └─ AuthRepository.CreateUser / GetUserByUsername
        └─ bcrypt confere/gera o hash da senha
        └─ startSession(): gera token aleatório, grava em "sessions",
           manda de volta um cookie HttpOnly (minikb_session)

2. Requisições seguintes (ex.: GET /tasks)
   └─ o navegador reenvia o cookie sozinho (fetch com credentials: "include")
   └─ AuthHandler.RequireAuth (middleware) confere o cookie contra a
      tabela "sessions" antes de deixar chegar no TaskHandler
   └─ sessão inválida/ausente → 401, sem chegar no handler de tarefas
```

Um usuário administrador (`ADMIN_USERNAME`/`ADMIN_PASSWORD`, padrão `admin`/`admin123`) é
recriado a cada início do servidor, para sempre existir um jeito de entrar mesmo sem cadastro
prévio. Além dele, `POST /register` permite criar novas contas livremente (sem confirmação por
e-mail — ver limitações em [MELHORIAS.md](MELHORIAS.md)).

### Banco de dados

Três tabelas (a tabela `tasks` é igual à da `main`; `users` e `sessions` são novas):

```sql
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  expires_at DATETIME NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

`database.go` também tem uma função `migrate()`, que existe só para ajustar bancos SQLite
criados antes da coluna `name` existir em `users` (adiciona a coluna se faltar, sem apagar
dados) — necessário porque `CREATE TABLE IF NOT EXISTS` não altera uma tabela que já existe.

### Paginação

`TaskHandler.GetTasks` continua devolvendo todas as tarefas por padrão. Só passa a paginar
(via `TaskRepository.GetPage`, que usa `LIMIT`/`OFFSET`) se a requisição incluir
`?page=`/`?pageSize=` na URL — o total de registros vai no cabeçalho `X-Total-Count`. O quadro
Kanban do frontend não usa isso hoje (ver o porquê em MELHORIAS.md); é uma capacidade pronta na
API para um uso futuro (ex.: uma tela de lista/administração).

## Frontend (`frontend/src/`)

```text
src/
├── main.jsx                    # ponto de entrada
├── App.jsx                     # componente raiz: sessão + tarefas + busca + drag and drop
├── services/api.js             # todas as chamadas fetch (com cookie de sessão)
├── test/setup.js               # configuração do ambiente de testes (Vitest)
└── components/
    ├── Login.jsx                 # tela de login/cadastro
    ├── TaskForm.jsx               # formulário do topo, cria tarefas
    ├── KanbanColumn.jsx           # uma coluna — também é "zona de soltar" do drag and drop
    └── TaskCard.jsx               # um card — visualização, edição e "alça" de arrastar
```

### Fluxo de dados

```text
App.jsx (estado: user, tasks, search, error, message)
  │
  ├── (sem user) → Login.jsx → onLogin / onRegister
  │
  └── (com user)
        ├── TaskForm            → onTaskCreated(task)
        │
        └── DndContext (@dnd-kit)
              └── KanbanColumn × 3 (uma por status, "zona de soltar")
                    │  filtra `filteredTasks` (já filtradas pela busca) pelo próprio status
                    └── TaskCard × N (cada um "arrastável")
                          → onUpdate(id, task)   edita / muda de coluna
                          → onDelete(id)         remove, com confirmação
```

### Autenticação no frontend

```text
1. App monta → GET /me (com credentials: "include")
     200 → guarda o usuário no estado, segue para o quadro
     401 → mostra <Login />

2. Login.jsx → onLogin/onRegister (props vindas de App.jsx)
     → services/api.js: login()/register() (POST, credentials: "include")
     → sucesso: App guarda o usuário retornado no estado
     → falha: Login.jsx mostra a mensagem de erro vinda do backend
```

Como o cookie de sessão é `HttpOnly`, o JavaScript do frontend nunca lê nem manipula o token
diretamente — só usa `credentials: "include"` para o navegador anexá-lo/recebê-lo sozinho em
cada requisição.

### Drag and drop

`App.jsx` envolve o quadro num `<DndContext>` do `@dnd-kit/core`. Cada `TaskCard` é arrastável
(`useDraggable`, id = id da tarefa) e cada `KanbanColumn` é uma zona de soltar (`useDroppable`,
id = o próprio status). Ao soltar um card sobre uma coluna, `handleDragEnd` (em `App.jsx`)
recebe os dois ids e, se forem de colunas diferentes, chama `handleTaskUpdate` com o novo
status — ou seja, arrastar reaproveita o mesmo caminho de "editar tarefa" que o formulário já
usava, só mudando o campo `status`.

### Busca

Um campo de texto (`search`, estado em `App.jsx`) filtra a lista de tarefas no navegador (sem
chamar a API de novo), usando `useMemo` para não refazer o filtro em todo re-render. O
resultado (`filteredTasks`) é o que as 3 colunas recebem.

### Mensagens de erro e sucesso

Diferente da `main`, aqui os avisos aparecem como um banner por cima do quadro
(`.error-banner` / `.success-banner`) em vez de substituir a tela inteira — e as mensagens de
erro vêm do texto de resposta do próprio backend (via `readErrorMessage` em `api.js`), em vez
de um texto genérico fixo no frontend.

## Testes automatizados

```text
backend/internal/repository/task_repository_test.go   Create, GetAll, Update, Delete, GetPage
backend/internal/handlers/task_handler_test.go         validação, criação, listagem paginada
backend/internal/handlers/auth_handler_test.go         login, cadastro, sessão, middleware
                                                        (todos usando SQLite em memória)

frontend/src/components/TaskForm.test.jsx              envio do formulário de criar tarefa
frontend/src/components/TaskCard.test.jsx               confirmação antes de excluir
frontend/src/components/Login.test.jsx                  login e cadastro (sucesso e erro)
```

Rodar: `go test ./...` (dentro de `backend/`) e `npm test` (dentro de `frontend/`).

## Docker

Os arquivos de Docker desta branch são os mesmos da `main`, com a adição das variáveis
`ADMIN_USERNAME`/`ADMIN_PASSWORD`. Ver [DOCKER.md](DOCKER.md), especialmente a
[seção 10](DOCKER.md#10-login-e-cadastro-de-usuários), sobre login.
