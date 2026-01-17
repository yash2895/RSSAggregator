package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/yash2895/RSSAggregator/internal/database"
)

type user struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string	`json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser (databaseuser database.User) user {
	return user {
		ID : databaseuser.ID,
		CreatedAt: databaseuser.CreatedAt,
		UpdatedAt: databaseuser.UpdatedAt,
		Name: databaseuser.Name,
		ApiKey: databaseuser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"` 
	UpdatedAt time.Time`json:"update_at"`
	Name      string`json:"name"`
	Url       string`json:"url"`
	UserID    uuid.UUID`json:"user_id"`
}

func databaseFeedToFeed (databasefeed database.Feed) Feed {
	return Feed {
		ID : databasefeed.ID,
		CreatedAt: databasefeed.CreatedAt,
		UpdatedAt: databasefeed.UpdatedAt,
		Name: databasefeed.Name,
		Url: databasefeed.Url,
		UserID: databasefeed.UserID,
	}
}
func databaseFeedsToFeeds (databaseFeeds []database.Feed) []Feed {
	var feeds []Feed
	for _,data_feed := range databaseFeeds {
		feeds = append(feeds,databaseFeedToFeed(data_feed))
	} 
	return feeds;
}


type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"` 
	UpdatedAt time.Time`json:"update_at"`
	UserID    uuid.UUID`json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow (databasefeed database.FeedFollow) FeedFollow {
	return FeedFollow {
		ID : databasefeed.ID,
		CreatedAt: databasefeed.CreatedAt,
		UpdatedAt: databasefeed.UpdatedAt,
		UserID: databasefeed.UserID,
		FeedID: databasefeed.FeedID,
	}
}
func databaseFeedFollowsToFeedFollows (databaseFeeds []database.FeedFollow) []FeedFollow {
	var feeds []FeedFollow
	for _,data_feed := range databaseFeeds {
		feeds = append(feeds,databaseFeedFollowToFeedFollow(data_feed))
	} 
	return feeds;
}
