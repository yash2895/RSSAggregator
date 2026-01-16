package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Home(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w,"Homepage")
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Error getting port")
	}
	log.Println("PORT :",portString)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz",handler_readiness)
	mux.HandleFunc("GET /error",handler_error)

	log.Printf("Server started on port %v \n", portString)

	err := http.ListenAndServe(":"+portString,mux);
	if (err != nil) {
		log.Fatalf("Error starting server on port %v",portString)
	}


}


