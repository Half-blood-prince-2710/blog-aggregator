package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/half-blood-prince-2710/blog-aggregator/internal/config"
	"github.com/half-blood-prince-2710/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Print("err", err.Error(), " \n")
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Print("Error: database not connected\nerr: ", err)

	}
	dbQueries := database.New(db)
	
	// Initializing application state
	s := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	//Initializing command registry
	cmds := &commands{}

	//registering our commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg",handlerAgg)
	// cmds.register("addfeed",handlerAddFeed)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds",handlerFeeds)
	cmds.register("follow",handlerFollow)
	cmds.register("following",handlerFollowing)
	cmds.register("unfollow",middlewareLoggedIn(handlerUnfollow))

	//checking if arguments are less than 2
	// fmt.Print(os.Args,"\n")
	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments\n")
		os.Exit(1)
	}
	// extracting command name and arguments
	cmd := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	// running the command
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Print("Error: ", err)
		os.Exit(1)
	}

}




