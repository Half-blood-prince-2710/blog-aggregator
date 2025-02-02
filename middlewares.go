package main

import (
	"context"
	"fmt"

	"github.com/half-blood-prince-2710/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state,cmd command,user database.User)error) func(*state,command)error{
	return func(s *state, c command) error {
		user,err:= s.db.GetUser(context.Background(),s.cfg.Username)
		if err!=nil {
			return fmt.Errorf("error: user not logged in or not found\nerr: %w\n", err)
		}
		return handler(s,c,user)
	}
}