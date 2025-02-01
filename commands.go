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
	
	user , err:=s.db.GetUser(context.Background(),cmd.arguments[0])
	if err!=nil {
		return fmt.Errorf("error: user is not registerd , err: %s",err)
	}
	err =s.cfg.SetUser(user.Name)
	if err!=nil {
		return err
	}
	fmt.Print("User has been Set to ",user.Name)
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

func handlerReset(s *state, cmd command) error{
	err:= s.db.DeleteAllUser(context.Background())
	if err!=nil {
		return fmt.Errorf("error: error deleting users, err: %s",err)
	}

	fmt.Print("Sucessfully deleted all users \n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err!=nil {
		return fmt.Errorf("error: error fetching list of users\nerr: %s",err)
	}
	
	for _,user:=range users {
		if s.cfg.Username == user.Name {
			fmt.Print("* ",user.Name," (current)\n")
			continue
		}
		fmt.Print("* ",user.Name,"\n")
	}

	return  nil
}



// feed handlers 


func handleAgg(s *state, cmd command) error {
	//fetch feed
	feed, err:=fetchFeed(context.Background(),"https://www.wagslane.dev/index.xml")
	if err!=nil{
		fmt.Print("err: ",err,"\n")
		return err
	}
	fmt.Print("feed: \n",feed,"\n")
	return nil
}


