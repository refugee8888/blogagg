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
	dbQueries := database.New(db)
	if err := c.register("login", handlerLogin); err != nil {
		log.Fatal(err)
	}
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
