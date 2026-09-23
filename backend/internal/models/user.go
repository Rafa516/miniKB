package models

// User é o formato de usuário exposto pela API (respostas de /login,
// /register e /me). O hash da senha nunca é serializado — nem existe
// campo pra ele aqui, só na linha da tabela "users" no banco.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
