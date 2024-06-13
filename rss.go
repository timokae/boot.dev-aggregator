package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/timokae/boot.dev-aggregator/internal/database"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func fetchRSS(url string) (RSS, error) {
	log.Printf("Fetching %s", url)

	res, err := http.Get(url)
	if err != nil {
		return RSS{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return RSS{}, fmt.Errorf("got status code %d", res.StatusCode)
	}

	rss := RSS{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return RSS{}, err
	}

	return rss, nil
}

func scrapeFeed(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	log.Printf("Fetching post of %s", feed.Url)

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rss, err := fetchRSS(feed.Url)
	if err != nil {
		log.Println("error fetching rss feed", err)
		return
	}

	for _, post := range rss.Channel.Item {
		description := sql.NullString{}
		if post.Description != "" {
			description.String = post.Description
			description.Valid = true
		}

		t, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			log.Printf("could not parse date %v with err %v", post.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       post.Title,
			Description: description,
			PublishedAt: t,
			Url:         post.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("failed to create post", err)
		}
	}

	log.Printf("Feed %s collected, found %d posts", feed.Name, len(rss.Channel.Item))
}

func scrapeFeeds(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds", err)
			continue
		}

		var waitGroup sync.WaitGroup
		for _, feed := range feeds {
			waitGroup.Add(1)

			go scrapeFeed(db, &waitGroup, feed)
		}

		waitGroup.Wait()
	}
}
