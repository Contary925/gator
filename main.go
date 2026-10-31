package main

import (
	"github.com/Contary925/gator/internal/config"
	"fmt"
	"os"
)

type state struct {
	Config config.Config
}

func main() {
	s := state{
		Config: config.Read(),
	} 
	var c commands
	c.handlers = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Not enough arguments!\n")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err := c.run(&s, cmd)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	return
}