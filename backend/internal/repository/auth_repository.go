package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

var ErrUsernameTaken = errors.New("nome de usuário já está em uso")

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUser(name, username, passwordHash string) (int, error) {
	result, err := r.DB.Exec(
		`INSERT INTO users (name, username, password_hash) VALUES (?, ?, ?)`,
		name, username, passwordHash,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return 0, ErrUsernameTaken
		}

		return 0, err
	}

	id, err := result.LastInsertId()

	return int(id), err
}

func (r *AuthRepository) GetUserByUsername(username string) (id int, name string, passwordHash string, err error) {
	err = r.DB.QueryRow(
		`SELECT id, name, password_hash FROM users WHERE username = ?`,
		username,
	).Scan(&id, &name, &passwordHash)

	return id, name, passwordHash, err
}

func (r *AuthRepository) CreateSession(token string, userID int, expiresAt time.Time) error {
	_, err := r.DB.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expiresAt,
	)

	return err
}

func (r *AuthRepository) GetSession(token string) (userID int, name string, username string, err error) {
	err = r.DB.QueryRow(`
		SELECT sessions.user_id, users.name, users.username
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token = ? AND sessions.expires_at > ?
	`, token, time.Now()).Scan(&userID, &name, &username)

	return userID, name, username, err
}

func (r *AuthRepository) DeleteSession(token string) error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)

	return err
}
