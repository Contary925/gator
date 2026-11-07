package main

import (
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
	"os"
	"github.com/Contary925/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments provided!\n")
	}
	ctx := context.Background()
	username := cmd.args[0]
	db := s.db
	if user, _ := db.GetUserByName(ctx, username); user.Name == "" {
		fmt.Printf("A user with this username does not exist!\n")
		os.Exit(1)
	}
	s.Config.SetUser(username)
	fmt.Printf("The user has been set to \"%s\".\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments provided!\n")
	}
	ctx := context.Background()
	username := cmd.args[0]
	id := uuid.New()
	time_now := time.Now()
	db := s.db
	if user, _ := db.GetUserByName(ctx, username); user.Name != "" {
		fmt.Printf("A user with this username already exists!\n")
		os.Exit(1)
	}
	_, err := db.CreateUser(ctx, database.CreateUserParams{
		ID: id,
		CreatedAt: time_now,
		UpdatedAt: time_now,
		Name: username,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created a user with username %s.\n", username)
	s.Config.SetUser(username)
	fmt.Printf(s.Config.Current_user_name+"\n")
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	db := s.db
	err := db.Reset(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Successfully reset the database.\n")
	return nil
}

func handlerList(s *state, cmd command) error {
	ctx := context.Background()
	db := s.db
	users, err := db.GetUsers(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.Name == s.Config.Current_user_name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	result_handler := func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.Config.Current_user_name)
		if err != nil {
			return err
		}
		handler(s, cmd, user)
		return err
	}
	return result_handler
}