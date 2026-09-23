package repository

import (
	"database/sql"
	"testing"

	"minikb/backend/internal/models"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("erro ao abrir banco de teste: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("erro ao criar tabela de teste: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	return db
}

func TestTaskRepository_CreateAndGetAll(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))

	created, err := repo.Create(models.Task{
		Title:       "Estudar Go",
		Description: "Praticar API REST",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("Create retornou erro: %v", err)
	}

	if created.ID == 0 {
		t.Fatalf("esperava um ID gerado, veio 0")
	}

	tasks, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll retornou erro: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("esperava 1 tarefa, veio %d", len(tasks))
	}

	if tasks[0].Title != "Estudar Go" {
		t.Errorf("título esperado 'Estudar Go', veio %q", tasks[0].Title)
	}
}

func TestTaskRepository_UpdateAndDelete(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))

	created, err := repo.Create(models.Task{Title: "Original", Status: "todo"})
	if err != nil {
		t.Fatalf("Create retornou erro: %v", err)
	}

	err = repo.Update(created.ID, models.Task{
		Title:       "Atualizado",
		Description: "nova descrição",
		Status:      "doing",
	})
	if err != nil {
		t.Fatalf("Update retornou erro: %v", err)
	}

	tasks, _ := repo.GetAll()
	if tasks[0].Title != "Atualizado" || tasks[0].Status != "doing" {
		t.Errorf("tarefa não foi atualizada corretamente: %+v", tasks[0])
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("Delete retornou erro: %v", err)
	}

	tasks, _ = repo.GetAll()
	if len(tasks) != 0 {
		t.Errorf("esperava 0 tarefas após excluir, veio %d", len(tasks))
	}
}

func TestTaskRepository_GetPage(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))

	for i := 0; i < 5; i++ {
		if _, err := repo.Create(models.Task{Title: "Tarefa", Status: "todo"}); err != nil {
			t.Fatalf("Create retornou erro: %v", err)
		}
	}

	page, total, err := repo.GetPage(1, 2)
	if err != nil {
		t.Fatalf("GetPage retornou erro: %v", err)
	}

	if total != 5 {
		t.Errorf("esperava total 5, veio %d", total)
	}

	if len(page) != 2 {
		t.Errorf("esperava 2 itens na página, veio %d", len(page))
	}

	page2, _, err := repo.GetPage(3, 2)
	if err != nil {
		t.Fatalf("GetPage retornou erro: %v", err)
	}

	if len(page2) != 1 {
		t.Errorf("esperava 1 item na última página, veio %d", len(page2))
	}
}
