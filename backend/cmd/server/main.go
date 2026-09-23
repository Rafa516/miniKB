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
	// Banco de dados: abre (ou cria) o arquivo SQLite e garante que a
	// tabela "tasks" exista.
	// ---------------------------------------------------------------
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	defer db.Close()

	// ---------------------------------------------------------------
	// Camadas da aplicação: o repositório fala com o banco, o handler
	// fala HTTP e usa o repositório para buscar/gravar dados.
	// ---------------------------------------------------------------
	taskRepository := repository.NewTaskRepository(db)

	taskHandler := handlers.NewTaskHandler(taskRepository)

	// ---------------------------------------------------------------
	// Rotas da API. Como o projeto usa só a biblioteca padrão do Go
	// (sem framework de rotas), cada caminho é registrado à mão e o
	// método HTTP (GET/POST/PUT/DELETE) é decidido dentro do handler.
	// ---------------------------------------------------------------

	// /tasks: listar (GET) e criar (POST) tarefas.
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)

		case http.MethodPost:
			taskHandler.CreateTask(w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	// /tasks/{id}: atualizar (PUT) e excluir (DELETE) uma tarefa específica.
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)

		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Servidor rodando em http://localhost:8080")

	// enableCORS envolve todas as rotas acima antes de subir o servidor.
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

// enableCORS libera o navegador a chamar essa API a partir de outra origem
// (o frontend, que roda em uma porta diferente). Sem isso, o navegador
// bloqueia a resposta por política de CORS.
//
// As origens permitidas vêm da variável de ambiente ALLOWED_ORIGINS (uma
// lista separada por vírgula); se ela não for definida, o padrão é só
// "http://localhost:5173", que é onde o frontend roda localmente.
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
