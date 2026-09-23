import TaskCard from "./TaskCard";

// KanbanColumn desenha uma coluna do quadro (TODO/DOING/DONE). Recebe a
// lista inteira de tarefas e filtra pelo próprio status — ou seja, as 3
// colunas recebem os mesmos dados e cada uma mostra só a fatia que lhe
// interessa.
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

  return (
    <section className={`kanban-column ${status}`}>
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
