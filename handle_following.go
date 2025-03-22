package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("error: wrong number of arguments passed, expected 1, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	feed_id, err := s.db.GetFeedIDByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed_id,
	}
	follow, err := s.db.CreateFeedFollows(context.Background(), params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Feed name: %s\nCurrent user: %s\n", follow.FeedName, currentUser.Name)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Printf("error: wrong number of arguments passed, expected 0, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, f := range feeds {
		println(f.FeedName)
	}
	return nil
}
