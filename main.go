package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	cmds.register("addfeed",handlerAddFeed)

	//checking if arguments are less than 2
	// fmt.Print(os.Args,"\n")
	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments\n")
		os.Exit(1)
	}
	cmd := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Print("Error: ", err)
		os.Exit(1)
	}

}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed RSSFeed

	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}
	// fmt.Print(feed.Channel.Description,"\n",feed.Channel.Title,"\n")
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i:= range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
	}

	return  &feed, nil
}


