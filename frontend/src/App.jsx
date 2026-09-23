import { useEffect, useState } from "react";
import {
  getTasks,
  createTask,
  updateTask,
  deleteTask,
} from "./services/api";
import TaskForm from "./components/TaskForm";
import KanbanColumn from "./components/KanbanColumn";

import "./App.css";

// App é o componente raiz: guarda a lista de tarefas em memória e passa
// callbacks para o formulário e para as colunas mexerem nela.
function App() {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // -----------------------------------------------------------------
  // Carrega as tarefas do backend uma vez, quando o componente monta
  // (array de dependências vazio "[]").
  // -----------------------------------------------------------------
  useEffect(() => {
    async function loadTasks() {
      try {
        const data = await getTasks();

        setTasks(data);
      } catch (err) {
        setError("Não foi possível carregar as tarefas.");
      } finally {
        setLoading(false);
      }
    }

    loadTasks();
  }, []);

  // Enquanto carrega, mostra só o texto de loading.
  if (loading) {
    return <p>Carregando tarefas...</p>;
  }

  // Se algum erro tiver acontecido (nessa carga ou em qualquer ação
  // abaixo), o quadro inteiro para de ser desenhado e só a mensagem
  // aparece na tela.
  if (error) {
    return <p>{error}</p>;
  }

  // -----------------------------------------------------------------
  // Criar tarefa: chama a API e, se der certo, insere a tarefa nova no
  // começo da lista local (sem precisar recarregar tudo do backend).
  // -----------------------------------------------------------------
  async function handleTaskCreated(task) {
    try {
      const createdTask = await createTask(task);

      setTasks((currentTasks) => [
        createdTask,
        ...currentTasks,
      ]);
    } catch (error) {
      setError("Não foi possível criar a tarefa.");
    }
  }

  // -----------------------------------------------------------------
  // Atualizar tarefa: chama a API e mescla os campos novos na tarefa
  // correspondente da lista local, pelo id.
  // -----------------------------------------------------------------
  async function handleTaskUpdate(id, task) {
    try {
      await updateTask(id, task);

      setTasks((currentTasks) =>
        currentTasks.map((currentTask) =>
          currentTask.id === id
            ? { ...currentTask, ...task }
            : currentTask
        )
      );
    } catch (error) {
      setError("Não foi possível atualizar a tarefa.");
    }
  }

  // -----------------------------------------------------------------
  // Excluir tarefa: chama a API e remove a tarefa da lista local.
  // -----------------------------------------------------------------
  async function handleTaskDelete(id) {
    try {
      await deleteTask(id);

      setTasks((currentTasks) =>
        currentTasks.filter(
          (task) => task.id !== id
        )
      );
    } catch (error) {
      setError("Não foi possível excluir a tarefa.");
    }
  }

  // -----------------------------------------------------------------
  // Layout: cabeçalho + formulário de criação + as 3 colunas do
  // Kanban, cada uma recebendo a lista inteira de tarefas e filtrando
  // pelo próprio status internamente (ver KanbanColumn.jsx).
  // -----------------------------------------------------------------
  return (
    <main className="app">
      <header className="app-header">
        <h1>Mini Kanban</h1>

        <p>Gerencie suas tarefas</p>
      </header>
      <TaskForm onTaskCreated={handleTaskCreated} />

      <div className="kanban-board">

        <KanbanColumn
         title="TODO"
         status="todo"
         tasks={tasks}
         onUpdate={handleTaskUpdate}
         onDelete={handleTaskDelete}
       />

        <KanbanColumn
         title="DOING"
         status="doing"
         tasks={tasks}
         onUpdate={handleTaskUpdate}
         onDelete={handleTaskDelete}
       />

        <KanbanColumn
        title="DONE"
        status="done"
        tasks={tasks}
        onUpdate={handleTaskUpdate}
        onDelete={handleTaskDelete}
       />
      </div>
    </main>
  );
}

export default App;
