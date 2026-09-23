package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// ErrUsernameTaken é devolvido por CreateUser quando o login escolhido já
// existe (a coluna "username" tem uma constraint UNIQUE no banco).
var ErrUsernameTaken = errors.New("nome de usuário já está em uso")

// AuthRepository concentra as consultas SQL de autenticação: usuários e
// sessões. Não sabe nada de HTTP — isso fica em AuthHandler.
type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

// ---------------------------------------------------------------------
// CreateUser insere um usuário novo com a senha já em hash (o hashing em
// si é responsabilidade do handler, não do repositório). Usa "RETURNING
// id" (Postgres) em vez de Result.LastInsertId(), que o driver do
// Postgres não implementa.
// ---------------------------------------------------------------------
func (r *AuthRepository) CreateUser(name, username, passwordHash string) (int, error) {
	var id int

	err := r.DB.QueryRow(
		`INSERT INTO users (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		name, username, passwordHash,
	).Scan(&id)

	if err != nil {
		// O driver do Postgres sinaliza violação de UNIQUE com o código
		// "23505"; como não queremos depender de um pacote extra só para
		// checar esse código, a forma prática é procurar pelo texto da
		// mensagem de erro.
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique") {
			return 0, ErrUsernameTaken
		}

		return 0, err
	}

	return id, nil
}

// ---------------------------------------------------------------------
// GetUserByUsername busca um usuário pelo login, usado no fluxo de login
// para conferir a senha (o hash é comparado no handler).
// ---------------------------------------------------------------------
func (r *AuthRepository) GetUserByUsername(username string) (id int, name string, passwordHash string, err error) {
	err = r.DB.QueryRow(
		`SELECT id, name, password_hash FROM users WHERE username = $1`,
		username,
	).Scan(&id, &name, &passwordHash)

	return id, name, passwordHash, err
}

// ---------------------------------------------------------------------
// CreateSession grava um novo token de sessão, associado a um usuário e
// com prazo de validade (ver AuthHandler.startSession).
// ---------------------------------------------------------------------
func (r *AuthRepository) CreateSession(token string, userID int, expiresAt time.Time) error {
	_, err := r.DB.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt,
	)

	return err
}

// ---------------------------------------------------------------------
// GetSession valida um token de sessão: só retorna sucesso se o token
// existir E ainda não tiver expirado (comparação com o horário atual
// direto na cláusula WHERE).
// ---------------------------------------------------------------------
func (r *AuthRepository) GetSession(token string) (userID int, name string, username string, err error) {
	err = r.DB.QueryRow(`
		SELECT sessions.user_id, users.name, users.username
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token = $1 AND sessions.expires_at > $2
	`, token, time.Now()).Scan(&userID, &name, &username)

	return userID, name, username, err
}

// ---------------------------------------------------------------------
// DeleteSession remove uma sessão (usado no logout).
// ---------------------------------------------------------------------
func (r *AuthRepository) DeleteSession(token string) error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)

	return err
}
