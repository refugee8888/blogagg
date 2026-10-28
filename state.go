package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

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
	getting, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("unknown username")
		os.Exit(1)
	}
	if err := s.cfg.SetUser(getting.Name); err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("username is required")
	}
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{uuid.New(), time.Now(), time.Now(), cmd.Arguments[0]})
	if err != nil {
		os.Exit(1)
		return err
	}

	if err := s.cfg.SetUser(newUser.Name); err != nil {
		return err
	}
	fmt.Println("user created")
	fmt.Println(newUser)

	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}

	if err := s.cfg.SetUser(""); err != nil {
		return err
	}
	fmt.Println("table has been reset")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	getter, err := s.db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}
	for _, v := range getter {
		if s.cfg.CurrentUserName == v.Name {
			fmt.Println("* " + v.Name + " (current)")
		} else {
			fmt.Println("* " + v.Name)
		}
	}

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
