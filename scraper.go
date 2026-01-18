package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yash2895/RSSAggregator/internal/database"
)


func startScrapping(db *database.Queries, concurrency int32, timeBtwRequest time.Duration) {
	log.Printf("Starting scrapping on $d goroutines on %d time interval",concurrency, timeBtwRequest)
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
	for _,rssitems := range rssfeed.Channel.Item {
		log.Printf("Found Item: %s \n",rssitems.Title)	
	}
	log.Printf("Fetched feed %s and found %d items \n",rssfeed.Channel.Title,len(rssfeed.Channel.Item))
}
