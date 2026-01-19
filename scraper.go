package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yash2895/RSSAggregator/internal/database"
)


func startScrapping(db *database.Queries, concurrency int32, timeBtwRequest time.Duration) {
	log.Printf("Starting scrapping on %d goroutines on %d time interval",concurrency, timeBtwRequest)
	ticker := time.Ticker{}
	for ;; <- ticker.C {
		feeds,err := db.GetFeedsToFetch(context.Background(),concurrency)
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v",err)
			continue;
		}

		wg := &sync.WaitGroup{}
		for _,feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db,wg,feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_,err := db.MarkFeedFetched(context.Background(),feed.ID)
	if err != nil {
		log.Printf("Not able to mark feed as fetched: %v",err)
	}

	rssfeed,err := getFeedFromURL(feed.Url)
	if err != nil {
		log.Printf("%v",err)
		return;
	}
	for _,item := range rssfeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		db.CreatePost(context.Background(),database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
		})

		t,err := time.Parse(time.RFC1123Z,item.PubDate)
		if err != nil {
			log.Printf("could'nt parse date %v with err %v",item.PubDate,err)
			continue
		}

		_,err = db.CreatePost(context.Background(),database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
			PublishedAt: t,
			Url: item.Link,
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(),"duplicate key") {
				continue
			}
			log.Println("Unable to save post ",err)
		}
	}

	log.Printf("Fetched feed %s and found %d items \n",rssfeed.Channel.Title,len(rssfeed.Channel.Item))
}
