package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"minikb/backend/internal/models"
	"minikb/backend/internal/repository"
)

// TaskHandler concentra as rotas HTTP relacionadas a tarefas. Ele não fala
// com o banco diretamente — delega isso ao Repository.
type TaskHandler struct {
	Repository *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{
		Repository: repo,
	}
}

// ---------------------------------------------------------------------
// GET /tasks — lista todas as tarefas.
// ---------------------------------------------------------------------
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repository.GetAll()

	if err != nil {
		http.Error(w, "Erro ao buscar tarefas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasks)
}

// ---------------------------------------------------------------------
// POST /tasks — cria uma tarefa nova a partir do corpo JSON da requisição.
// ---------------------------------------------------------------------
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Lê e decodifica o corpo da requisição diretamente no struct Task.
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valida os campos obrigatórios antes de gravar no banco.
	err = validateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.Repository.Create(task)
	if err != nil {
		http.Error(w, "Erro ao criar tarefa", http.StatusInternalServerError)
		return
	}

	// 201 Created + a tarefa já com o ID e created_at gerados pelo banco.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdTask)
}

// ---------------------------------------------------------------------
// PUT /tasks/{id} — atualiza uma tarefa existente.
// ---------------------------------------------------------------------
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// O roteador não tem parâmetros de rota nomeados (tipo :id), então o
	// ID é extraído "na mão" removendo o prefixo "/tasks/" do caminho.
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var task models.Task

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = validateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repository.Update(id, task)
	if err != nil {
		http.Error(w, "Erro ao atualizar tarefa", http.StatusInternalServerError)
		return
	}

	// 204: atualizado com sucesso, sem corpo de resposta.
	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------------
// DELETE /tasks/{id} — remove uma tarefa.
// ---------------------------------------------------------------------
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Repository.Delete(id)
	if err != nil {
		http.Error(w, "Erro ao excluir tarefa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------------
// validateTask centraliza as regras de validação usadas tanto na criação
// quanto na atualização de tarefas.
// ---------------------------------------------------------------------
func validateTask(task models.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return fmt.Errorf("título é obrigatório")
	}

	if strings.TrimSpace(task.Status) == "" {
		return fmt.Errorf("status é obrigatório")
	}

	switch task.Status {
	case "todo", "doing", "done":
		return nil

	default:
		return fmt.Errorf("status inválido")
	}
}
