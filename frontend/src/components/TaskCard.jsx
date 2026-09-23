import { useState } from "react";
import { useDraggable } from "@dnd-kit/core";

function TaskCard({ task, onUpdate, onDelete }) {
  const [editing, setEditing] = useState(false);

  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description);
  const [status, setStatus] = useState(task.status);

  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: task.id,
    disabled: editing,
  });

  const dragStyle = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
        zIndex: isDragging ? 10 : undefined,
        opacity: isDragging ? 0.6 : undefined,
      }
    : undefined;

  async function handleUpdate(event) {
    event.preventDefault();

    await onUpdate(task.id, {
      title,
      description,
      status,
    });

    setEditing(false);
  }

  function handleDelete() {
    const confirmed = window.confirm(
      `Excluir a tarefa "${task.title}"? Essa ação não pode ser desfeita.`
    );

    if (confirmed) {
      onDelete(task.id);
    }
  }

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

  return (
    <div
      className="task-card"
      ref={setNodeRef}
      style={dragStyle}
      {...listeners}
      {...attributes}
    >
      <h3>{task.title}</h3>

      {task.description && (
        <p>{task.description}</p>
      )}

      <div className="task-actions">
        <button
          onPointerDown={(event) => event.stopPropagation()}
          onClick={() => setEditing(true)}
        >
          Editar
        </button>

        <button
          onPointerDown={(event) => event.stopPropagation()}
          onClick={handleDelete}
        >
          Excluir
        </button>
      </div>
    </div>
  );
}

export default TaskCard;
