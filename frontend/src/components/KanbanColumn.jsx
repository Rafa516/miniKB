import { useDroppable } from "@dnd-kit/core";
import TaskCard from "./TaskCard";

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

  const { setNodeRef, isOver } = useDroppable({ id: status });

  return (
    <section
      className={`kanban-column ${status}${isOver ? " drag-over" : ""}`}
      ref={setNodeRef}
    >
      <div className="column-header">
        <h2>{title}</h2>

        <span>{columnTasks.length}</span>
      </div>

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
