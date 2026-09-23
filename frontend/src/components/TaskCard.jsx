import { useState } from "react";

// TaskCard mostra uma tarefa dentro da coluna e também funciona como o
// formulário de edição dela — os dois modos vivem no mesmo componente,
// controlados pelo estado "editing".
function TaskCard({ task, onUpdate, onDelete }) {
  const [editing, setEditing] = useState(false);

  // Estado local de edição, pré-preenchido com os valores atuais da
  // tarefa. Só é enviado para o backend quando o formulário é salvo.
  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description);
  const [status, setStatus] = useState(task.status);

  // Salva a edição: chama onUpdate (recebido do App.jsx) e volta pro
  // modo de visualização.
  async function handleUpdate(event) {
    event.preventDefault();

    await onUpdate(task.id, {
      title,
      description,
      status,
    });

    setEditing(false);
  }

  // -----------------------------------------------------------------
  // Modo edição: formulário com os mesmos campos do TaskForm.
  // -----------------------------------------------------------------
  if (editing) {
    return (
      <form className="task-card" onSubmit={handleUpdate}>
        <input
          value={title}
          onChange={(event) => setTitle(event.target.value)}
        />

        <input
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

        <div className="task-actions">
          <button type="submit">
            Salvar
          </button>

          <button
            type="button"
            onClick={() => setEditing(false)}
          >
            Cancelar
          </button>
        </div>
      </form>
    );
  }

  // -----------------------------------------------------------------
  // Modo visualização: título, descrição (se houver) e os botões de
  // ação. "Excluir" chama onDelete direto, sem confirmação.
  // -----------------------------------------------------------------
  return (
    <div className="task-card">
      <h3>{task.title}</h3>

      {task.description && (
        <p>{task.description}</p>
      )}

      <div className="task-actions">
        <button onClick={() => setEditing(true)}>
          Editar
        </button>

        <button onClick={() => onDelete(task.id)}>
          Excluir
        </button>
      </div>
    </div>
  );
}

export default TaskCard;
