package repository

import (
	"database/sql"

	"minikb/backend/internal/models"
)

// TaskRepository é a única camada do código que sabe escrever SQL. Os
// handlers chamam esses métodos sem saber que por baixo existe SQLite.
type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		DB: db,
	}
}

// ---------------------------------------------------------------------
// GetAll busca todas as tarefas, das mais recentes para as mais antigas.
// ---------------------------------------------------------------------
func (r *TaskRepository) GetAll() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Começa com uma slice vazia (não nil) para que a API sempre devolva
	// "[]" em vez de "null" quando não houver tarefas.
	tasks := make([]models.Task, 0)

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// ---------------------------------------------------------------------
// Create insere uma tarefa nova e devolve ela já com o ID e a data de
// criação gerados pelo banco.
// ---------------------------------------------------------------------
func (r *TaskRepository) Create(task models.Task) (models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status)
		VALUES (?, ?, ?)
	`

	result, err := r.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
	)

	if err != nil {
		return models.Task{}, err
	}

	// AUTOINCREMENT: o ID da linha recém-inserida.
	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, err
	}

	task.ID = int(id)

	// created_at é preenchido pelo banco (DEFAULT CURRENT_TIMESTAMP), então
	// precisa de uma segunda consulta para saber o valor exato gravado.
	err = r.DB.QueryRow(`
		SELECT created_at
		FROM tasks
		WHERE id = ?
	`, id).Scan(&task.CreatedAt)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// ---------------------------------------------------------------------
// Update substitui título, descrição e status de uma tarefa existente.
// ---------------------------------------------------------------------
func (r *TaskRepository) Update(id int, task models.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?
		WHERE id = ?
	`

	_, err := r.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		id,
	)

	return err
}

// ---------------------------------------------------------------------
// Delete remove uma tarefa pelo ID.
// ---------------------------------------------------------------------
func (r *TaskRepository) Delete(id int) error {
	query := `
		DELETE FROM tasks
		WHERE id = ?
	`

	_, err := r.DB.Exec(query, id)

	return err
}
