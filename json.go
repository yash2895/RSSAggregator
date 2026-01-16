package main

import (
	"encoding/json"
	"log"
	"net/http"
)
func responseWithError(w http.ResponseWriter,status_code int, msg string) {
	if status_code > 499 {
		log.Printf("error with status code 5XX : %s \n", msg)

		type errorMsg struct {
			status_code int
			Error string
		}

		responseWithJSON(w,status_code,errorMsg{
			status_code : status_code, 
			Error :msg,
		})
	}
}
func responseWithJSON(w http.ResponseWriter,status_code int,payload interface{}) {
	dat,err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal JSON :%v", payload)
		w.WriteHeader(500)
	}
	w.WriteHeader(status_code)
	w.Header().Add("Content-type","application/json")
	w.Write(dat)

}
