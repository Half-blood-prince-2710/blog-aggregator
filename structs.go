package main

import (
	"fmt"

	"github.com/half-blood-prince-2710/blog-aggregator/internal/config"
	"github.com/half-blood-prince-2710/blog-aggregator/internal/database"
)


// state and commands
type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	arguments []string
} 

type commands struct{
	cmdHandler map[string]func(*state, command) error
}


func (c *commands) register(name string, f func(*state, command)error) {
		if c.cmdHandler ==nil {
			c.cmdHandler = make(map[string]func(*state, command) error)
		}
		c.cmdHandler[name] = f
}

func (c *commands)  run(s *state,cmd command) error {
	handler , exists := c.cmdHandler[cmd.name]
	if exists {
		return handler(s,cmd)
	}
	return fmt.Errorf("command not found: %s",cmd.name)
}





//RSS 

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}



