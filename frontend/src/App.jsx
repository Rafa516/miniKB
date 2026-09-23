import { useEffect, useMemo, useState } from "react";
import { DndContext } from "@dnd-kit/core";
import {
  getTasks,
  createTask,
  updateTask,
  deleteTask,
  login,
  register,
  logout,
  getCurrentUser,
} from "./services/api";
import TaskForm from "./components/TaskForm";
import KanbanColumn from "./components/KanbanColumn";
import Login from "./components/Login";

import "./App.css";

function App() {
  const [checkingSession, setCheckingSession] = useState(true);
  const [user, setUser] = useState(null);

  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [search, setSearch] = useState("");

  useEffect(() => {
    async function checkSession() {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
      setCheckingSession(false);
    }

    checkSession();
  }, []);

  useEffect(() => {
    if (!user) {
      return;
    }

    async function loadTasks() {
      try {
        const data = await getTasks();

        setTasks(data);
      } catch (err) {
        setError(err.message || "Não foi possível carregar as tarefas.");
      } finally {
        setLoading(false);
      }
    }

    loadTasks();
  }, [user]);

  const filteredTasks = useMemo(() => {
    const term = search.trim().toLowerCase();

    if (!term) {
      return tasks;
    }

    return tasks.filter(
      (task) =>
        task.title.toLowerCase().includes(term) ||
        (task.description || "").toLowerCase().includes(term)
    );
  }, [tasks, search]);

  async function handleLogin(username, password) {
    const loggedUser = await login(username, password);
    setUser(loggedUser);
  }

  async function handleRegister(name, username, password) {
    const registeredUser = await register(name, username, password);
    setUser(registeredUser);
  }

  async function handleLogout() {
    await logout();
    setUser(null);
    setTasks([]);
  }

  async function handleTaskCreated(task) {
    setError("");

    try {
      const createdTask = await createTask(task);

      setTasks((currentTasks) => [createdTask, ...currentTasks]);
      showMessage("Tarefa criada com sucesso.");
    } catch (err) {
      setError(err.message || "Não foi possível criar a tarefa.");
    }
  }

  async function handleTaskUpdate(id, task) {
    setError("");

    try {
      await updateTask(id, task);

      setTasks((currentTasks) =>
        currentTasks.map((currentTask) =>
          currentTask.id === id
            ? { ...currentTask, ...task }
            : currentTask
        )
      );
      showMessage("Tarefa atualizada com sucesso.");
    } catch (err) {
      setError(err.message || "Não foi possível atualizar a tarefa.");
    }
  }

  async function handleTaskDelete(id) {
    setError("");

    try {
      await deleteTask(id);

      setTasks((currentTasks) =>
        currentTasks.filter((task) => task.id !== id)
      );
      showMessage("Tarefa excluída com sucesso.");
    } catch (err) {
      setError(err.message || "Não foi possível excluir a tarefa.");
    }
  }

  function handleDragEnd(event) {
    const { active, over } = event;

    if (!over) {
      return;
    }

    const task = tasks.find((t) => t.id === active.id);
    const newStatus = over.id;

    if (!task || task.status === newStatus) {
      return;
    }

    handleTaskUpdate(task.id, { ...task, status: newStatus });
  }

  function showMessage(text) {
    setMessage(text);
    setTimeout(() => setMessage(""), 3000);
  }

  if (checkingSession) {
    return <p>Carregando...</p>;
  }

  if (!user) {
    return <Login onLogin={handleLogin} onRegister={handleRegister} />;
  }

  if (loading) {
    return <p>Carregando tarefas...</p>;
  }

  return (
    <main className="app">
      <header className="app-header">
        <div className="app-header-top">
          <div>
            <h1>Mini Kanban</h1>
            <p>Gerencie suas tarefas</p>
          </div>

          <div className="user-bar">
            <span>Olá, {user.name}</span>
            <button onClick={handleLogout}>Sair</button>
          </div>
        </div>
      </header>

      {error && <p className="error-banner">{error}</p>}
      {message && <p className="success-banner">{message}</p>}

      <TaskForm onTaskCreated={handleTaskCreated} />

      <input
        type="text"
        className="search-input"
        placeholder="Buscar por título ou descrição..."
        value={search}
        onChange={(event) => setSearch(event.target.value)}
      />

      <DndContext onDragEnd={handleDragEnd}>
        <div className="kanban-board">
          <KanbanColumn
            title="TODO"
            status="todo"
            tasks={filteredTasks}
            onUpdate={handleTaskUpdate}
            onDelete={handleTaskDelete}
          />

          <KanbanColumn
            title="DOING"
            status="doing"
            tasks={filteredTasks}
            onUpdate={handleTaskUpdate}
            onDelete={handleTaskDelete}
          />

          <KanbanColumn
            title="DONE"
            status="done"
            tasks={filteredTasks}
            onUpdate={handleTaskUpdate}
            onDelete={handleTaskDelete}
          />
        </div>
      </DndContext>
    </main>
  );
}

export default App;
