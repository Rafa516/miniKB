package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"minikb/backend/internal/models"
	"minikb/backend/internal/repository"

	_ "modernc.org/sqlite"
)

func newTestHandler(t *testing.T) *TaskHandler {
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

	return NewTaskHandler(repository.NewTaskRepository(db))
}

func TestValidateTask(t *testing.T) {
	cases := []struct {
		name    string
		task    models.Task
		wantErr bool
	}{
		{"válida", models.Task{Title: "Fazer algo", Status: "todo"}, false},
		{"sem título", models.Task{Title: "  ", Status: "todo"}, true},
		{"sem status", models.Task{Title: "Fazer algo", Status: ""}, true},
		{"status inválido", models.Task{Title: "Fazer algo", Status: "invalido"}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateTask(c.task)
			if (err != nil) != c.wantErr {
				t.Errorf("validateTask(%+v) erro = %v, wantErr = %v", c.task, err, c.wantErr)
			}
		})
	}
}

func TestCreateTask_RejectsEmptyTitle(t *testing.T) {
	handler := newTestHandler(t)

	body := bytes.NewBufferString(`{"title":"","status":"todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, veio %d (corpo: %s)", rec.Code, rec.Body.String())
	}
}

func TestCreateTask_ThenGetTasks(t *testing.T) {
	handler := newTestHandler(t)

	body := bytes.NewBufferString(`{"title":"Nova tarefa","description":"desc","status":"todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava 201, veio %d (corpo: %s)", rec.Code, rec.Body.String())
	}

	var created models.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("resposta não é um JSON de tarefa válido: %v", err)
	}

	if created.Title != "Nova tarefa" {
		t.Errorf("título esperado 'Nova tarefa', veio %q", created.Title)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	listRec := httptest.NewRecorder()

	handler.GetTasks(listRec, listReq)

	var tasks []models.Task
	if err := json.Unmarshal(listRec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("resposta da listagem não é um JSON válido: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("esperava 1 tarefa na listagem, veio %d", len(tasks))
	}
}

func TestGetTasks_Paginated(t *testing.T) {
	handler := newTestHandler(t)

	for i := 0; i < 3; i++ {
		body := bytes.NewBufferString(`{"title":"Tarefa","status":"todo"}`)
		req := httptest.NewRequest(http.MethodPost, "/tasks", body)
		handler.CreateTask(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks?page=1&pageSize=2", nil)
	rec := httptest.NewRecorder()

	handler.GetTasks(rec, req)

	if rec.Header().Get("X-Total-Count") != "3" {
		t.Errorf("esperava X-Total-Count 3, veio %q", rec.Header().Get("X-Total-Count"))
	}

	var tasks []models.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("resposta não é um JSON válido: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("esperava 2 tarefas na página, veio %d", len(tasks))
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, veio %d", rec.Code)
	}
}
