package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments provided!\n")
	}
	username := cmd.args[0]
	s.Config.SetUser(username)
	fmt.Printf("The user has been set to \"%s\".\n", username)
	return nil
}