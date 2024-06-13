package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/timokae/boot.dev-aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.Apikey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	var lastFetchedAt *time.Time
	if feed.LastFetchedAt.Valid {
		lastFetchedAt = &feed.LastFetchedAt.Time
	}

	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: lastFetchedAt,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	feedsToReturn := make([]Feed, 0)

	for _, feed := range feeds {
		feedsToReturn = append(feedsToReturn, databaseFeedToFeed(feed))
	}

	return feedsToReturn
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(feed database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feed.ID,
		FeedID:    feed.FeedID,
		UserID:    feed.UserID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	}
}

func databaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollow {
	feedFollowsToReturn := make([]FeedFollow, 0)

	for _, feedFollow := range feedFollows {
		feedFollowsToReturn = append(feedFollowsToReturn, databaseFeedFollowToFeedFollow(feedFollow))
	}

	return feedFollowsToReturn
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	FeedId      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(post database.Post) Post {
	var description *string
	if post.Description.Valid {
		description = &post.Description.String
	}

	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
		Title:       post.Title,
		Description: description,
		URL:         post.Url,
		FeedId:      post.FeedID,
	}
}

func databasePostsToPosts(posts []database.Post) []Post {
	postsToReturn := make([]Post, 0)

	for _, post := range posts {
		postsToReturn = append(postsToReturn, databasePostToPost(post))
	}

	return postsToReturn
}
