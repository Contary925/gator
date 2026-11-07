package main

import _ "github.com/lib/pq"

import (
	"github.com/Contary925/gator/internal/config"
	"fmt"
	"os"
	"database/sql"
	"github.com/Contary925/gator/internal/database"
)

type state struct {
	db  *database.Queries
	Config config.Config
}

func main() {
	dbURL := config.Read().DbURL
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	s := state{
		db: dbQueries,
		Config: config.Read(),
	} 
	var c commands
	c.handlers = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerList)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowsList))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Not enough arguments!\n")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	return
}