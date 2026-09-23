# Docker — guia passo a passo (para quem nunca usou)

Este documento explica, com calma, como rodar o Mini Kanban usando Docker. Esta branch
(`melhorias`) também inclui um conjunto de melhorias na aplicação em si — login/cadastro,
arrastar e soltar, busca, paginação, testes automatizados etc. — detalhadas em
[MELHORIAS.md](MELHORIAS.md). Este documento foca só na parte de Docker.

## Índice

1. [O que é Docker, em palavras simples](#1-o-que-é-docker-em-palavras-simples)
2. [Conceitos básicos que você vai ver aparecer](#2-conceitos-básicos-que-você-vai-ver-aparecer)
3. [Instalando o Docker](#3-instalando-o-docker)
4. [Passo a passo: rodando o Mini Kanban](#4-passo-a-passo-rodando-o-mini-kanban)
5. [Como parar, ver o que está rodando, e apagar tudo](#5-como-parar-ver-o-que-está-rodando-e-apagar-tudo)
6. [O que cada arquivo novo faz](#6-o-que-cada-arquivo-novo-faz)
7. [Perguntas frequentes / problemas comuns](#7-perguntas-frequentes--problemas-comuns)
8. [Detalhes técnicos (para quem quiser entender o "porquê")](#8-detalhes-técnicos-para-quem-quiser-entender-o-porquê)
9. [Acessando de outra máquina na rede](#9-acessando-de-outra-máquina-na-rede)
10. [Login e cadastro de usuários](#10-login-e-cadastro-de-usuários)

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

Isso já deve te levar direto pra tela de login do Mini Kanban.

> Dica: se preferir que o terminal fique livre (sem ficar preso mostrando os logs), use
> `docker compose up --build -d` — o `-d` roda tudo "em segundo plano". Nesse caso, pra ver os
> logs depois, use `docker compose logs -f`.

**Passo 4** — Entre com o usuário administrador que já vem configurado por padrão:

- Usuário: `admin`
- Senha: `admin123`

Ou clique em "Não tem conta? Cadastre-se" pra criar seu próprio usuário (nome, login e senha).
Veja a seção [10](#10-login-e-cadastro-de-usuários) pra entender como isso funciona e como
trocar a senha padrão do administrador.

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

## 9. Acessando de outra máquina na rede

Por padrão, o Mini Kanban só funciona corretamente quando acessado como `http://localhost:5173`,
no mesmo computador onde o Docker está rodando. Tentar abrir pelo IP da máquina (por exemplo,
`http://10.40.68.69:5173`, a partir de outro computador da rede) dava erro de "Não foi possível
carregar as tarefas.", por dois motivos:

1. O endereço do backend usado pelo site (`http://localhost:8080`) estava fixo dentro do código
   do frontend. Quando o site é aberto a partir de outra máquina, "localhost" naquele contexto é
   a própria máquina de quem está acessando — não o servidor —, então a chamada nunca chegava
   em lugar nenhum.
2. O backend só aceitava pedidos vindos exatamente de `http://localhost:5173` (checagem de CORS
   fixa no código). Um acesso vindo de outro endereço era bloqueado pelo navegador.

Para resolver isso, **foram feitas duas pequenas alterações de código** (documentadas com
detalhes na mensagem do commit correspondente):

- **[`frontend/src/services/api.js`](frontend/src/services/api.js):** em vez de sempre apontar
  pra `http://localhost:8080`, o endereço do backend agora é calculado a partir do endereço que
  a pessoa usou pra abrir o site (`window.location.hostname`). Assim, funciona automaticamente
  tanto em `localhost` quanto em qualquer IP da rede, sem precisar reconfigurar nada.
- **[`backend/cmd/server/main.go`](backend/cmd/server/main.go):** a origem liberada no CORS
  deixou de ser um valor fixo no código e passou a vir da variável de ambiente
  `ALLOWED_ORIGINS` (lista separada por vírgula), com o mesmo valor de antes
  (`http://localhost:5173`) como padrão — ou seja, quem não mexer em nada continua com o
  comportamento de sempre.

### Como habilitar o acesso pela rede

1. Descubra o IP da máquina que roda o Docker (no Windows: `ipconfig`, procure por "Endereço
   IPv4").
2. Crie um arquivo chamado `.env` na raiz do projeto (pode copiar o `.env.example` que já existe)
   com o conteúdo:

   ```text
   ALLOWED_ORIGINS=http://localhost:5173,http://<seu-ip-aqui>:5173
   ```

3. Suba os containers de novo:

   ```bash
   docker compose up -d --build
   ```

4. Nas outras máquinas da rede (mesma rede local/Wi-Fi), acesse `http://<seu-ip>:5173`.

O arquivo `.env` não é enviado ao Git (está no `.gitignore`), porque o IP é específico da sua
máquina/rede — cada pessoa que for rodar o projeto configura o seu.

## 10. Login e cadastro de usuários

Nesta branch, a aplicação passou a exigir login (ver [MELHORIAS.md](MELHORIAS.md) para o porquê
e os detalhes técnicos). Dois jeitos de entrar:

**1. Usuário administrador padrão**, criado automaticamente ao subir o backend:
- Usuário: `admin`
- Senha: `admin123`

**2. Criar sua própria conta**, pelo link "Não tem conta? Cadastre-se" na tela de login
(informe nome, usuário e senha — sem confirmação por e-mail, é um cadastro simples).

### Trocando a senha do administrador padrão

Adicione ao seu arquivo `.env` (o mesmo mencionado na seção 9):

```text
ADMIN_USERNAME=admin
ADMIN_PASSWORD=uma-senha-mais-forte-aqui
```

E suba os containers de novo (`docker compose up -d --build`). Isso é importante caso você
exponha a aplicação além da sua própria máquina — a senha padrão `admin123` é só para
facilitar o primeiro acesso local.

### Onde ficam os usuários cadastrados

Os usuários (nome, login e senha com hash) ficam salvos nas mesmas tabelas do banco SQLite
descrito na seção 8 — ou seja, no volume `backend-data`, persistindo entre reinícios dos
containers.
