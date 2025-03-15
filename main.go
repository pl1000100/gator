package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pl1000100/gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	NameFunc map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.NameFunc[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.NameFunc[cmd.Name]
	if !ok {
		return fmt.Errorf("error: command %s doesn't exists", cmd.Name)
	}
	err := f(s, cmd)
	if err != nil {
		return err
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
	err := s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %s\n", cmd.Arguments[0])
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	s := state{
		Config: &cfg,
	}
	commands := commands{
		NameFunc: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	if len(os.Args) < 3 {
		log.Fatalf("error not enough arguments")
	}
	command := command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}
	commands.run(&s, command)

}
