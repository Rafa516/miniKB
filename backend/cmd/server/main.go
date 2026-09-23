package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"minikb/backend/internal/database"
	"minikb/backend/internal/handlers"
	"minikb/backend/internal/repository"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	defer db.Close()

	taskRepository := repository.NewTaskRepository(db)

	taskHandler := handlers.NewTaskHandler(taskRepository)

	authRepository := repository.NewAuthRepository(db)

	authHandler := handlers.NewAuthHandler(authRepository)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Register(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Login(w, r)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Logout(w, r)
	})

	http.HandleFunc("/me", authHandler.Me)

	http.HandleFunc("/tasks", authHandler.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)

		case http.MethodPost:
			taskHandler.CreateTask(w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/tasks/", authHandler.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)

		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("Servidor rodando em http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

// enableCORS libera as origens listadas em ALLOWED_ORIGINS (separadas por
// vírgula). Sem essa variável, o padrão é http://localhost:5173, igual ao
// comportamento original. Access-Control-Allow-Credentials precisa estar
// presente porque o login agora depende de cookie de sessão.
func enableCORS(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if origin == strings.TrimSpace(allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
