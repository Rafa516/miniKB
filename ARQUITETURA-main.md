# Arquitetura — branch `main`

Este documento explica como o Mini Kanban funciona nesta branch (`main`): a versão "base" do
projeto — sem login, sem drag and drop, sem os outros extras que existem na branch `melhorias`
(ver [ARQUITETURA-melhorias.md](ARQUITETURA-melhorias.md) para essa versão).

## Visão geral

```text
Navegador (React)  ──HTTP/JSON──▶  Backend (Go)  ──SQL──▶  SQLite (kanban.db)
   localhost:5173                    localhost:8080
```

- **Frontend**: React + Vite, fala com o backend via `fetch` (JSON puro, sem biblioteca de
  requisição).
- **Backend**: Go, só com a biblioteca padrão (`net/http`), sem framework web.
- **Banco**: SQLite — um arquivo (`kanban.db`), sem servidor de banco separado.

## O que foi feito nesta branch (além do código original)

O código de aplicação (rotas, componentes, modelo de dados) é o mesmo do desafio original. As
adições desta branch foram:

1. **Containerização com Docker** (Dockerfiles, `docker-compose.yml`, Nginx para servir o
   frontend) — ver [DOCKER.md](DOCKER.md) para o passo a passo completo.
2. **CORS configurável** (`backend/cmd/server/main.go`): a origem liberada para o frontend
   deixou de ser fixa no código e passou a vir da variável de ambiente `ALLOWED_ORIGINS`, para
   permitir acesso de outras máquinas na rede sem precisar recompilar.
3. **Endereço do backend calculado automaticamente** (`frontend/src/services/api.js`): em vez
   de sempre apontar para `localhost:8080`, o frontend descobre o endereço do backend a partir
   de `window.location.hostname` — funciona tanto em `localhost` quanto acessando por outro IP.
4. **Comentários explicando os blocos de código** em todos os arquivos do backend e do
   frontend (este documento é o complemento em texto corrido).

## Backend (`backend/`)

```text
backend/
├── cmd/server/main.go              # ponto de entrada: rotas HTTP + CORS
└── internal/
    ├── database/database.go        # conexão com o SQLite + criação da tabela
    ├── models/task.go              # struct Task (linha do banco / JSON da API)
    ├── repository/task_repository.go  # SQL: GetAll, Create, Update, Delete
    └── handlers/task_handler.go    # HTTP: decodifica JSON, valida, chama o repository
```

O fluxo de uma requisição segue sempre essa ordem:

```text
main.go (rota)  →  handler (valida, decodifica JSON)  →  repository (SQL)  →  SQLite
```

### Rotas da API

| Método | Rota          | Handler               | O que faz                         |
|--------|---------------|------------------------|-----------------------------------|
| GET    | `/tasks`      | `GetTasks`             | Lista todas as tarefas            |
| POST   | `/tasks`      | `CreateTask`           | Cria uma tarefa                   |
| PUT    | `/tasks/{id}` | `UpdateTask`           | Atualiza uma tarefa                |
| DELETE | `/tasks/{id}` | `DeleteTask`           | Exclui uma tarefa                  |

Não existe roteador externo (tipo `gorilla/mux` ou `chi`) — as rotas são registradas direto com
`http.HandleFunc`, e o método HTTP é decidido dentro de cada handler com um `switch`.

### CORS

Toda a API passa pela função `enableCORS` (em `main.go`) antes de chegar nos handlers. Ela:
- Lê a lista de origens permitidas da variável `ALLOWED_ORIGINS` (padrão:
  `http://localhost:5173`).
- Só libera a resposta se a origem da requisição bater com alguma da lista.
- Responde sozinha às requisições `OPTIONS` (o "preflight" que o navegador manda antes de
  requisições não triviais).

### Banco de dados

Uma única tabela, criada automaticamente na primeira execução:

```sql
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL,       -- "todo" | "doing" | "done"
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

Não existe um "servidor de banco" separado — o SQLite é uma biblioteca embutida no próprio
processo Go, lendo e escrevendo direto no arquivo `kanban.db`.

## Frontend (`frontend/src/`)

```text
src/
├── main.jsx                    # ponto de entrada, monta <App /> na página
├── App.jsx                     # componente raiz: estado das tarefas + layout
├── services/api.js             # todas as chamadas fetch para o backend
└── components/
    ├── TaskForm.jsx             # formulário do topo, cria tarefas
    ├── KanbanColumn.jsx         # uma coluna (TODO/DOING/DONE)
    └── TaskCard.jsx             # um card — visualização e edição
```

### Fluxo de dados

`App.jsx` é o único lugar que guarda a lista de tarefas (`useState`). Os componentes filhos não
têm estado de dados próprio — só recebem `tasks` e funções (`onTaskCreated`, `onUpdate`,
`onDelete`) via props e chamam essas funções quando o usuário interage:

```text
App.jsx (estado: tasks)
  │
  ├── TaskForm      → onTaskCreated(task)         cria e insere no início da lista
  │
  └── KanbanColumn × 3 (uma por status)
        │  filtra `tasks` pelo próprio status
        └── TaskCard × N
              → onUpdate(id, task)   edita e mescla na lista
              → onDelete(id)         remove da lista
```

Cada ação (`handleTaskCreated`, `handleTaskUpdate`, `handleTaskDelete`, em `App.jsx`) segue o
mesmo padrão: chama a função correspondente de `services/api.js`, e se der certo, atualiza o
estado local `tasks` diretamente (sem recarregar tudo do backend de novo).

### Um comportamento a saber (bug conhecido)

Em `App.jsx`, os `if (loading)` / `if (error)` no topo do componente fazem parte de **todo**
render, não só do carregamento inicial. Ou seja: se qualquer ação (criar, editar ou excluir)
falhar, a tela inteira passa a mostrar só a mensagem de erro, escondendo o quadro, até a página
ser recarregada. Isso está corrigido na branch `melhorias` — ver
[MELHORIAS.md](https://github.com/Rafa516/miniKB/blob/melhorias/MELHORIAS.md) nessa branch.

### Estilo (`App.css`)

Um único arquivo CSS, sem biblioteca de componentes — classes como `.task-form`,
`.kanban-column`, `.task-card` correspondem diretamente às classes usadas no JSX.

## Docker

Ver [DOCKER.md](DOCKER.md) para o guia completo (o que é Docker, como instalar, como rodar,
perguntas frequentes e os detalhes técnicos de cada decisão tomada).
