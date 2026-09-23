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

export async function logout() {
  await fetch(`${API_URL}/logout`, {
    method: "POST",
    credentials: "include",
  });
}

export async function getCurrentUser() {
  const response = await fetch(`${API_URL}/me`, {
    credentials: "include",
  });

  if (!response.ok) {
    return null;
  }

  return response.json();
}

export async function getTasks() {
  const response = await fetch(`${API_URL}/tasks`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao buscar tarefas"));
  }

  return response.json();
}

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

export async function deleteTask(id) {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await readErrorMessage(response, "Erro ao excluir tarefa"));
  }
}
