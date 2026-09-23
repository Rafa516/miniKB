# Deploy no Render — guia passo a passo (branch `deploy`)

Este documento é o equivalente ao [DEPLOY-RAILWAY.md](DEPLOY-RAILWAY.md), mas para o
[Render](https://render.com/). Como fiz a branch `deploy` de um jeito propositalmente genérico
(porta via `$PORT`, banco via `DATABASE_URL`, endereço do backend via `VITE_API_URL`), **nenhum
código precisou mudar** para funcionar no Render — só a configuração da plataforma é diferente.
Assim como no Railway, nada aqui publica nada sozinho; são passos para você seguir manualmente.

## O que é diferente do guia do Railway

| | Railway | Render |
|---|---|---|
| Ligar o banco ao serviço | "Add Reference" no painel (automático) | copiar/colar a "Internal Database URL" manualmente (ou usar um `render.yaml`, ver seção 7) |
| Descobrir a URL pública | só existe depois de gerar o domínio | previsível desde antes: `https://<nome-do-serviço>.onrender.com` |
| Variável de build (`VITE_API_URL`) | configurada como variável do serviço | mesma coisa — o Render repassa variáveis do painel como `--build-arg` automaticamente pro Docker |
| Plano grátis | créditos limitados, não expira o projeto | plano free "de verdade", mas o serviço web **dorme** após 15 min sem uso, e o Postgres free **expira em 30 dias** |

Fonte dessas informações: documentação oficial do Render — [Docker on Render](https://render.com/docs/docker),
[Monorepo support](https://render.com/docs/monorepo-support),
[Environment Variables and Secrets](https://render.com/docs/configure-environment-variables),
[Free tier](https://render.com/docs/free) e o
[changelog sobre expiração do Postgres grátis](https://render.com/changelog/free-postgresql-instances-now-expire-after-30-days-previously-90).

## 1. Pré-requisitos

- O código já enviado para o GitHub, nesta branch (`deploy`)
- Uma conta no [Render](https://render.com/) (dá para criar com login do GitHub, sem cartão
  para o plano free)

## 2. Criar o banco PostgreSQL

1. No painel do Render: **New → PostgreSQL**.
2. Dê um nome (ex.: `minikb-db`), escolha a região e o plano (o **Free** expira em 30 dias —
   veja a nota sobre isso mais abaixo) e crie.
3. Na página do banco criado, copie a **Internal Database URL** (vai usar no passo 4). Ela só
   funciona para serviços Render na mesma conta/região — é a que você quer usar, não a
   "External" (mais lenta, pensada para acesso de fora do Render).

## 3. Criar o serviço do backend

1. **New → Web Service**, conecte o repositório `Rafa516/miniKB`.
2. Configure:
   - **Branch**: `deploy`
   - **Root Directory**: `backend`
   - **Runtime**: Docker (o Render detecta o `backend/Dockerfile` sozinho)
   - **Instance Type**: Free (ou pago, se preferir sem o "dormir" explicado abaixo)
3. Escolha um nome para o serviço (ex.: `minikb-backend`) — isso já define a URL pública dele:
   `https://minikb-backend.onrender.com` (ajuste o nome no exemplo pelo que você escolher).

## 4. Variáveis de ambiente do backend

Na aba **Environment** do serviço do backend, adicione:

| Variável | Valor |
|---|---|
| `DATABASE_URL` | a Internal Database URL copiada no passo 2 |
| `ALLOWED_ORIGINS` | a URL pública do serviço do frontend (que você vai escolher no passo 5), ex.: `https://minikb-frontend.onrender.com` |
| `ADMIN_USERNAME` | opcional, padrão `admin` |
| `ADMIN_PASSWORD` | **troque isso** — não deixe a senha padrão `admin123` em algo publicado |

Não defina `PORT` manualmente — o Render já injeta essa variável sozinho (o backend já lê ela,
com `8080` como padrão só para uso local).

## 5. Criar o serviço do frontend

1. **New → Web Service** de novo, mesmo repositório.
2. Configure:
   - **Branch**: `deploy`
   - **Root Directory**: `frontend`
   - **Runtime**: Docker (detecta `frontend/Dockerfile`)
3. Nome do serviço (ex.: `minikb-frontend`) → URL pública `https://minikb-frontend.onrender.com`.
4. Na aba **Environment** deste serviço, adicione:

   | Variável | Valor |
   |---|---|
   | `VITE_API_URL` | a URL pública do backend (passo 3), ex.: `https://minikb-backend.onrender.com` |

   O Render repassa essa variável automaticamente como `--build-arg VITE_API_URL=...` para o
   `docker build` (é assim que a documentação deles descreve o comportamento padrão para
   serviços Docker) — o `frontend/Dockerfile` já declara `ARG VITE_API_URL` esperando por ela.

## 6. Deploy

Como as URLs do Render são previsíveis pelo nome do serviço (diferente do Railway, que só gera
o domínio depois de criado), dá para preencher `ALLOWED_ORIGINS` e `VITE_API_URL` **já na
criação dos serviços**, sem precisar voltar depois para corrigir — só usar os nomes que você
escolheu nos passos 3 e 5. Depois de criar os dois serviços com as variáveis certas, o Render já
builda e sobe tudo sozinho.

## 7. (Opcional) Automatizando com um Blueprint (`render.yaml`)

Em vez de configurar tudo pelo painel, o Render suporta um arquivo `render.yaml` na raiz do
repositório que descreve os serviços e o banco de uma vez (chamado de "Blueprint"). Com ele, dá
pra referenciar a `DATABASE_URL` automaticamente (sem copiar/colar):

```yaml
databases:
  - name: minikb-db
    plan: free

services:
  - type: web
    name: minikb-backend
    runtime: docker
    rootDir: backend
    dockerfilePath: Dockerfile
    envVars:
      - key: DATABASE_URL
        fromDatabase:
          name: minikb-db
          property: connectionString
      - key: ALLOWED_ORIGINS
        value: https://minikb-frontend.onrender.com
      - key: ADMIN_PASSWORD
        sync: false   # pede o valor no painel na hora de aplicar o Blueprint, em vez de fixar aqui

  - type: web
    name: minikb-frontend
    runtime: docker
    rootDir: frontend
    dockerfilePath: Dockerfile
    envVars:
      - key: VITE_API_URL
        value: https://minikb-backend.onrender.com
```

Esse arquivo **não está incluído nesta branch** — é só um exemplo de como ficaria, caso você
prefira esse caminho no lugar da configuração manual pelo painel (passos 2–5). Se quiser, posso
adicionar esse `render.yaml` de verdade ao repositório.

## Testando localmente antes de publicar

O mesmo `docker-compose.yml` desta branch (já preparado com Postgres local) serve para testar
antes de publicar no Render — ver a seção "Testando localmente" do
[DEPLOY-RAILWAY.md](DEPLOY-RAILWAY.md#testando-localmente-antes-de-publicar), o procedimento é
idêntico (não depende de qual plataforma você vai usar depois).

## Coisas a saber sobre o plano gratuito do Render

- **O serviço web "dorme" após 15 minutos sem receber tráfego**, e o primeiro acesso depois
  disso demora cerca de 1 minuto para responder (o Render precisa religar o container). Se for
  mostrar a aplicação pra alguém, considere acessar a URL um pouco antes, ou usar um plano pago
  (que não dorme).
- **O banco Postgres grátis expira 30 dias após criado** (política atual do Render, reduzida de
  90 para 30 dias). Depois disso, ainda tem 14 dias de prazo para fazer upgrade antes dos dados
  serem apagados de vez. Para um projeto de estudo/portfólio de longa duração, isso significa
  ter que recriar o banco periodicamente (ou pagar pelo plano com banco permanente).
