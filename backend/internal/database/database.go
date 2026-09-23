package database

import (
	"database/sql"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "kanban.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	if err := seedAdminUser(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err := db.Exec(query)

	return err
}

// migrate ajusta bancos criados antes da coluna "name" existir na tabela
// users (CREATE TABLE IF NOT EXISTS não altera tabelas já existentes). Em um
// banco novo a coluna já nasce criada e o ALTER abaixo só é ignorado.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`ALTER TABLE users ADD COLUMN name TEXT NOT NULL DEFAULT ''`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}

	_, err = db.Exec(`UPDATE users SET name = username WHERE name = ''`)

	return err
}

// seedAdminUser garante que exista um usuário administrador padrão, definido
// pelas variáveis de ambiente ADMIN_USERNAME/ADMIN_PASSWORD (com valores
// padrão para uso local), para já dar acesso de cara. Além dele, novos
// usuários podem se cadastrar normalmente pela tela de cadastro (ver
// AuthRepository.CreateUser / AuthHandler.Register).
func seedAdminUser(db *sql.DB) error {
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "admin123"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users (name, username, password_hash)
		VALUES (?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET password_hash = excluded.password_hash
	`, "Administrador", username, string(hash))

	return err
}
