package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("error: too many arguments passed, expected 0")
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("error: couldn't delete users")
		return err
	}
	fmt.Println("Deleted all users successfully")
	return nil
}
