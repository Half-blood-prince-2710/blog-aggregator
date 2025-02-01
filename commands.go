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


func handlerRegister(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
    return fmt.Errorf("error: username required for registration\n")
}

	user:= database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.arguments[0],
	}
	data, err :=s.db.CreateUser(context.Background(),user)
	if err!=nil{
		return fmt.Errorf("error: user is not registered\n")
	}
	err= s.cfg.SetUser(data.Name)
	if err!=nil {
		return fmt.Errorf("error: Setting user\n")
	}
	fmt.Print("User succesfully created: \n",data,"\n")
	return nil
}