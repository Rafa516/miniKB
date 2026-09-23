import { useState } from "react";

function Login({ onLogin, onRegister }) {
  const [mode, setMode] = useState("login");

  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  function switchMode(nextMode) {
    setMode(nextMode);
    setError("");
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      if (mode === "login") {
        await onLogin(username, password);
      } else {
        await onRegister(name, username, password);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="login-screen">
      <form className="login-card" onSubmit={handleSubmit}>
        <h1>Mini Kanban</h1>
        <p>
          {mode === "login"
            ? "Entre para gerenciar suas tarefas"
            : "Crie sua conta"}
        </p>

        {error && <p className="error-banner">{error}</p>}

        {mode === "register" && (
          <input
            type="text"
            placeholder="Nome"
            value={name}
            onChange={(event) => setName(event.target.value)}
            autoFocus
          />
        )}

        <input
          type="text"
          placeholder="Usuário"
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          autoFocus={mode === "login"}
        />

        <input
          type="password"
          placeholder="Senha"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
        />

        <button type="submit" disabled={submitting}>
          {submitting
            ? "Aguarde..."
            : mode === "login"
            ? "Entrar"
            : "Criar conta"}
        </button>

        {mode === "login" ? (
          <button
            type="button"
            className="link-button"
            onClick={() => switchMode("register")}
          >
            Não tem conta? Cadastre-se
          </button>
        ) : (
          <button
            type="button"
            className="link-button"
            onClick={() => switchMode("login")}
          >
            Já tem conta? Entrar
          </button>
        )}
      </form>
    </div>
  );
}

export default Login;
