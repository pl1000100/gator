package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/gator/internal/database"
)

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("error: too many arguments passed, expected 0")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("error: no arguments passed")
	}
	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("error: too many arguments passed, expected 1")
	}
	name := cmd.Arguments[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("error: no arguments passed")
	}
	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("error: too many arguments passed, expected 1")
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("User was created successfully")
	log.Println(user)

	return nil
}
