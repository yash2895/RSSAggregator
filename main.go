package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yash2895/RSSAggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Error getting port")
	}
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatal("Error getting postgres URL")
	}
	log.Println("PORT :", portString)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScrapping(apiCfg.DB, 10, time.Minute)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", apiCfg.handler_readiness)
	mux.HandleFunc("GET /error", apiCfg.handler_error)
	mux.HandleFunc("POST /users", apiCfg.handler_users)
	mux.HandleFunc("GET /users", apiCfg.handlerGetUsers)
	mux.HandleFunc("POST /feeds", apiCfg.middlewareAuth(apiCfg.handlerFeed))
	mux.HandleFunc("POST /feed_follow", apiCfg.middlewareAuth(apiCfg.handlerFeedFollows))
	mux.HandleFunc("GET /feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	mux.HandleFunc("DELETE /feed_follow/{id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))
	mux.HandleFunc("GET /post", apiCfg.middlewareAuth(apiCfg.handlerGetPosts))

	server := http.Server{
		Handler: mux,
		Addr:    ":" + portString,
	}

	// gracefull shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {

		log.Printf("Server started on port %v \n", portString)

		err = server.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed,err){
			log.Fatalf("Error starting server on port %v", portString)
		}
	}()

	<- done
	log.Println("Shutting down server")
	timer,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err = server.Shutdown(timer)
	if err != nil {
		log.Fatalf("Error shutting down server %v",err)
	}
}
