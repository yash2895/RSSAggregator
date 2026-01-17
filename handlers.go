package main

import (
	"encoding/json"

	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yash2895/RSSAggregator/internal/auth"
	"github.com/yash2895/RSSAggregator/internal/database"
)

func (api *apiConfig) handler_readiness(w http.ResponseWriter,r *http.Request) {
	responseWithJSON(w,200,struct{}{})
}

func(api *apiConfig) handler_error(w http.ResponseWriter,r *http.Request) {
	responseWithError(w,500,"something went wrong")
}

func (api *apiConfig) handler_users(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	body := parameters{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responseWithError(w,400,fmt.Sprintln("Bad Body"))
		return
	}

	user,err := api.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID : uuid.New(),
		CreatedAt : time.Now().UTC(), 
		UpdatedAt : time.Now().UTC(),  
		Name : body.Name,         
	})

	if err != nil {
		responseWithError(w,500,"Unable to create a user")
	}
	
	responseWithJSON(w,200,databaseUserToUser(user))
}

func (api *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	apiKey,err := auth.GetAPIKey(r.Header)

	if err != nil {
		responseWithError(w,403,err.Error())
		return
	}

	user, err := api.DB.GetUserByAPIKey(r.Context(),apiKey)
	if err!=nil {
		responseWithError(w,500,err.Error())
	}

	responseWithJSON(w,201,databaseUserToUser(user))
}

func (api *apiConfig) handlerFeed(w http.ResponseWriter,r *http.Request, user database.User) {
	
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	body := parameters{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responseWithError(w,400,fmt.Sprintln("Bad Body"))
		return
	}

	feed,err := api.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID : uuid.New(),
		CreatedAt : time.Now().UTC(), 
		UpdatedAt : time.Now().UTC(),  
		Name : body.Name,         
		Url: body.Url,
		UserID: user.ID,
	})

	if err != nil {
		responseWithError(w,500,"Unable to create a feed")
	}
	
	responseWithJSON(w,200,databaseFeedToFeed(feed))
}

func (api *apiConfig) handlerGetFeeds(w http.ResponseWriter,r *http.Request) {
	feeds, err := api.DB.GetFeeds(r.Context())

	if err != nil {
		responseWithError(w,500,fmt.Sprintf("Error fetching all feeds: %v", err))
	}
	responseWithJSON(w,200,databaseFeedsToFeeds(feeds))
}

func (api *apiConfig) handlerFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()

	type parameters struct {
		Feed_id uuid.UUID `json:"feed_id"`
	}

	body := parameters{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Feed_id == uuid.Nil{
		responseWithError(w,400,fmt.Sprintln("Bad Body"))
		return
	}

	feed,err := api.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID : uuid.New(),
		CreatedAt : time.Now().UTC(), 
		UpdatedAt : time.Now().UTC(),  
		UserID: user.ID,
		FeedID: body.Feed_id,
	})

	if err != nil {
		responseWithError(w,500,"Unable to create a feed")
	}
	
	responseWithJSON(w,200,databaseFeedFollowToFeedFollow(feed))

}

func (api *apiConfig) handlerGetFeedFollows(w http.ResponseWriter,r *http.Request, user database.User) {
	feeds ,err := api.DB.GetFeedFollows(r.Context(),user.ID)

	if err != nil {
		responseWithError(w,400,fmt.Sprintf("Cannot feeds with the follows: %v",err))
		return
	}

	responseWithJSON(w,200,databaseFeedFollowsToFeedFollows(feeds))
}
