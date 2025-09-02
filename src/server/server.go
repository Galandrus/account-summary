package server

import (
	"account-summary/src/config"
	"account-summary/src/interfaces/handlers"
	"account-summary/src/interfaces/server"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func NewServer(cfg *config.Config, handler handlers.MainApiHandlerInterface) server.ServerInterface {
	port := cfg.Port
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "frontend.html")
		} else {
			http.NotFound(w, r)
		}
	})

	// Rutas de la API
	mux.HandleFunc("/load-transactions", handler.LoadTransactions)
	mux.HandleFunc("/transactions", handler.GetTransactions)
	mux.HandleFunc("/summary", handler.GetSummary)
	mux.HandleFunc("/send-email", handler.SendSummaryEmail)

	return &Server{mux: mux, port: port}
}

func (s *Server) Start() {
	fmt.Printf("server running on http://localhost:%s\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.port), corsMiddleware(s.mux)))
}

// Middleware CORS para permitir peticiones desde el frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir todos los orígenes (para desarrollo)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
