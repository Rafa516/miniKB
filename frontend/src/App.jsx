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

// App é o componente raiz. Além do estado das tarefas (como na versão
// simples do projeto), aqui ele também controla a sessão do usuário
// (login/cadastro) e a busca — tudo centralizado aqui e repassado para
// os componentes filhos via props.
function App() {
  // Sessão: enquanto `checkingSession` é true, não sabemos ainda se o
  // usuário já está logado (cookie válido) ou não.
  const [checkingSession, setCheckingSession] = useState(true);
  const [user, setUser] = useState(null);

  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [search, setSearch] = useState("");

  // -----------------------------------------------------------------
  // Ao montar o app, pergunta pro backend (GET /me) se já existe uma
  // sessão válida — assim quem já logou antes não precisa logar de novo
  // a cada vez que abre a página.
  // -----------------------------------------------------------------
  useEffect(() => {
    async function checkSession() {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
      setCheckingSession(false);
    }

    checkSession();
  }, []);

  // -----------------------------------------------------------------
  // Só carrega as tarefas depois de confirmar que há um usuário logado
  // (esse efeito roda de novo sempre que `user` muda — inclusive ao
  // logar/deslogar).
  // -----------------------------------------------------------------
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

  // -----------------------------------------------------------------
  // Busca: filtra a lista de tarefas no navegador mesmo (sem chamar a
  // API de novo), por título ou descrição. useMemo evita refazer o
  // filtro em todo re-render, só quando `tasks` ou `search` mudam.
  // -----------------------------------------------------------------
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

  // -----------------------------------------------------------------
  // Login e cadastro: em caso de sucesso, o backend já devolve os
  // dados do usuário (nome/usuário) e cria o cookie de sessão — só
  // precisamos guardar esse usuário no estado.
  // -----------------------------------------------------------------
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

  // -----------------------------------------------------------------
  // Ações sobre tarefas: cada uma limpa o erro anterior antes de
  // tentar, e mostra um aviso de sucesso/erro no lugar de esconder a
  // tela inteira (diferente da versão sem essas melhorias).
  // -----------------------------------------------------------------
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

  // -----------------------------------------------------------------
  // Drag and drop: o @dnd-kit chama isso ao soltar um card. `active.id`
  // é o id da tarefa arrastada, `over.id` é o id da coluna onde foi
  // solta (as colunas usam o próprio status como id — ver
  // KanbanColumn.jsx). Se soltou fora de uma coluna, ou na mesma coluna
  // de onde saiu, não faz nada; senão, reaproveita handleTaskUpdate
  // para persistir o novo status.
  // -----------------------------------------------------------------
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

  // Mostra um aviso de sucesso por 3 segundos e depois some sozinho.
  function showMessage(text) {
    setMessage(text);
    setTimeout(() => setMessage(""), 3000);
  }

  // -----------------------------------------------------------------
  // Ordem de decisão do que renderizar: checando sessão → sem login →
  // carregando tarefas → quadro completo.
  // -----------------------------------------------------------------
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

          {/* Nome do usuário logado + botão de sair */}
          <div className="user-bar">
            <span>Olá, {user.name}</span>
            <button onClick={handleLogout}>Sair</button>
          </div>
        </div>
      </header>

      {/* Avisos: aparecem por cima do quadro, sem esconder o resto da tela */}
      {error && <p className="error-banner">{error}</p>}
      {message && <p className="success-banner">{message}</p>}

      <TaskForm onTaskCreated={handleTaskCreated} />

      {/* Campo de busca: filtra `filteredTasks`, usado nas 3 colunas abaixo */}
      <input
        type="text"
        className="search-input"
        placeholder="Buscar por título ou descrição..."
        value={search}
        onChange={(event) => setSearch(event.target.value)}
      />

      {/* DndContext: envolve as colunas para habilitar arrastar/soltar
          cards entre elas (ver handleDragEnd acima) */}
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
