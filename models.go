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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
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
