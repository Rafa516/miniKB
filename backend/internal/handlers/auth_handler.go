package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"minikb/backend/internal/repository"
)

// Nome do cookie de sessão e por quanto tempo ele vale.
const sessionCookieName = "minikb_session"
const sessionDuration = 7 * 24 * time.Hour

// AuthHandler concentra as rotas HTTP de autenticação (registro, login,
// logout, "quem sou eu") e o middleware que protege as rotas de tarefas.
type AuthHandler struct {
	Repository *repository.AuthRepository
}

func NewAuthHandler(repo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		Repository: repo,
	}
}

// registerRequest é o formato esperado no corpo de POST /register.
type registerRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ---------------------------------------------------------------------
// POST /register — cria um usuário novo e já efetua login (cria a sessão
// na hora, sem exigir um segundo passo de "agora faça login").
// ---------------------------------------------------------------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)

	// Validação simples: nome e login não podem ser vazios, login e
	// senha têm um tamanho mínimo.
	if req.Name == "" {
		http.Error(w, "nome é obrigatório", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 {
		http.Error(w, "login deve ter pelo menos 3 caracteres", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "senha deve ter pelo menos 6 caracteres", http.StatusBadRequest)
		return
	}

	// A senha em texto puro nunca é gravada — só o hash bcrypt dela.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "erro ao processar senha", http.StatusInternalServerError)
		return
	}

	userID, err := h.Repository.CreateUser(req.Name, req.Username, string(hash))
	if err != nil {
		if err == repository.ErrUsernameTaken {
			http.Error(w, "esse login já está em uso", http.StatusConflict)
			return
		}

		http.Error(w, "erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	h.startSession(w, userID, req.Name, req.Username)
}

// loginRequest é o formato esperado no corpo de POST /login.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ---------------------------------------------------------------------
// POST /login — confere usuário e senha, e cria uma sessão nova.
// ---------------------------------------------------------------------
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	userID, name, passwordHash, err := h.Repository.GetUserByUsername(req.Username)
	if err != nil {
		// Mesma mensagem tanto para "usuário não existe" quanto para
		// "senha errada" — não dá pra um atacante descobrir logins
		// válidos só testando usuários.
		http.Error(w, "usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, "usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	h.startSession(w, userID, name, req.Username)
}

// ---------------------------------------------------------------------
// startSession gera um token aleatório, grava a sessão no banco e manda
// o cookie HttpOnly na resposta. Usado tanto por Register quanto por
// Login, já que os dois terminam com o usuário logado.
// ---------------------------------------------------------------------
func (h *AuthHandler) startSession(w http.ResponseWriter, userID int, name, username string) {
	token, err := generateToken()
	if err != nil {
		http.Error(w, "erro ao criar sessão", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(sessionDuration)

	if err := h.Repository.CreateSession(token, userID, expiresAt); err != nil {
		http.Error(w, "erro ao criar sessão", http.StatusInternalServerError)
		return
	}

	// HttpOnly: o cookie não pode ser lido por JavaScript no navegador
	// (protege contra roubo de sessão via XSS). SameSite=Lax é o padrão
	// razoável para uma aplicação comum como essa.
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"name": name, "username": username})
}

// ---------------------------------------------------------------------
// POST /logout — apaga a sessão no banco e expira o cookie no navegador.
// ---------------------------------------------------------------------
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		h.Repository.DeleteSession(cookie.Value)
	}

	// Reenvia o mesmo cookie com data de expiração no passado — é assim
	// que se "apaga" um cookie HttpOnly do lado do servidor.
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------------
// GET /me — diz quem está logado a partir do cookie de sessão. Usado
// pelo frontend ao carregar a página para saber se já existe uma sessão
// válida, sem precisar pedir login de novo.
// ---------------------------------------------------------------------
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	_, name, username, err := h.Repository.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"name": name, "username": username})
}

// ---------------------------------------------------------------------
// RequireAuth é o middleware que protege as rotas de tarefas: recebe o
// handler "de verdade" e devolve outro handler que primeiro confere se
// existe um cookie de sessão válido, só chamando o handler original se
// a checagem passar. Usado em main.go envolvendo /tasks e /tasks/{id}.
// ---------------------------------------------------------------------
func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Error(w, "não autenticado", http.StatusUnauthorized)
			return
		}

		if _, _, _, err := h.Repository.GetSession(cookie.Value); err != nil {
			http.Error(w, "não autenticado", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// generateToken cria um token de sessão aleatório (32 bytes, em hexa) a
// partir do gerador de números aleatórios criptograficamente seguro do
// Go — não é um contador nem algo previsível.
func generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
