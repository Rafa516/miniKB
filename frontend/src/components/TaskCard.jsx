import { useState } from "react";
import { useDraggable } from "@dnd-kit/core";

// TaskCard mostra uma tarefa dentro da coluna e também funciona como o
// formulário de edição dela — os dois modos vivem no mesmo componente,
// controlados pelo estado "editing". Também é o card "arrastável" do
// drag and drop entre colunas.
function TaskCard({ task, onUpdate, onDelete }) {
  const [editing, setEditing] = useState(false);

  // Estado local de edição, pré-preenchido com os valores atuais da
  // tarefa. Só é enviado para o backend quando o formulário é salvo.
  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description);
  const [status, setStatus] = useState(task.status);

  // -----------------------------------------------------------------
  // @dnd-kit: torna este card arrastável. `id: task.id` é o que chega
  // em `active.id` no handleDragEnd do App.jsx. Fica desabilitado
  // enquanto o card está em modo de edição, para não competir com o
  // clique nos campos do formulário. `transform` é a posição do card
  // enquanto está sendo arrastado (aplicada no style abaixo).
  // -----------------------------------------------------------------
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

  // Pede confirmação antes de excluir — evita apagar por engano com um
  // clique sem querer.
  function handleDelete() {
    const confirmed = window.confirm(
      `Excluir a tarefa "${task.title}"? Essa ação não pode ser desfeita.`
    );

    if (confirmed) {
      onDelete(task.id);
    }
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
  // ação. O `ref`/`listeners`/`attributes` do @dnd-kit ficam no `div`
  // raiz, para o card inteiro ser a "alça" de arrastar — por isso os
  // botões usam onPointerDown com stopPropagation, senão um clique
  // neles seria interpretado como início de um arraste.
  // -----------------------------------------------------------------
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
