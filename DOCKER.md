# Docker — guia passo a passo (para quem nunca usou)

Este documento explica, com calma, como rodar o Mini Kanban usando Docker. Nenhum arquivo de
código da aplicação (`.go`, `.jsx`, `.js`, etc.) foi alterado — só foram adicionados os
arquivos de configuração do Docker, listados na seção [Arquivos adicionados](#arquivos-adicionados).

## Índice

1. [O que é Docker, em palavras simples](#1-o-que-é-docker-em-palavras-simples)
2. [Conceitos básicos que você vai ver aparecer](#2-conceitos-básicos-que-você-vai-ver-aparecer)
3. [Instalando o Docker](#3-instalando-o-docker)
4. [Passo a passo: rodando o Mini Kanban](#4-passo-a-passo-rodando-o-mini-kanban)
5. [Como parar, ver o que está rodando, e apagar tudo](#5-como-parar-ver-o-que-está-rodando-e-apagar-tudo)
6. [O que cada arquivo novo faz](#6-o-que-cada-arquivo-novo-faz)
7. [Perguntas frequentes / problemas comuns](#7-perguntas-frequentes--problemas-comuns)
8. [Detalhes técnicos (para quem quiser entender o "porquê")](#8-detalhes-técnicos-para-quem-quiser-entender-o-porquê)

---

## 1. O que é Docker, em palavras simples

Hoje, pra rodar o Mini Kanban "na mão", você precisa:

- Instalar o Go na versão certa
- Instalar o Node.js na versão certa
- Rodar `go run ./cmd/server` num terminal
- Rodar `npm install` e `npm run dev` em outro terminal
- Torcer pra nenhuma versão de programa instalado no seu computador conflitar com o que o
  projeto espera

O Docker existe pra resolver exatamente essa dor de cabeça. Em vez de instalar tudo isso
diretamente no seu Windows, o Docker cria uma espécie de **"caixinha" isolada** (chamada de
**container**) que já vem com tudo o que o programa precisa pra rodar — o Go certo, o Node
certo, as bibliotecas certas — sem misturar nada com o resto do seu computador.

Analogia simples: pense no Docker como se fosse enviar o projeto inteiro **junto com o
computador que sabe rodar ele**, dentro de uma caixa fechada. Você só liga a caixa, e ela já
sabe fazer tudo funcionar do jeito certo, sempre igual, não importa em qual computador você
está.

Vantagens práticas disso aqui:

- Você não precisa instalar Go nem Node.js na sua máquina — só o Docker.
- Roda igualzinho no seu PC, no PC de outra pessoa, ou num servidor na nuvem.
- Pra "desligar" o projeto, é só desligar a caixinha — não fica nada instalado espalhado pelo
  seu sistema.

## 2. Conceitos básicos que você vai ver aparecer

| Termo | O que significa aqui |
|---|---|
| **Imagem** (*image*) | É a "receita pronta" — um pacote com tudo que o programa precisa (sistema, bibliotecas, o próprio código já compilado). Não está rodando ainda, é só o molde. |
| **Container** | É a imagem realmente **rodando** — como se fosse "ligar" a caixinha. Um container é criado a partir de uma imagem. |
| **Dockerfile** | Um arquivo de texto com a receita de como montar a imagem, passo a passo (`backend/Dockerfile` e `frontend/Dockerfile`). |
| **Docker Compose** | Uma ferramenta pra ligar **vários containers juntos** de uma vez só, com um único comando. O arquivo que descreve isso aqui é o `docker-compose.yml`, na raiz do projeto. |
| **Porta (port)** | O "endereço" numérico pelo qual você acessa o programa no navegador. Ex.: `localhost:5173` é o frontend, `localhost:8080` é o backend. |
| **Volume** | Um espaço de armazenamento que sobrevive mesmo se você desligar/apagar o container — é onde o banco de dados fica salvo (ver seção 8). |

## 3. Instalando o Docker

1. Baixe e instale o **Docker Desktop** em: https://www.docker.com/products/docker-desktop/
2. Abra o Docker Desktop e espere ele terminar de iniciar (o ícone da baleia na barra de
   tarefas fica "parado", sem animação, quando está pronto).
3. Pra confirmar que instalou certo, abra um terminal (PowerShell, por exemplo) e digite:

```bash
docker --version
```

Se aparecer um número de versão, está tudo certo.

## 4. Passo a passo: rodando o Mini Kanban

**Passo 1** — Abra um terminal (PowerShell) na pasta do projeto:

```bash
cd C:\Aplicações\miniKB
```

**Passo 2** — Rode este comando único:

```bash
docker compose up --build
```

O que acontece quando você roda isso:

1. O Docker lê o arquivo `docker-compose.yml`.
2. Ele monta a imagem do **backend** (segue a receita em `backend/Dockerfile`): baixa uma
   imagem base do Go, copia o código, compila o programa em Go.
3. Ele monta a imagem do **frontend** (segue a receita em `frontend/Dockerfile`): baixa uma
   imagem base do Node.js, instala as dependências (`npm install`), gera os arquivos finais do
   site (`npm run build`), e depois copia esses arquivos pra dentro de um servidor web leve
   (Nginx).
4. Ele "liga" as duas caixinhas (containers): uma pro backend, outra pro frontend.

Na primeira vez isso demora alguns minutos (precisa baixar as imagens base da internet). Nas
próximas vezes é bem mais rápido, porque o Docker guarda o que já baixou/montou.

Quando o terminal parar de "rolar" texto e mostrar algo como:

```text
minikb-backend-1   | Servidor rodando em http://localhost:8080
```

está tudo pronto.

**Passo 3** — Abra o navegador e acesse:

```text
http://localhost:5173
```

Isso já deve te levar direto pra tela do Mini Kanban.

> Dica: se preferir que o terminal fique livre (sem ficar preso mostrando os logs), use
> `docker compose up --build -d` — o `-d` roda tudo "em segundo plano". Nesse caso, pra ver os
> logs depois, use `docker compose logs -f`.

## 5. Como parar, ver o que está rodando, e apagar tudo

**Parar os containers** (mas manter os dados salvos e as imagens já montadas):

```bash
docker compose down
```

**Ver o que está rodando no momento:**

```bash
docker compose ps
```

Ou, visualmente, abra o **Docker Desktop** → aba **Containers**: vai aparecer um grupo
chamado `minikb` com `backend-1` e `frontend-1` dentro.

**Ligar de novo** (já estando tudo montado, não precisa do `--build` de novo a não ser que
algum Dockerfile tenha mudado):

```bash
docker compose up -d
```

**Apagar tudo mesmo** (containers, rede, e também os dados do banco salvos no volume):

```bash
docker compose down -v
```

Use esse último com cuidado — o `-v` apaga as tarefas que estiverem salvas no banco.

## 6. O que cada arquivo novo faz

| Arquivo | Pra que serve, em palavras simples |
|---|---|
| [`backend/Dockerfile`](backend/Dockerfile) | Receita pra montar a "caixinha" do backend (API em Go). |
| [`backend/.dockerignore`](backend/.dockerignore) | Lista de arquivos que o Docker deve ignorar ao montar a imagem do backend (ex.: o banco de dados local, que não deve ir junto). |
| [`frontend/Dockerfile`](frontend/Dockerfile) | Receita pra montar a "caixinha" do frontend (o site em React, servido por um Nginx). |
| [`frontend/nginx.conf`](frontend/nginx.conf) | Configuração de como o Nginx (o servidor web dentro da caixinha do frontend) deve entregar as páginas do site. |
| [`frontend/.dockerignore`](frontend/.dockerignore) | Mesma ideia do backend: arquivos que não devem ir pra dentro da imagem do frontend (ex.: `node_modules`). |
| [`docker-compose.yml`](docker-compose.yml) | O "controle remoto" que liga o backend e o frontend juntos, com um comando só. |

## 7. Perguntas frequentes / problemas comuns

**"Erro: porta já está em uso" (`port is already allocated`)**
Alguma coisa no seu computador já está usando a porta 5173 ou 8080 (pode ser outro projeto, ou
até uma execução anterior do próprio Mini Kanban rodando "na mão", fora do Docker). Feche o que
estiver usando essa porta, ou pare os containers antigos com `docker compose down`.

**"Não abre nada no navegador"**
Confira se o terminal mostrou alguma mensagem de erro. Rode `docker compose ps` pra ver se os
dois containers (`backend-1` e `frontend-1`) estão com status `running`/`Up`.

**"Mudei um arquivo do projeto e não vejo a mudança no navegador"**
O Docker, do jeito que está configurado aqui, empacota uma **versão fixa** do código dentro da
imagem (pensado pra "produção", não pra desenvolvimento do dia a dia). Se você alterar o
código, precisa reconstruir a imagem:

```bash
docker compose up --build
```

Pra desenvolvimento do dia a dia (com atualização automática ao salvar), continue usando
`npm run dev` e `go run ./cmd/server` normalmente, fora do Docker — o Docker aqui é mais voltado
pra "empacotar e rodar" a aplicação pronta.

**"Quero ver o que tem dentro do container"**
Dá pra "entrar" dentro da caixinha que está rodando com:

```bash
docker exec -it minikb-backend-1 sh
```

(troque `minikb-backend-1` por `minikb-frontend-1` pra entrar no outro). Digite `exit` pra sair.

**"Perdi as tarefas que tinha criado"**
Só acontece se você rodar `docker compose down -v` (o `-v` apaga o volume onde o banco de
dados fica guardado) ou se remover o volume manualmente pelo Docker Desktop. Sem o `-v`, os
dados ficam salvos entre uma execução e outra.

## 8. Detalhes técnicos (para quem quiser entender o "porquê")

Esta seção é mais técnica — explica as decisões tomadas ao montar os Dockerfiles.

### Por que as portas são fixas (5173 e 8080) e não podem ser trocadas livremente

O código do backend ([`backend/cmd/server/main.go`](backend/cmd/server/main.go)) libera CORS
(a permissão de um site acessar o outro) apenas para a origem `http://localhost:5173`, e o
backend escuta na porta `8080` — ambos fixos no código-fonte. Como a instrução foi não alterar
nenhum código, o `docker-compose.yml` publica os serviços exatamente nessas portas do host:

```yaml
backend:
  ports:
    - "8080:8080"
frontend:
  ports:
    - "5173:80"
```

Se você mudar essas portas no `docker-compose.yml` sem também mudar o código-fonte, o
navegador vai bloquear as requisições do frontend para o backend (erro de CORS).

### Frontend: por que existe o caminho `/miniKB/`

O `frontend/vite.config.js` já definia `base: '/miniKB/'` antes desta configuração de Docker
(usado no deploy via GitHub Pages). Para não alterar esse arquivo, o Nginx dentro do container
do frontend serve os arquivos dentro desse mesmo caminho e redireciona `/` para `/miniKB/`
(ver `frontend/nginx.conf`). Assim o mesmo build funciona sem adaptação nenhuma, tanto no
Docker quanto no GitHub Pages.

### Backend: build sem CGO

O driver `modernc.org/sqlite`, usado pelo projeto, é implementado em Go puro (não depende de
C). Por isso o backend é compilado com `CGO_ENABLED=0`, o que permite usar uma imagem final
baseada em `alpine` — bem menor e sem precisar instalar ferramentas de compilação C.

### Por que não existe um container "de banco de dados"

Diferente de PostgreSQL, MySQL etc., o SQLite não é um servidor que roda separado — é uma
**biblioteca embutida dentro do próprio programa**. O driver `modernc.org/sqlite` lê e escreve
direto num arquivo (`kanban.db`) em disco, sem processo próprio e sem porta de rede. Por isso
o `docker-compose.yml` tem só dois serviços (`backend` e `frontend`), sem um terceiro para
"banco de dados".

### Como os dados do SQLite são mantidos entre execuções

`backend/internal/database/database.go` abre o banco num caminho relativo (`kanban.db`), ou
seja, relativo à pasta onde o programa está rodando no momento. O `Dockerfile` do backend
define essa pasta como `/app/data`, que é ligada a um **volume** chamado `backend-data`. Um
volume é uma área de armazenamento gerenciada pelo Docker que continua existindo mesmo que o
container seja apagado e recriado — é o que garante que suas tarefas não somem quando você
reinicia os containers (a não ser que rode `docker compose down -v`, que apaga o volume junto).

Para conferir o arquivo do banco dentro do container em execução:

```bash
docker exec minikb-backend-1 ls -la /app/data
```

### Rodando os serviços separadamente (sem Docker Compose)

Caso queira montar e rodar cada parte na mão, sem o `docker-compose.yml`:

```bash
# Backend
cd backend
docker build -t minikb-backend .
docker run -p 8080:8080 -v minikb-backend-data:/app/data minikb-backend

# Frontend (em outro terminal)
cd frontend
docker build -t minikb-frontend .
docker run -p 5173:80 minikb-frontend
```

### Validação

Esta configuração foi testada localmente:
- `docker compose up --build` sobe os dois serviços sem erros.
- `POST`/`GET` em `/tasks` funcionam via `curl` e via interface (criação, listagem).
- Dados do SQLite persistem após recriar os containers (`docker compose up --build` de novo).
- A interface carrega corretamente em `http://localhost:5173` (com redirecionamento para
  `/miniKB/`).
