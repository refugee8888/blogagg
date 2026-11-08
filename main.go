package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/refugee8888/blogagg/internal/config"
	"github.com/refugee8888/blogagg/internal/database"
)

func main() {
	var s state
	var c commands
	var cmd command
	args := os.Args
	c.NameToHandler = make(map[string]func(*state, command) error)

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s.cfg = &cfg
	db, err := sql.Open("postgres", s.cfg.DBurl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	if err := c.register("login", handlerLogin); err != nil {
		log.Fatal(err)
	}
	if err := c.register("register", handlerRegister); err != nil {
		log.Fatal(err)
	}
	if err := c.register("reset", handlerResetUsers); err != nil {
		log.Fatal(err)
	}
	if err := c.register("users", handlerUsers); err != nil {
		log.Fatal(err)
	}
	if err := c.register("agg", handlerFeed); err != nil {
		log.Fatal(err)
	}
	if err := c.register("addfeed", middlewareLoggedIn(handlerAddFeed)); err != nil {
		log.Fatal(err)
	}
	if err := c.register("feeds", handlerFeeds); err != nil {
		log.Fatal(err)
	}
	if err := c.register("follow", middlewareLoggedIn(handlerFeedFollow)); err != nil {
		log.Fatal(err)
	}
	if err := c.register("following", middlewareLoggedIn(handlerUserFollowing)); err != nil {
		log.Fatal(err)
	}
	if err := c.register("unfollow", middlewareLoggedIn(handlerFeedUnfollow)); err != nil {
		log.Fatal(err)
	}
	if err := c.register("browse", middlewareLoggedIn(handlerBrowse)); err != nil {
		log.Fatal(err)
	}
	s.db = dbQueries
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}
	cmd.Name = args[1]

	cmd.Arguments = args[2:]
	if err := c.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
