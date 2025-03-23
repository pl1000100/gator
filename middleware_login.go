package main

import (
	"context"
	"fmt"

	"github.com/pl1000100/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return handler(s, cmd, currentUser)
	}
}
