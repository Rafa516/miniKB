import { useDroppable } from "@dnd-kit/core";
import TaskCard from "./TaskCard";

// KanbanColumn desenha uma coluna do quadro (TODO/DOING/DONE). Recebe a
// lista inteira de tarefas (já filtrada pela busca, se houver) e filtra
// de novo pelo próprio status — ou seja, as 3 colunas recebem os mesmos
// dados e cada uma mostra só a fatia que lhe interessa. Também é a "zona
// de soltar" (drop target) do drag and drop.
function KanbanColumn({
  title,
  status,
  tasks,
  onUpdate,
  onDelete,
}) {
  const columnTasks = tasks.filter(
    (task) => task.status === status
  );

  // -----------------------------------------------------------------
  // @dnd-kit: torna a coluna uma área onde um card pode ser solto.
  // `id: status` é o que chega em `over.id` no handleDragEnd do
  // App.jsx — por isso as colunas usam o próprio status como id, sem
  // precisar de um id separado. `isOver` fica true só enquanto um card
  // está sendo arrastado por cima dela (usado para o destaque visual).
  // -----------------------------------------------------------------
  const { setNodeRef, isOver } = useDroppable({ id: status });

  return (
    <section
      className={`kanban-column ${status}${isOver ? " drag-over" : ""}`}
      ref={setNodeRef}
    >
      {/* Cabeçalho: nome da coluna + contador de tarefas */}
      <div className="column-header">
        <h2>{title}</h2>

        <span>{columnTasks.length}</span>
      </div>

      {/* Um TaskCard por tarefa da coluna */}
      <div className="column-content">
        {columnTasks.map((task) => (
          <TaskCard
            key={task.id}
            task={task}
            onUpdate={onUpdate}
            onDelete={onDelete}
          />
        ))}
      </div>
    </section>
  );
}

export default KanbanColumn;
