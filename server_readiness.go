package main

import "net/http"

func handler_readiness(w http.ResponseWriter,r *http.Request) {
	responseWithJSON(w,200,struct{}{})
}

func handler_error(w http.ResponseWriter,r *http.Request) {
	responseWithError(w,500,"something went wrong")
}
