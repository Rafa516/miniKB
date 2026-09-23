# Deploy no Railway — guia passo a passo (branch `deploy`)

Esta branch deixa o Mini Kanban pronto para publicar no [Railway](https://railway.app/), com
banco de dados PostgreSQL (em vez do SQLite usado nas branches `main` e `melhorias`). Nada
aqui publica nada sozinho — os passos abaixo são para você seguir manualmente, já que criar
conta e configurar serviços em um provedor externo só pode ser feito por você.

## O que muda nesta branch em relação à `melhorias`

| Item | `melhorias` | `deploy` |
|---|---|---|
| Banco de dados | SQLite (arquivo `kanban.db`, num volume Docker) | PostgreSQL (via `DATABASE_URL`) |
| Porta do backend | fixa em `8080` | lida de `$PORT` (padrão `8080` se não definida) |
| Porta do frontend (Nginx) | fixa em `80` | lida de `$PORT` (padrão `80` se não definida) |
| Endereço do backend no frontend | calculado automaticamente (`window.location.hostname`) | pode ser fixado no build via `VITE_API_URL` (necessário no Railway, onde cada serviço tem seu próprio domínio) |

Por que essas mudanças eram necessárias — e o que quebraria sem elas — está detalhado nas
seções abaixo. Se quiser rodar localmente antes de publicar, o `docker-compose.yml` desta
branch já sobe um Postgres junto (ver seção "Testando localmente").

## 1. Pré-requisitos

- O código já enviado para o GitHub, nesta branch (`deploy`)
- Uma conta no [Railway](https://railway.app/) (dá para criar com login do GitHub)

O Railway é pago por uso (tem um plano de teste com créditos grátis, mas cartão de crédito
pode ser pedido dependendo da promoção vigente) — confira as condições atuais no site deles
antes de prosseguir.

## 2. Criar o projeto e os serviços

Como o repositório é um **monorepo** (backend e frontend na mesma pasta), o Railway precisa de
**dois serviços dentro do mesmo projeto**, um para cada pasta:

1. No painel do Railway, clique em **New Project → Deploy from GitHub repo** e selecione
   `Rafa516/miniKB`.
2. O Railway vai criar um serviço a partir da raiz do repositório. Configure esse primeiro
   serviço como o **backend**:
   - **Settings → Root Directory**: `backend`
   - **Settings → Branch**: `deploy`
   - O Railway detecta o `backend/Dockerfile` e o `backend/railway.json` automaticamente.
3. Adicione um **segundo serviço** (botão **+ New → GitHub Repo**, mesmo repositório) para o
   **frontend**:
   - **Settings → Root Directory**: `frontend`
   - **Settings → Branch**: `deploy`
   - Detecta `frontend/Dockerfile` e `frontend/railway.json`.

## 3. Adicionar o PostgreSQL

1. No mesmo projeto, **+ New → Database → Add PostgreSQL**.
2. O Railway cria o banco e já disponibiliza a variável `DATABASE_URL` — mas ela precisa ser
   **conectada ao serviço do backend**: no serviço do backend, vá em **Variables → Add
   Reference** e selecione a variável `DATABASE_URL` do serviço do Postgres.

Não é preciso rodar nenhum script de criação de tabela manualmente — o backend cria (e migra)
as tabelas sozinho ao iniciar (ver `backend/internal/database/database.go`).

## 4. Variáveis de ambiente

No serviço do **backend** (Settings → Variables):

| Variável | Valor |
|---|---|
| `DATABASE_URL` | referência ao Postgres (passo 3) |
| `ALLOWED_ORIGINS` | a URL pública do serviço do frontend (ver passo 5), ex.: `https://minikb-frontend-production.up.railway.app` |
| `ADMIN_USERNAME` | opcional, padrão `admin` |
| `ADMIN_PASSWORD` | **troque isso** — não deixe a senha padrão `admin123` em algo publicado |

`PORT` **não precisa ser definida** — o Railway já injeta essa variável automaticamente, e o
backend já lê ela sozinho (ver `backend/cmd/server/main.go`).

No serviço do **frontend** (Settings → Variables, estas contam como *build-time*, então
precisam existir antes do build rodar):

| Variável | Valor |
|---|---|
| `VITE_API_URL` | a URL pública do serviço do backend, ex.: `https://minikb-backend-production.up.railway.app` |

## 5. Gerar os domínios públicos

Em cada serviço (backend e frontend): **Settings → Networking → Generate Domain**. O Railway
gera uma URL tipo `https://<nome>-production.up.railway.app`.

Depois de gerar os dois domínios, **volte no passo 4** e confira se `ALLOWED_ORIGINS` (no
backend) e `VITE_API_URL` (no frontend) estão apontando um para o outro corretamente — esses
dois valores dependem um do outro, então normalmente só ficam certos depois que ambos os
domínios já existem.

## 6. Redeploy

Depois de ajustar as variáveis, force um redeploy dos dois serviços (o frontend precisa
recompilar para aplicar `VITE_API_URL`, já que é uma variável usada em tempo de build, não de
execução). No painel de cada serviço: **Deployments → ⋮ → Redeploy**.

## 7. Testando

Acesse a URL pública do frontend. Deve aparecer a tela de login — entre com `admin` e a senha
que você configurou em `ADMIN_PASSWORD` (ou crie uma conta pelo "Cadastre-se").

## Testando localmente antes de publicar

Esta branch já vem com um `docker-compose.yml` que sobe um Postgres local, então dá para testar
tudo exatamente como vai rodar no Railway (só que com portas fixas em vez de domínios):

```bash
docker compose up --build
```

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- Postgres: exposto em `localhost:5432` (usuário/senha/banco: `minikb`/`minikb`/`minikb`), só
  para facilitar rodar os testes automatizados fora do Docker

Ver [DOCKER.md](DOCKER.md) para o guia completo do Docker (o que é, como instalar, etc.) — os
comandos são os mesmos, só o banco por trás mudou de SQLite para Postgres.

## Diferenças técnicas, em detalhe

### Por que a porta não podia ficar fixa

O Railway (como a maioria dos serviços de deploy) decide sozinho em qual porta seu serviço vai
escutar, e informa isso via variável de ambiente `PORT` — o container é obrigado a escutar
nela, senão o roteamento de tráfego externo não encontra a aplicação. Como o código original
tinha a porta `8080` fixa em `http.ListenAndServe(":8080", ...)`, isso foi trocado para ler de
`os.Getenv("PORT")`, com `8080` como padrão (mantendo o comportamento local de sempre, sem
precisar dessa variável).

O mesmo vale para o Nginx do frontend: em vez de um `nginx.conf` fixo com `listen 80;`, agora é
um **template** (`nginx.conf.template`) com `listen ${PORT};`, que a própria imagem oficial do
Nginx substitui pelo valor real ao iniciar o container (usando `envsubst`, um recurso já
embutido na imagem — não foi preciso escrever nenhum script novo para isso).

### Por que o endereço do backend precisa ser fixado no build do frontend

Nas branches `main`/`melhorias`, o frontend descobre o endereço do backend sozinho
(`window.location.hostname` + porta 8080), porque as duas partes rodam no mesmo host,
só em portas diferentes. No Railway, cada serviço ganha um **domínio próprio e completo**
(não é o mesmo host em portas diferentes) — então esse truque automático não funciona, e o
endereço do backend precisa ser passado explicitamente como variável de build (`VITE_API_URL`),
via `--build-arg` no `frontend/Dockerfile`.

### Por que o SQLite virou PostgreSQL

O SQLite grava os dados em um arquivo local (`kanban.db`). O Railway (como a maioria dos
serviços de deploy) **não garante disco persistente** para o sistema de arquivos do container
por padrão — a cada novo deploy, o disco é recriado do zero, e os dados se perderiam. O Railway
oferece um addon de PostgreSQL gerenciado (com armazenamento persistente de verdade) pronto
para conectar em qualquer serviço, então essa branch troca o driver e as consultas SQL para
usar Postgres em vez de SQLite:

- Driver: `github.com/jackc/pgx/v5/stdlib` no lugar de `modernc.org/sqlite` (ambos em Go puro,
  sem precisar de CGO).
- Os placeholders das queries SQL mudaram de `?` (SQLite) para `$1, $2, ...` (Postgres).
- `INSERT ... RETURNING id` no lugar de `Result.LastInsertId()` — o driver do Postgres não
  implementa esse método do `database/sql` (o SQLite implementa; o Postgres, historicamente,
  resolve isso com a cláusula `RETURNING`, que é justamente o jeito recomendado).
- `SERIAL PRIMARY KEY` no lugar de `INTEGER PRIMARY KEY AUTOINCREMENT`.
- A migração da coluna `name` (adicionada depois na tabela `users`) ficou mais simples: o
  Postgres suporta nativamente `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`, sem precisar do
  truque de "tentar e ignorar erro de coluna duplicada" que o SQLite exigia.

### Testes automatizados

Os testes do backend (`go test ./...`) agora rodam contra um Postgres de verdade, em vez de um
SQLite em memória — isso pega problemas específicos do Postgres (como os dois pontos acima)
que só apareceriam em produção, senão. Para isso, o `docker-compose.yml` cria automaticamente
um banco separado, `minikb_test` (via `backend/db/init-test-db.sql`), isolado dos dados reais
da aplicação (banco `minikb`) — os testes limpam as tabelas desse banco de teste ao final de
cada execução.

Para rodar os testes localmente:

```bash
docker compose up -d postgres
cd backend
DATABASE_URL="postgres://minikb:minikb@localhost:5432/minikb_test?sslmode=disable" go test ./...
```

## O que ainda não está automatizado

- **CI/CD**: o Railway já faz redeploy automático a cada push na branch conectada (`deploy`),
  mas isso não inclui rodar os testes automaticamente antes — hoje isso é manual.
- **Backups do Postgres**: o plano do Railway determina a política de backup do addon; vale
  conferir a documentação deles antes de depender disso para dados importantes.
