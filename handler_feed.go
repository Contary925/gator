package main

import (
	"fmt"
	"context"
	"github.com/Contary925/gator/internal/feed"
	"github.com/Contary925/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"os"
	"database/sql"
	"strings"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Printf("Not enough arguments!\n")
		os.Exit(1)
	}
	var time_between_reqs time.Duration
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s.db)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	db := s.db
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough arguments provided!\n")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	id := uuid.New()
	time_now := time.Now()
	feed, err := db.CreateFeed(ctx, database.CreateFeedParams{
		ID: id,
		CreatedAt: time_now,
		UpdatedAt: time_now,
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created a new feed:\nName:%s\n", feed.Name)
	_, err = db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: id, 
		CreatedAt: time_now, 
		UpdatedAt: time_now, 
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	db := s.db
	feeds, err := db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Feeds:\n--------------------\n")
	for _, feed := range feeds {
		user, err := db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s\nUrl: %s\nCreated by: %s\n", feed.Name, feed.Url, user.Name)
		fmt.Printf("--------------------\n")
	}
	return nil
}

func scrapeFeeds(db *database.Queries) error {
	ctx := context.Background()
	feed_to_fetch, err := db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	time_now := sql.NullTime{Time: time.Now(), Valid: true}
	err = db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{feed_to_fetch.ID, time_now})
	if err != nil {
		return err
	}
	fetched_feeds, err := feed.FetchFeed(ctx, feed_to_fetch.Url)
	if err != nil {
		return err
	}
	for _, fetched_feed := range fetched_feeds.Channel.Item {
		id := uuid.New()
		_, err := db.CreatePost(ctx, database.CreatePostParams{
			ID: id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: fetched_feed.Title,
			Url: fetched_feed.Link,
			Description: fetched_feed.Description,
			PublishedAt: fetched_feed.PubDate,
			FeedID: feed_to_fetch.ID,
		})
		if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
			fmt.Printf("%v\n", err)
		}
	}
	return nil
}