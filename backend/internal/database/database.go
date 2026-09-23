package database

import (
	"database/sql"

	// Driver do SQLite em Go puro (sem depender de C/CGO). O "_" na frente
	// importa o pacote só pelo efeito colateral de registrar o driver
	// "sqlite" no database/sql — o código aqui não chama nada dele
	// diretamente.
	_ "modernc.org/sqlite"
)

// Connect abre (criando se não existir) o arquivo kanban.db e garante que
// a estrutura de tabelas esteja pronta antes de devolver a conexão.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "kanban.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables cria a tabela "tasks" caso ela ainda não exista. Como usa
// "IF NOT EXISTS", é seguro chamar isso toda vez que o servidor inicia.
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)

	return err
}
