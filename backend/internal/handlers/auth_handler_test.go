package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"minikb/backend/internal/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// newTestAuthDB conecta num Postgres real (DATABASE_URL — ver
// docker-compose.yml) com as tabelas "users"/"sessions" limpas para cada
// teste.
func newTestAuthDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL não definida — suba o Postgres local (docker compose up -d postgres) para rodar estes testes")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("erro ao abrir conexão com o banco de teste: %v", err)
	}

	_, err = db.Exec(`
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
	`)
	if err != nil {
		t.Fatalf("erro ao criar tabelas de teste: %v", err)
	}

	t.Cleanup(func() {
		db.Exec(`TRUNCATE sessions, users RESTART IDENTITY CASCADE`)
		db.Close()
	})

	return db
}

func newTestAuthHandler(t *testing.T, name, username, password string) *AuthHandler {
	t.Helper()

	db := newTestAuthDB(t)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("erro ao gerar hash: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO users (name, username, password_hash) VALUES ($1, $2, $3)`, name, username, string(hash)); err != nil {
		t.Fatalf("erro ao inserir usuário de teste: %v", err)
	}

	return NewAuthHandler(repository.NewAuthRepository(db))
}

func TestLogin_WrongPassword(t *testing.T) {
	handler := newTestAuthHandler(t, "Administrador", "admin", "senha-certa")

	body := bytes.NewBufferString(`{"username":"admin","password":"senha-errada"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, veio %d", rec.Code)
	}
}

func TestLogin_ThenMe(t *testing.T) {
	handler := newTestAuthHandler(t, "Administrador", "admin", "senha-certa")

	body := bytes.NewBufferString(`{"username":"admin","password":"senha-certa"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, veio %d (corpo: %s)", rec.Code, rec.Body.String())
	}

	var loginResp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("resposta do login não é um JSON válido: %v", err)
	}

	if loginResp["name"] != "Administrador" {
		t.Errorf("esperava name 'Administrador', veio %q", loginResp["name"])
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("esperava um cookie de sessão na resposta do login")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/me", nil)
	meReq.AddCookie(cookies[0])
	meRec := httptest.NewRecorder()

	handler.Me(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em /me, veio %d", meRec.Code)
	}
}

func TestMe_WithoutCookie(t *testing.T) {
	handler := newTestAuthHandler(t, "Administrador", "admin", "senha-certa")

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, veio %d", rec.Code)
	}
}

func TestRequireAuth_BlocksWithoutSession(t *testing.T) {
	handler := newTestAuthHandler(t, "Administrador", "admin", "senha-certa")

	called := false
	protected := handler.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	protected(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, veio %d", rec.Code)
	}

	if called {
		t.Error("o handler protegido não deveria ter sido chamado")
	}
}

func TestRegister_CreatesUserAndSession(t *testing.T) {
	db := newTestAuthDB(t)
	handler := NewAuthHandler(repository.NewAuthRepository(db))

	body := bytes.NewBufferString(`{"name":"Fulano","username":"fulano","password":"senha123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, veio %d (corpo: %s)", rec.Code, rec.Body.String())
	}

	if len(rec.Result().Cookies()) == 0 {
		t.Fatalf("esperava sessão criada automaticamente após o cadastro")
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = 'fulano'`).Scan(&count); err != nil {
		t.Fatalf("erro ao consultar usuário criado: %v", err)
	}

	if count != 1 {
		t.Fatalf("esperava 1 usuário 'fulano' no banco, veio %d", count)
	}
}

func TestRegister_RejectsDuplicateUsername(t *testing.T) {
	handler := newTestAuthHandler(t, "Administrador", "admin", "senha-certa")

	body := bytes.NewBufferString(`{"name":"Outro Admin","username":"admin","password":"senha123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("esperava 409, veio %d (corpo: %s)", rec.Code, rec.Body.String())
	}
}

func TestRegister_RejectsShortPassword(t *testing.T) {
	db := newTestAuthDB(t)
	handler := NewAuthHandler(repository.NewAuthRepository(db))

	body := bytes.NewBufferString(`{"name":"Fulano","username":"fulano","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, veio %d", rec.Code)
	}
}
