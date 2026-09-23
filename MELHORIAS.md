# Sugestões de melhoria

Anotações levantadas durante a containerização do projeto (ver [DOCKER.md](DOCKER.md)),
enquanto testava a aplicação rodando em Docker. Nenhuma dessas mudanças foi aplicada —
é só a lista para repassar pro pessoal que mantém o projeto.

## 🐞 Bug: a tela trava permanentemente após qualquer erro

Em [`frontend/src/App.jsx`](frontend/src/App.jsx), o componente `App` funciona assim:

```jsx
if (loading) {
  return <p>Carregando tarefas...</p>;
}

if (error) {
  return <p>{error}</p>;
}
```

O problema: esses dois `if` valem pra **qualquer** render do componente, não só pro
carregamento inicial. Assim que `error` é preenchido — seja por falha ao carregar, criar,
editar ou excluir uma tarefa — o componente para de desenhar o cabeçalho, o formulário e as
colunas do quadro, e passa a mostrar só a frase do erro, pra sempre. Não tem como voltar pro
quadro sem recarregar a página inteira (F5).

Foi exatamente isso que aconteceu durante os testes. **Como reproduzir:** clique em
"Adicionar tarefa" sem preencher o campo "Título" — o backend responde `400 Bad Request`
("título é obrigatório"), e a tela trava mostrando só "Não foi possível criar a tarefa.", sem
quadro, sem formulário, em qualquer navegador. Só recarregando a página (F5) o quadro volta.
Ou seja, não é um problema de rede/Docker: é um erro simples de validação sendo tratado de um
jeito que quebra a tela inteira.

**Sugestão:** guardar erro e dados em estados separados, e mostrar o erro como um aviso/banner
por cima do quadro — sem esconder o resto da interface. Por exemplo:

```jsx
{error && <p className="error-banner">{error}</p>}

<TaskForm onTaskCreated={handleTaskCreated} />
<div className="kanban-board">...</div>
```

E limpar o erro (`setError("")`) no início de cada ação, antes de tentar de novo.

## 🗑️ Modal de confirmação ao excluir

Em [`frontend/src/components/TaskCard.jsx`](frontend/src/components/TaskCard.jsx), o botão
"Excluir" chama `onDelete(task.id)` direto, sem nenhuma confirmação:

```jsx
<button onClick={() => onDelete(task.id)}>
  Excluir
</button>
```

Um clique sem querer apaga a tarefa na hora, sem chance de desfazer. Vale adicionar uma
confirmação (um `window.confirm(...)` simples já resolve, ou um modal customizado se quiser
algo mais bonito) antes de chamar `onDelete`.

## Outras sugestões

- **Botões não desabilitam durante o envio**: em `TaskForm` e `TaskCard`, nada impede o
  usuário de clicar em "Adicionar"/"Salvar" várias vezes seguidas enquanto a requisição ainda
  está em andamento, podendo criar/duplicar registros. Vale desabilitar o botão (ou mostrar um
  spinner) enquanto a promise não resolve.

- **`Update` e `Delete` não verificam se a tarefa existe**: em
  [`backend/internal/repository/task_repository.go`](backend/internal/repository/task_repository.go),
  os métodos `Update` e `Delete` não checam `RowsAffected()`. Se o `id` não existir, a API
  mesmo assim responde `204 No Content` como se tivesse dado certo, em vez de `404 Not Found`.

- **CORS e porta fixos no código**: em
  [`backend/cmd/server/main.go`](backend/cmd/server/main.go), tanto a porta (`8080`) quanto a
  origem liberada no CORS (`http://localhost:5173`) estão hardcoded. Isso é justamente o que
  obrigou a fixar as portas do Docker Compose (ver DOCKER.md). Ler esses valores de variáveis
  de ambiente (com esses mesmos valores como padrão) deixaria o backend mais flexível para
  rodar em outras portas/ambientes sem precisar recompilar.

- **Sem testes automatizados**: já está listado como melhoria futura no `README.MD`, mas vale
  reforçar — nem backend nem frontend têm testes hoje.

- **Sem limite de tamanho para título/descrição**: a validação em
  [`backend/internal/handlers/task_handler.go`](backend/internal/handlers/task_handler.go) só
  checa se o título e o status estão vazios, mas não limita o tamanho de nenhum campo.
