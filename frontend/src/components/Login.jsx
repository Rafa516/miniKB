import { useState } from "react";

// Login é a tela mostrada quando ninguém está autenticado (ver App.jsx).
// Um único componente cobre os dois fluxos — entrar e criar conta — só
// alternando o que é mostrado conforme o estado `mode`.
function Login({ onLogin, onRegister }) {
  const [mode, setMode] = useState("login"); // "login" | "register"

  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  // Troca de modo e limpa qualquer erro da tentativa anterior.
  function switchMode(nextMode) {
    setMode(nextMode);
    setError("");
  }

  // -----------------------------------------------------------------
  // Envio do formulário: chama onLogin ou onRegister (recebidos do
  // App.jsx) dependendo do modo atual. Os dois vêm de services/api.js
  // e lançam erro em caso de falha (usuário/senha errados, login já
  // em uso, etc.), que é capturado aqui e mostrado na tela.
  // -----------------------------------------------------------------
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

        {/* Campo "Nome" só aparece no modo de cadastro */}
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

        {/* Link para alternar entre login e cadastro */}
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
