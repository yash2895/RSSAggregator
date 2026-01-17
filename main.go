package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	
    _ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/yash2895/RSSAggregator/internal/database"
)
type apiConfig struct {
	DB *database.Queries
}

func Home(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w,"Homepage")
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Error getting port")
	}
	dbURL := os.Getenv("POSTGRES_URL")
	if  dbURL == "" {
		log.Fatal("Error getting postgres URL")
	}
	log.Println("PORT :",portString)

	conn,err := sql.Open("postgres",dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	
	apiCfg := apiConfig{
		DB : database.New(conn),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz",apiCfg.handler_readiness)
	mux.HandleFunc("GET /error",apiCfg.handler_error)
	mux.HandleFunc("POST /users",apiCfg.handler_users)
	mux.HandleFunc("GET /users",apiCfg.handlerGetUsers)
	mux.HandleFunc("POST /feeds",apiCfg.middlewareAuth(apiCfg.handlerFeed))
	mux.HandleFunc("POST /feed_follow",apiCfg.middlewareAuth(apiCfg.handlerFeedFollows))
	mux.HandleFunc("GET /feed_follow",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))

	log.Printf("Server started on port %v \n", portString)

	err = http.ListenAndServe(":"+portString,mux);
	if (err != nil) {
		log.Fatalf("Error starting server on port %v",portString)
	}


}


