package database

import (
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	// Driver do PostgreSQL em Go puro (sem CGO). O "_" na frente importa
	// o pacote só pelo efeito colateral de registrar o driver "pgx" no
	// database/sql — o código aqui não chama nada dele diretamente.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect abre a conexão com o PostgreSQL (endereço vindo da variável de
// ambiente DATABASE_URL — no Railway, isso é preenchido automaticamente
// quando o addon de Postgres é conectado ao serviço) e deixa o banco
// pronto para uso: cria as tabelas que faltarem, aplica migrações em
// bancos mais antigos e garante que o usuário administrador exista.
func Connect() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("variável de ambiente DATABASE_URL não definida")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao Postgres: %w", err)
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

// createTables cria as tabelas "tasks", "users" e "sessions" caso ainda
// não existam. Como usa "IF NOT EXISTS", é seguro chamar isso toda vez
// que o servidor inicia — em um banco já existente, não faz nada.
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		expires_at TIMESTAMPTZ NOT NULL
	);
	`

	_, err := db.Exec(query)

	return err
}

// migrate ajusta bancos criados antes da coluna "name" existir na tabela
// users. Diferente do SQLite, o Postgres já suporta "ADD COLUMN IF NOT
// EXISTS" nativamente, então não precisa de nenhum truque para detectar
// se a coluna já existe.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT ''`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE users SET name = username WHERE name = ''`)

	return err
}

// seedAdminUser garante que exista um usuário administrador padrão,
// definido pelas variáveis de ambiente ADMIN_USERNAME/ADMIN_PASSWORD (com
// valores padrão para uso local), para já dar acesso de cara. Além dele,
// novos usuários podem se cadastrar normalmente pela tela de cadastro
// (ver AuthRepository.CreateUser / AuthHandler.Register).
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
		VALUES ($1, $2, $3)
		ON CONFLICT (username) DO UPDATE SET password_hash = excluded.password_hash
	`, "Administrador", username, string(hash))

	return err
}
