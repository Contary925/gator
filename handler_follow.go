package main

import (
	"fmt"
	"context"
	"github.com/Contary925/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"os"
	"strconv"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	db := s.db
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments provided!\n")
	}
	url := cmd.args[0]
	id := uuid.New()
	time_now := time.Now()
	feed, err := db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
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
	fmt.Printf("The user \"%s\" is now following the feed \"%s\"\n", user.Name, feed.Name)
	return nil
}

func handlerFollowsList(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	db := s.db
	res, err := db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("List of news you are following:\n")
	for index, feed_follow := range res {
		fmt.Printf("%v. %s\n", index+1, feed_follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		fmt.Printf("Not enough arguments!")
		os.Exit(1)
	}
	url := cmd.args[0]
	ctx := context.Background()
	db := s.db
	feed, err := db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	err = db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{user.ID, feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully unfollowed.\n")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		parsed_limit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			limit = 2
		} else {
			limit = int32(parsed_limit)
		}
	}
	ctx := context.Background()
	db := s.db
	posts, err := db.GetPostsByUserID(ctx, database.GetPostsByUserIDParams{user.ID, limit})
	if err != nil {
		return err
	}
	fmt.Printf("====================\nSaved posts:\n")
	for idx, post := range posts {
		fmt.Printf("%v. %s\n", idx+1, post.Title)
	}
	return nil
}