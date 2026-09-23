package models

import "time"

// Task representa uma tarefa do quadro Kanban, tanto na linha do banco
// quanto no JSON trocado com o frontend (as tags `json:"..."` controlam os
// nomes dos campos na resposta da API).
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // um de: "todo", "doing", "done"
	CreatedAt   time.Time `json:"createdAt"`
}
