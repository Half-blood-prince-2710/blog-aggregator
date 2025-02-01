package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/half-blood-prince-2710/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: username required\n")
		
	}
	username := cmd.arguments[0]
	user , err:=s.db.GetUser(context.Background(),username)
	if err!=nil {

	}
	err:=s.cfg.SetUser(username)
	if err!=nil {
		return err
	}
	fmt.Print("User has been Set to ",username)
	return nil
}


