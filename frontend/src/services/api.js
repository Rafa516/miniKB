// API_URL: endereço do backend. Se a variável de ambiente VITE_API_URL
// não for definida no build, é calculada a partir do endereço usado para
// abrir o site (protocolo + hostname) na porta 8080 — assim funciona tanto
// em localhost quanto acessando por outro IP da rede, sem precisar mexer
// em código.
const API_URL =
  import.meta.env.VITE_API_URL ||
  `${window.location.protocol}//${window.location.hostname}:8080`;

// Lê o corpo da resposta de erro (texto simples, é o que o backend Go
// devolve) para mostrar mensagens específicas em vez de um texto genérico.
async function readErrorMessage(response, fallback) {
  try {
    const text = await response.text();
    return text && text.trim() ? text.trim() : fallback;
  } catch {
    return fallback;
  }
}

// ---------------------------------------------------------------------
// Autenticação. "credentials: include" é obrigatório em todas as
// chamadas (inclusive nas de tarefas, mais abaixo) para o navegador
// enviar/receber o cookie de sessão HttpOnly do backend.
// ---------------------------------------------------------------------

// POST /register — cria a conta e já efetua login (o backend devolve o
// cookie de sessão junto).
export async function register(name, username, password) {
  const response = await fetch(`${API_URL}/register`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, username, password }),
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Não foi possível criar a conta"));
  }

  return response.json();
}

// POST /login
export async function login(username, password) {
  const response = await fetch(`${API_URL}/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Usuário ou senha inválidos"));
  }

  return response.json();
}

// POST /logout — encerra a sessão atual no backend e limpa o cookie.
export async function logout() {
  await fetch(`${API_URL}/logout`, {
    method: "POST",
    credentials: "include",
  });
}

// GET /me — usado ao carregar a página para saber se já existe uma sessão
// válida (cookie ainda não expirado). Devolve null em vez de lançar erro,
// já que "não estar logado" é um resultado esperado aqui, não uma falha.
export async function getCurrentUser() {
  const response = await fetch(`${API_URL}/me`, {
    credentials: "include",
  });

  if (!response.ok) {
    return null;
  }

  return response.json();
}

// ---------------------------------------------------------------------
// Tarefas. Todas as chamadas exigem sessão válida (o backend responde
// 401 sem o cookie), por isso também usam credentials: "include".
// ---------------------------------------------------------------------

// GET /tasks
export async function getTasks() {
  const response = await fetch(`${API_URL}/tasks`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao buscar tarefas"));
  }

  return response.json();
}

// POST /tasks
export async function createTask(task) {
  const response = await fetch(`${API_URL}/tasks`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao criar tarefa"));
  }

  return response.json();
}

// PUT /tasks/:id
export async function updateTask(id, task) {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao atualizar tarefa"));
  }
}

// DELETE /tasks/:id
export async function deleteTask(id) {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao excluir tarefa"));
  }
}
