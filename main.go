package main

import _ "github.com/lib/pq"

import (
	"github.com/Contary925/gator/internal/config"
	"fmt"
	"os"
	"github.com/Contary925/gator/internal/sql"
	"github.com/Contary925/gator/internal/database"
)

type state struct {
	db  *database.Queries
	Config config.Config
}

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/gator"
	db, err := sql.Open("postgres", dbURL)
	//dbQueries := database.New(db)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	s := state{
		db: db,
		Config: config.Read(),
	} 
	var c commands
	c.handlers = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
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