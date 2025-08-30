package main

import (
	"account-summary/src/connections"
	"account-summary/src/handlers"
	"account-summary/src/libs"
	"account-summary/src/repository"
	"account-summary/src/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func main() {
	mongoClient := connections.NewMongoConn(os.Getenv("MONGO_URI"))

	transactionRepository := repository.NewTransactionRepository(mongoClient)
	accountRepository := repository.NewAccountRepository(mongoClient)

	summaryProcessor := libs.NewSummaryProcessor()
	csvReader := libs.NewCsvReader()

	service := services.NewTransactionsService(transactionRepository, accountRepository, csvReader, summaryProcessor)
	handler := handlers.NewTransactionsHandler(service)

	mux := http.NewServeMux()

	// Servir archivos estáticos desde la carpeta assets
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Servir el frontend en la ruta raíz
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
	mux.HandleFunc("/send-email", handler.SendEmail)

	port := os.Getenv("PORT")
	fmt.Printf("server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), corsMiddleware(mux)))

}
