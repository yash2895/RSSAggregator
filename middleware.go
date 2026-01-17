package main

import (
	"net/http"

	"github.com/yash2895/RSSAggregator/internal/auth"
	"github.com/yash2895/RSSAggregator/internal/database"
)

type authedHandler func(w http.ResponseWriter,r *http.Request,user database.User)

func (api *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		 
	apiKey,err := auth.GetAPIKey(r.Header)

	if err != nil {
		responseWithError(w,403,err.Error())
		return
	}

	user, err := api.DB.GetUserByAPIKey(r.Context(),apiKey)
	if err!=nil {
		responseWithError(w,500,err.Error())
	}


	handler(w,r,user)
	}
}
