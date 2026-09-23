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

	// ---------------------------------------------------------------
	// Banco de dados: abre (ou cria) o arquivo SQLite, garante que as
	// tabelas existam, aplica migrações pendentes e cria/atualiza o
	// usuário administrador padrão (ver internal/database/database.go).
	// ---------------------------------------------------------------
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	defer db.Close()

	// ---------------------------------------------------------------
	// Camadas da aplicação: repositórios falam com o banco, handlers
	// falam HTTP e usam os repositórios para buscar/gravar dados.
	// ---------------------------------------------------------------
	taskRepository := repository.NewTaskRepository(db)

	taskHandler := handlers.NewTaskHandler(taskRepository)

	authRepository := repository.NewAuthRepository(db)

	authHandler := handlers.NewAuthHandler(authRepository)

	// ---------------------------------------------------------------
	// Rotas de autenticação (públicas — não exigem sessão).
	// ---------------------------------------------------------------

	// /register: cria um usuário novo e já retorna logado.
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Register(w, r)
	})

	// /login: confere usuário/senha e cria a sessão (cookie).
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Login(w, r)
	})

	// /logout: encerra a sessão atual.
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Logout(w, r)
	})

	// /me: diz quem está logado a partir do cookie de sessão (usado
	// pelo frontend ao carregar a página, para saber se já tem sessão).
	http.HandleFunc("/me", authHandler.Me)

	// ---------------------------------------------------------------
	// Rotas de tarefas — protegidas por authHandler.RequireAuth, que
	// exige um cookie de sessão válido antes de deixar a requisição
	// chegar no handler de verdade (responde 401 caso contrário).
	// ---------------------------------------------------------------

	// /tasks: listar (GET) e criar (POST) tarefas.
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

	// /tasks/{id}: atualizar (PUT) e excluir (DELETE) uma tarefa específica.
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

	// enableCORS envolve todas as rotas acima antes de subir o servidor.
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

// enableCORS libera as origens listadas em ALLOWED_ORIGINS (separadas por
// vírgula). Sem essa variável, o padrão é http://localhost:5173, igual ao
// comportamento original. Access-Control-Allow-Credentials precisa estar
// presente porque o login agora depende de cookie de sessão — sem esse
// header, o navegador não guarda nem reenvia o cookie entre domínios.
func enableCORS(next http.Handler) http.Handler {
	// ---------------------------------------------------------------
	// Lê a lista de origens permitidas uma única vez, na inicialização
	// do servidor (não a cada requisição).
	// ---------------------------------------------------------------
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// -----------------------------------------------------------
		// Só ecoa o header de CORS de volta se a origem da requisição
		// bater exatamente com alguma da lista permitida.
		// -----------------------------------------------------------
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

		// -----------------------------------------------------------
		// Requisição de "preflight": o navegador manda um OPTIONS antes
		// do pedido de verdade, só para checar se tem permissão. Aqui
		// só respondemos 204 (sem corpo) e paramos — não repassamos
		// esse OPTIONS pro handler de verdade.
		// -----------------------------------------------------------
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
