package models

// User é o formato exposto pela API. O hash da senha nunca é serializado.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
