package main

import "fmt"

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
