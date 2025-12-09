package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/fatih-sonay/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %v duration", concurency, timeBetweenRequest)

	tickeer := time.NewTicker(timeBetweenRequest)

	for ; ; <-tickeer.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurency))

		if err != nil {
			log.Printf("error fetching feeds: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, feed, wg)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("error fetching feed %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Printf("Found post %s on feed %s", item.Title, feed.Name)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
