// API_URL: endereço do backend. Se a variável de ambiente VITE_API_URL
// não for definida no build, é calculada a partir do endereço usado para
// abrir o site (protocolo + hostname) na porta 8080 — assim funciona tanto
// em localhost quanto acessando por outro IP da rede, sem precisar mexer
// em código.
const API_URL =
  import.meta.env.VITE_API_URL ||
  `${window.location.protocol}//${window.location.hostname}:8080`;

// ---------------------------------------------------------------------
// GET /tasks
// ---------------------------------------------------------------------
export async function getTasks() {
  const response = await fetch(`${API_URL}/tasks`);

  if (!response.ok) {
    throw new Error("Erro ao buscar tarefas");
  }

  return response.json();
}

// ---------------------------------------------------------------------
// POST /tasks
// ---------------------------------------------------------------------
export async function createTask(task) {
  const response = await fetch(`${API_URL}/tasks`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  });

  if (!response.ok) {
    throw new Error("Erro ao criar tarefa");
  }

  return response.json();
}

// ---------------------------------------------------------------------
// PUT /tasks/:id
// ---------------------------------------------------------------------
export async function updateTask(id, task) {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  });

  if (!response.ok) {
    throw new Error("Erro ao atualizar tarefa");
  }
}

// ---------------------------------------------------------------------
// DELETE /tasks/:id
// ---------------------------------------------------------------------
export async function deleteTask(id) {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Erro ao excluir tarefa");
  }
}
