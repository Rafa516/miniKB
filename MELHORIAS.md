# Melhorias implementadas (branch `melhorias`)

Esta branch implementa a maior parte da lista de "Possíveis melhorias" do `README.MD`. Ela foi
criada a partir de um ponto anterior à containerização (branch `main`), então também inclui os
arquivos de Docker (replicados de `main`, ver [DOCKER.md](DOCKER.md)) — mas o foco deste
documento é a aplicação em si.

**Esta branch não foi mesclada (merge) com `main`, de propósito**, pra manter o `main` estável
enquanto essas mudanças (bem maiores que só "colocar em Docker") são avaliadas.

Itens da lista do README que ficaram de fora desta rodada, por envolverem decisões de arquitetura
maiores (escolha de provedor de hospedagem, migração de banco):
- **Deploy da aplicação** — publicar em um serviço como Render/Fly.io/Railway exige criar conta
  nesses serviços, algo que só o dono do projeto pode fazer. A aplicação já está pronta pra isso
  (tem Dockerfiles prontos para produção), só falta escolher onde hospedar.
- **Banco PostgreSQL** — trocar o SQLite por Postgres é uma migração de banco (mudaria como os
  dados são armazenados e exigiria rodar um serviço de banco separado), não um ajuste pequeno.

## 🔐 Autenticação de usuários (login e cadastro)

A aplicação inteira (tarefas) agora exige estar logado. Foi implementado do jeito mais simples
possível que ainda é seguro:

- **Backend**: duas tabelas novas no SQLite, `users` (nome, login, hash da senha com bcrypt) e
  `sessions` (token de sessão, usuário, validade). Endpoints novos:
  - `POST /register` — cria um usuário (nome, usuário, senha) e já faz login automaticamente.
  - `POST /login` — confere usuário/senha e cria uma sessão (cookie `HttpOnly`, válido por 7 dias).
  - `POST /logout` — encerra a sessão.
  - `GET /me` — diz quem está logado (usado pelo frontend pra saber se já tem sessão válida ao
    abrir a página).
  - Todas as rotas de `/tasks` agora passam por um middleware (`AuthHandler.RequireAuth`) que
    exige uma sessão válida, senão responde `401`.
- Um usuário administrador é criado automaticamente ao iniciar o backend, com as credenciais das
  variáveis de ambiente `ADMIN_USERNAME`/`ADMIN_PASSWORD` (padrão local: `admin` / `admin123`).
  Ver [DOCKER.md, seção 10](DOCKER.md#10-login-e-cadastro-de-usuários) pra trocar essa senha.
- **Frontend**: tela de login/cadastro nova (`Login.jsx`), com alternância entre os dois modos.
  Enquanto não há sessão válida, nenhuma tela do quadro é mostrada. O nome do usuário aparece no
  topo ("Olá, Fulano") com botão de sair.
- Os cookies de sessão exigem `credentials: 'include'` nas chamadas do frontend e
  `Access-Control-Allow-Credentials: true` no backend — por isso o CORS (que já era restrito a
  uma lista de origens permitidas) precisou desse header a mais.

**Limitação conhecida:** não existe recuperação de senha nem confirmação por e-mail — é
propositalmente simples. Se um usuário esquecer a senha, hoje só dá pra resolver mexendo direto
no banco.

## 🖱️ Drag and Drop entre colunas

Usando `@dnd-kit` (já estava nas dependências do `package.json`, mas não era usado em lugar
nenhum). Agora dá pra arrastar um card de uma coluna pra outra pra mudar o status da tarefa —
por baixo dos panos, isso dispara a mesma chamada de "atualizar tarefa" que o formulário de
edição já usava, só que automaticamente com o novo status.

## 🔍 Filtros e busca de tarefas

Um campo de busca acima do quadro filtra as tarefas por título ou descrição (busca simples,
sem diferenciar maiúsculas/minúsculas), feita no navegador mesmo — não precisou mexer na API,
já que o quadro inteiro (as 3 colunas) já carrega todas as tarefas de uma vez.

## 📄 Paginação

Implementada na API, de um jeito que não quebra quem já usava sem paginação: `GET /tasks`
continua devolvendo todas as tarefas por padrão. Se vierem os parâmetros `?page=` e/ou
`?pageSize=`, a resposta passa a ser só aquela página, com o total de registros no cabeçalho
`X-Total-Count`.

**Por que o quadro Kanban continua carregando tudo de uma vez:** um quadro Kanban separa as
tarefas em 3 colunas por status. Paginar a lista linear (ex.: 20 tarefas por página) faria uma
página inteira vir só de tarefas "feitas", por exemplo, deixando as outras colunas vazias sem
realmente estarem vazias — uma experiência ruim. Por isso a paginação ficou pronta e testada na
API (pra um caso de uso futuro, tipo uma tela de administração em lista), mas o quadro em si
continua pedindo tudo de uma vez, que é o comportamento correto pra esse tipo de interface.

## 🗑️ Confirmação antes da exclusão

O botão "Excluir" agora pergunta antes (`window.confirm`), citando o título da tarefa. Só chama
a exclusão de verdade se a pessoa confirmar.

## 💬 Mensagens de sucesso e erro mais detalhadas

Duas mudanças aqui:

1. **O bug da tela travando foi corrigido.** Antes, qualquer erro (mesmo uma validação simples,
   tipo título vazio) fazia o `App.jsx` parar de desenhar o quadro inteiro e mostrar só a
   mensagem de erro, pra sempre, até recarregar a página. Agora o erro aparece como um aviso
   (banner vermelho) por cima do quadro, sem esconder nada, e some quando você tenta de novo.
2. **Mensagens mais específicas.** Antes, qualquer falha de rede virava um texto genérico tipo
   "Erro ao criar tarefa". Agora o frontend lê o corpo da resposta de erro do backend (que já
   trazia mensagens específicas, como "título é obrigatório") e mostra ela direto. Ações bem
   sucedidas (criar/editar/excluir) também mostram um aviso verde de confirmação por alguns
   segundos.

## 🧪 Testes automatizados

**Backend** (Go, `go test ./...`): testes de repositório (criar/listar/atualizar/excluir/
paginar tarefas, usando SQLite em memória) e de handlers HTTP (validação, criação, listagem
paginada, login, cadastro, proteção de rotas por sessão).

**Frontend** (Vitest + Testing Library, `npm test`): testes de componentes — envio do formulário
de criar tarefa, confirmação antes de excluir (aceitando e cancelando), login e cadastro
(sucesso e erro).

## 🐳 Docker

Os arquivos de Docker desta branch são os mesmos da `main` (Dockerfiles, `docker-compose.yml`,
`nginx.conf`, `.env.example`) — ver [DOCKER.md](DOCKER.md) para o guia completo. A única
diferença é que o `docker-compose.yml` agora também repassa `ADMIN_USERNAME`/`ADMIN_PASSWORD`
pro backend, por causa do login.

**Atenção — migração de banco:** se você já tinha rodado uma versão anterior (sem login) e
tem um volume `backend-data` antigo, o backend detecta e ajusta a tabela `users` sozinho ao
iniciar (adiciona a coluna que faltava) — não precisa apagar o volume nem perder as tarefas
salvas. Isso foi testado nesta branch: as tarefas de teste continuaram lá depois da atualização.
