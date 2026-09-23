import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, test, expect, vi } from "vitest";
import Login from "./Login";

describe("Login", () => {
  test("chama onLogin com usuário e senha digitados", async () => {
    const onLogin = vi.fn().mockResolvedValue();

    render(<Login onLogin={onLogin} onRegister={vi.fn()} />);

    fireEvent.change(screen.getByPlaceholderText("Usuário"), {
      target: { value: "admin" },
    });

    fireEvent.change(screen.getByPlaceholderText("Senha"), {
      target: { value: "admin123" },
    });

    fireEvent.click(screen.getByText("Entrar"));

    await waitFor(() => {
      expect(onLogin).toHaveBeenCalledWith("admin", "admin123");
    });
  });

  test("mostra a mensagem de erro quando o login falha", async () => {
    const onLogin = vi.fn().mockRejectedValue(new Error("usuário ou senha inválidos"));

    render(<Login onLogin={onLogin} onRegister={vi.fn()} />);

    fireEvent.click(screen.getByText("Entrar"));

    expect(await screen.findByText("usuário ou senha inválidos")).toBeInTheDocument();
  });

  test("alterna para o cadastro e chama onRegister com nome, usuário e senha", async () => {
    const onRegister = vi.fn().mockResolvedValue();

    render(<Login onLogin={vi.fn()} onRegister={onRegister} />);

    fireEvent.click(screen.getByText("Não tem conta? Cadastre-se"));

    fireEvent.change(screen.getByPlaceholderText("Nome"), {
      target: { value: "Fulano" },
    });

    fireEvent.change(screen.getByPlaceholderText("Usuário"), {
      target: { value: "fulano" },
    });

    fireEvent.change(screen.getByPlaceholderText("Senha"), {
      target: { value: "senha123" },
    });

    fireEvent.click(screen.getByText("Criar conta"));

    await waitFor(() => {
      expect(onRegister).toHaveBeenCalledWith("Fulano", "fulano", "senha123");
    });
  });
});
