package repository

import (
	"database/sql"

	"minikb/backend/internal/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		DB: db,
	}
}

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

// GetPage retorna uma página de tarefas (mais recentes primeiro) e o total
// de tarefas existentes, para montar a paginação no header da resposta.
func (r *TaskRepository) GetPage(page, pageSize int) ([]models.Task, int, error) {
	offset := (page - 1) * pageSize

	rows, err := r.DB.Query(`
		SELECT id, title, description, status, created_at
		FROM tasks
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

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
			return nil, 0, err
		}

		tasks = append(tasks, task)
	}

	var total int

	if err := r.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&total); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

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

	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, err
	}

	task.ID = int(id)

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
func (r *TaskRepository) Delete(id int) error {
	query := `
		DELETE FROM tasks
		WHERE id = ?
	`

	_, err := r.DB.Exec(query, id)

	return err
}
