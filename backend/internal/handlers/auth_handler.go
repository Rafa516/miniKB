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

const sessionCookieName = "minikb_session"
const sessionDuration = 7 * 24 * time.Hour

type AuthHandler struct {
	Repository *repository.AuthRepository
}

func NewAuthHandler(repo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		Repository: repo,
	}
}

type registerRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)

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

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	userID, name, passwordHash, err := h.Repository.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, "usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	h.startSession(w, userID, name, req.Username)
}

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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		h.Repository.DeleteSession(cookie.Value)
	}

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

// RequireAuth protege uma rota, exigindo um cookie de sessão válido.
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

func generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
