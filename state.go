package main

import (
	"fmt"

	"github.com/refugee8888/blogagg/internal/config"
	"github.com/refugee8888/blogagg/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	NameToHandler map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("username is required")
	}
	if err := s.cfg.SetUser(cmd.Arguments[0]); err != nil {
		return fmt.Errorf("there was an error setting up user: %w", err)
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister() error {
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	v, ok := c.NameToHandler[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return v(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if name == "" {
		return fmt.Errorf("command name can't be empty")
	}
	c.NameToHandler[name] = f
	return nil
}
