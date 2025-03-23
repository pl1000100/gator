package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pl1000100/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) > 1 {
		fmt.Printf("error: wrong number of arguments passed, expected 0 or 1, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	limit := 2
	if len(cmd.Arguments) == 1 {
		limit, _ = strconv.Atoi(cmd.Arguments[0])
	}

	feed, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, f := range feed {
		fmt.Printf("%s - %s\n%s\n\n", f.PublishedAt.Time, f.Title.String, f.Description.String)
	}

	return nil
}
