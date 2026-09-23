import { render, screen, fireEvent } from "@testing-library/react";
import { describe, test, expect, vi } from "vitest";
import TaskForm from "./TaskForm";

describe("TaskForm", () => {
  test("chama onTaskCreated com os dados digitados e limpa o formulário", () => {
    const onTaskCreated = vi.fn();

    render(<TaskForm onTaskCreated={onTaskCreated} />);

    fireEvent.change(screen.getByPlaceholderText("Título da tarefa"), {
      target: { value: "Estudar Go" },
    });

    fireEvent.change(screen.getByPlaceholderText("Descrição"), {
      target: { value: "Praticar API REST" },
    });

    fireEvent.click(screen.getByText("Adicionar tarefa"));

    expect(onTaskCreated).toHaveBeenCalledWith({
      title: "Estudar Go",
      description: "Praticar API REST",
      status: "todo",
    });

    expect(screen.getByPlaceholderText("Título da tarefa").value).toBe("");
  });
});
