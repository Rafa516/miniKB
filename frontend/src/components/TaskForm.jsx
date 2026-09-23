import { useState } from "react";

// TaskForm é o formulário do topo da página, usado só para criar tarefas
// novas (editar é feito direto no card, em TaskCard.jsx).
function TaskForm({ onTaskCreated }) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("todo");

  // Ao enviar: monta o objeto da tarefa, avisa o componente pai (App.jsx)
  // via onTaskCreated, e limpa os campos para o próximo cadastro.
  function handleSubmit(event) {
    event.preventDefault();

    const task = {
      title,
      description,
      status,
    };

    onTaskCreated(task);

    setTitle("");
    setDescription("");
    setStatus("todo");
  }

  return (
    <form className="task-form" onSubmit={handleSubmit}>
      {/* Campos controlados: o valor sempre vem do estado do React */}
      <input
        type="text"
        placeholder="Título da tarefa"
        value={title}
        onChange={(event) => setTitle(event.target.value)}
      />

      <input
        type="text"
        placeholder="Descrição"
        value={description}
        onChange={(event) => setDescription(event.target.value)}
      />

      <select
        value={status}
        onChange={(event) => setStatus(event.target.value)}
      >
        <option value="todo">TODO</option>
        <option value="doing">DOING</option>
        <option value="done">DONE</option>
      </select>

      <button type="submit">
        Adicionar tarefa
      </button>
    </form>
  );
}

export default TaskForm;
