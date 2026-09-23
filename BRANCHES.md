# As branches deste repositório

Este repositório (fork) tem três branches, cada uma representando um nível diferente do
projeto. Este documento explica o que tem em cada uma, e — se você é o dono do projeto
original (ou qualquer outra pessoa) e está pensando em trazer alguma dessas mudanças pro seu
próprio repositório — como fazer isso com segurança, testando localmente antes de decidir.

## As três branches

### `main`

A versão mais próxima do projeto original, com uma adição: **Docker**. Dá pra rodar a
aplicação inteira (frontend + backend + banco) com um único comando, sem precisar instalar Go
nem Node.js na máquina. Nenhuma funcionalidade da aplicação foi alterada nesta branch — só a
forma de rodar.

- [DOCKER.md](DOCKER.md) — guia completo de Docker, passo a passo, para quem nunca usou
- [ARQUITETURA-main.md](ARQUITETURA-main.md) — como o frontend e o backend funcionam nesta versão
- [MELHORIAS.md](MELHORIAS.md) — sugestões de melhoria levantadas durante os testes (nada
  aplicado nesta branch, é só a lista)

### `melhorias`

Parte de um ponto anterior à `main` (antes do Docker existir) e implementa boa parte da lista
de "Possíveis melhorias" do `README.MD`: login e cadastro de usuários, arrastar e soltar entre
colunas, busca/filtro de tarefas, paginação na API, confirmação antes de excluir, mensagens de
erro mais claras, testes automatizados — e também tem Docker (implementado de novo aqui, já que
partiu de antes dele existir).

- [ARQUITETURA-melhorias.md](https://github.com/Rafa516/miniKB/blob/melhorias/ARQUITETURA-melhorias.md) — como tudo isso funciona
- [MELHORIAS.md](https://github.com/Rafa516/miniKB/blob/melhorias/MELHORIAS.md) (versão desta
  branch) — detalha cada melhoria implementada e por quê

### `deploy`

Parte da `melhorias` e deixa o projeto pronto para publicar num serviço de hospedagem: troca o
banco SQLite por PostgreSQL, torna a porta configurável, e documenta o passo a passo completo
para publicar no Railway ou no Render (nenhuma conta foi criada nem nada foi publicado — só
preparado).

- [DEPLOY-RAILWAY.md](https://github.com/Rafa516/miniKB/blob/deploy/DEPLOY-RAILWAY.md)
- [DEPLOY-RENDER.md](https://github.com/Rafa516/miniKB/blob/deploy/DEPLOY-RENDER.md)
- [COMPARATIVO-DEPLOY.md](https://github.com/Rafa516/miniKB/blob/deploy/COMPARATIVO-DEPLOY.md) —
  compara Railway, Render e outras plataformas

## Por que elas não estão mescladas entre si

`main` e `melhorias` partiram do **mesmo ponto** do projeto, mas de forma independente — cada
uma resolveu Docker/CORS à sua maneira, e os documentos (`DOCKER.md`, `MELHORIAS.md`) foram
escritos separadamente nas duas. Na prática, isso significa que juntar as duas gera bastante
conflito (checado com `git merge-tree`, que simula o merge sem aplicar nada): 12 arquivos em
conflito entre `main` e `melhorias`, 14 entre `main` e `deploy`. Já `melhorias` → `deploy` não
tem conflito nenhum, porque `deploy` já é construída em cima da `melhorias` (contém tudo dela).

Isso não significa que as mudanças sejam incompatíveis — é só retrabalho de resolver manualmente
algo que a `melhorias`/`deploy` já cobre (a `main` não tem nada de funcionalidade que as outras
duas não tenham; só um documento de arquitetura próprio).

## Se você quiser trazer alguma dessas branches pro seu próprio repositório

Vale testar localmente **antes** de decidir mesclar qualquer coisa na sua `main`. Duas formas:

### Opção 1 — pelo GitHub, sem usar terminal

1. No seu repositório, crie uma branch nova (ex.: `testando-melhorias`), pela própria interface
   do GitHub (botão de branches → digite o nome → Create branch).
2. Abra um Pull Request comparando **essa branch nova sua** (não a `main`) contra a branch deste
   fork que você quer testar. Um jeito rápido de montar essa comparação é editar a URL:
   ```text
   https://github.com/<seu-usuário>/<seu-repo>/compare/testando-melhorias...Rafa516:miniKB:melhorias
   ```
   (troque `melhorias` por `main` ou `deploy` conforme o que quiser olhar)
3. Mescle esse PR na sua branch de teste (não na `main`). Agora dá pra rodar/testar essa branch
   à vontade, sem nenhum risco pra sua `main`.
4. Gostou? Só então abra um segundo PR levando `testando-melhorias` pra sua `main`, quando
   estiver confiante.

### Opção 2 — pelo terminal (git), sem precisar de PR nenhum

```bash
# dentro do seu próprio repositório, já clonado localmente
git remote add rafael https://github.com/Rafa516/miniKB.git
git fetch rafael

# cria uma branch local sua, a partir da branch do fork que quer testar
git checkout -b testando-melhorias rafael/melhorias
```

Isso baixa o conteúdo da branch `melhorias` deste fork e cria uma branch local sua chamada
`testando-melhorias`, sem mexer em nada da sua `main`. Daí é só rodar normalmente (ver
`DOCKER.md`/`docker-compose.yml` da branch escolhida) e decidir depois se quer mesclar.

Se, ao tentar mesclar de verdade na sua `main`, aparecer conflito, o Git avisa exatamente quais
arquivos e onde — dá pra resolver direto pela interface web do GitHub (tem um editor simples
pra isso) ou localmente, abrindo os arquivos marcados com `<<<<<<<` / `=======` / `>>>>>>>` e
escolhendo o que manter.

### Por que testar antes de mesclar, mesmo sendo "só" um Pull Request

Um PR é uma **proposta** — nada é aplicado no seu repositório até alguém com permissão clicar em
"Merge". Mas isso não substitui testar de verdade: o diff de um PR mostra *o que* mudou, não
necessariamente *se funciona* na prática (por exemplo, a branch `deploy` só funciona com um
banco PostgreSQL rodando — isso não aparece óbvio só lendo o código). Testar numa branch local
antes, como descrito acima, evita descobrir esse tipo de coisa só depois de já ter mesclado.
