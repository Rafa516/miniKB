import { render, screen, fireEvent } from "@testing-library/react";
import { describe, test, expect, vi, afterEach } from "vitest";
import { DndContext } from "@dnd-kit/core";
import TaskCard from "./TaskCard";

const task = {
  id: 1,
  title: "Tarefa de teste",
  description: "descrição",
  status: "todo",
};

function renderCard(props = {}) {
  return render(
    <DndContext>
      <TaskCard task={task} onUpdate={vi.fn()} onDelete={vi.fn()} {...props} />
    </DndContext>
  );
}

describe("TaskCard", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  test("pede confirmação antes de excluir, e não exclui se cancelado", () => {
    vi.spyOn(window, "confirm").mockReturnValue(false);
    const onDelete = vi.fn();

    renderCard({ onDelete });

    fireEvent.click(screen.getByText("Excluir"));

    expect(window.confirm).toHaveBeenCalled();
    expect(onDelete).not.toHaveBeenCalled();
  });

  test("exclui a tarefa quando a confirmação é aceita", () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    const onDelete = vi.fn();

    renderCard({ onDelete });

    fireEvent.click(screen.getByText("Excluir"));

    expect(onDelete).toHaveBeenCalledWith(task.id);
  });
});
