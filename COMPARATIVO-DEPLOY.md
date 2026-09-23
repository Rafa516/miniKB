# Comparativo de plataformas de deploy

Este documento compara as opções de onde publicar o Mini Kanban, considerando o formato do
projeto: **dois serviços Docker** (backend em Go + frontend estático servido por Nginx) e um
**banco PostgreSQL** gerenciado. O guia já pronto pra uma delas é o
[DEPLOY-RAILWAY.md](DEPLOY-RAILWAY.md) — este documento existe pra justificar essa escolha e
deixar registrado o que se ganha/perde optando por outra.

## Tabela rápida

| Plataforma | Plano grátis | Deploy via Docker | Postgres gerenciado | "Dorme" sem uso | Cartão exigido no cadastro |
|---|:---:|:---:|:---:|:---:|:---:|
| **Railway** | Créditos grátis limitados (não é "grátis pra sempre") | Sim | Sim, no mesmo projeto | Não | Depende da promoção vigente |
| **Render** | Sim, plano free permanente | Sim | Sim, addon separado | Sim (serviços free) | Não, pro plano free |
| **Fly.io** | Créditos grátis limitados | Sim (nativo, `flyctl deploy`) | Sim (`fly postgres`) | Configurável (auto-stop) | Sim, desde 2023 |
| **Koyeb** | Sim, plano free com limitações | Sim | Não tem addon próprio (precisa de Postgres externo) | Não no free | Não, pro plano free |
| **Heroku** | Não tem mais plano grátis (desde 2022) | Sim (Container Registry) | Sim, addon pago | Não | Sim |
| **DigitalOcean App Platform** | Não | Sim | Sim, addon separado (pago) | Não | Sim |
| **Google Cloud Run + Cloud SQL** | Camada grátis generosa (uso baixo) | Sim (nativo) | Sim (Cloud SQL, separado) | Sim (escala a zero) | Sim |

## Railway (escolhido — guia em [DEPLOY-RAILWAY.md](DEPLOY-RAILWAY.md))

**Pontos altos**
- Um projeto só, com os dois serviços (backend/frontend) e o Postgres juntos, tudo se
  enxergando pela rede interna — configuração mais simples pra esse formato de monorepo.
- Detecta Dockerfile automaticamente, incluindo `railway.json` por serviço (já preparado nesta
  branch).
- Painel simples, boa experiência pra quem nunca fez deploy.
- Não "dorme": a aplicação fica sempre no ar, sem demora de "acordar" no primeiro acesso.

**Pontos baixos**
- Não é gratuito de verdade — funciona com créditos (um valor grátis por mês, geralmente
  pequeno), depois cobra por uso. Pode pedir cartão dependendo da promoção vigente no momento
  do cadastro.
- Menos maduro/conhecido que Heroku ou DigitalOcean pra quem busca uma plataforma consolidada.

## Render

**Pontos altos**
- Plano free permanente de verdade (sem cartão), inclusive pra banco Postgres (com limite de
  tempo de vida do banco free — 90 dias antes de precisar upgrade ou recriar).
- Configuração parecida com o Railway: detecta Dockerfile, painel simples, suporta monorepo via
  "Root Directory" por serviço.
- Bom para portfólio/demonstração, já que não exige cartão.

**Pontos baixos**
- Serviços no plano free **dormem** após ~15 minutos sem tráfego, e o primeiro acesso depois
  disso demora bastante (30s–1min) pra "acordar" — péssima experiência se for mostrar a
  aplicação pra alguém sem avisar.
- Banco Postgres free expira em 90 dias (precisa recriar ou pagar).

## Fly.io

**Pontos altos**
- Roda o Dockerfile quase sem adaptação (é a plataforma mais "Docker puro" da lista).
- Permite rodar réplicas em várias regiões do mundo, se algum dia isso importar.
- Config declarativa por serviço (`fly.toml`), versionável junto do código.

**Pontos baixos**
- Fluxo é mais via linha de comando (`flyctl`) do que painel visual — curva de aprendizado
  maior pra quem está começando.
- Pede cartão de crédito no cadastro desde 2023, mesmo pra usar os créditos gratuitos.
- Configurar Postgres é um passo a mais (`fly postgres create` + `fly postgres attach`), não é
  "clicar e já vem pronto" como no Railway/Render.

## Koyeb

**Pontos altos**
- Plano free existe e não pede cartão.
- Também detecta Dockerfile / conecta direto no GitHub.

**Pontos baixos**
- Não tem addon de Postgres próprio — seria preciso usar um Postgres de outro provedor (ex.:
  Neon ou Supabase, que têm free tier) e conectar via `DATABASE_URL`, adicionando mais uma
  peça pra gerenciar.
- Comunidade e documentação bem menores que as opções acima — menos exemplos prontos se algo
  der errado.

## Heroku

**Pontos altos**
- O mais tradicional/conhecido da lista, muita documentação e exemplos prontos na internet.
- Suporta deploy via Container Registry (Docker) além do fluxo clássico de buildpacks.

**Pontos baixos**
- **Não tem mais plano gratuito** desde novembro de 2022 — todo uso (mesmo mínimo) é pago,
  incluindo o addon de Postgres.
- Para um projeto pessoal/de estudo como este, é a opção mais cara da lista sem nenhuma
  vantagem técnica que compense, frente às alternativas gratuitas.

## DigitalOcean App Platform

**Pontos altos**
- Plataforma robusta, usada por empresas de verdade — bom se o projeto crescer e precisar de
  mais controle (escalabilidade, múltiplas regiões, etc.).
- Suporte a Docker e monorepo via "Source Directory" por componente.

**Pontos baixos**
- Sem plano gratuito — cobra por serviço rodando, incluindo o banco Postgres gerenciado
  (avulso, mais caro que os créditos do Railway pra um projeto pequeno).
- Mais configuração inicial que Railway/Render pra quem nunca usou.

## Google Cloud Run + Cloud SQL

**Pontos altos**
- Camada gratuita generosa pra uso baixo (Cloud Run cobra por requisição/tempo de CPU
  realmente usado — um projeto pessoal com pouco tráfego pode ficar essencialmente grátis).
- Escala automaticamente, inclusive a zero (não cobra nada parado).
- Infraestrutura de nível "produção de verdade" (é o que empresas grandes usam).

**Pontos baixos**
- De longe a opção mais complexa de configurar da lista: Cloud Run e Cloud SQL são serviços
  separados que não se enxergam direto — precisa configurar o **Cloud SQL Auth Proxy** (ou VPC
  connector) pra eles conversarem, gerenciar IAM, etc.
- Exige cartão de crédito e uma conta Google Cloud (com todo o painel/console que isso implica
  — bem menos amigável pra quem só quer publicar um projeto rápido).
- "Escalar a zero" tem uma contrapartida: se a aplicação ficar muito tempo sem uso, o próximo
  acesso demora pra "esquentar" (efeito parecido com o Render dormindo, só que geralmente mais
  rápido de voltar).

## Recomendação

Pra este projeto especificamente (Docker + Postgres, aplicação pessoal/de estudo, sem tráfego
alto esperado):

- **Railway** (já preparado) é o equilíbrio mais simples entre "fácil de configurar" e "fica
  sempre no ar" — o ponto fraco é não ser gratuito pra sempre.
- Se o objetivo for **nunca gastar nada**, **Render** é a alternativa mais direta — só aceitar
  que a aplicação demora pra responder após ficar parada, e que o banco expira em 90 dias.
- Se quiser aprender uma ferramenta mais "séria"/usada no mercado, **Google Cloud Run** é a que
  mais ensina sobre infraestrutura de verdade, ao custo de bem mais configuração inicial.
