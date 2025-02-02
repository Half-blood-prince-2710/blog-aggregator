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

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error: user is not registerd , err: %s", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Print("User has been Set to ", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: username required for registration\n")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	data, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error: user is not registered\n")
	}
	err = s.cfg.SetUser(data.Name)
	if err != nil {
		return fmt.Errorf("error: Setting user\n")
	}
	fmt.Print("User succesfully created: \n", data, "\n")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUser(context.Background())
	if err != nil {
		return fmt.Errorf("error: error deleting users, err: %s", err)
	}

	fmt.Print("Sucessfully deleted all users \n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: error fetching list of users\nerr: %s", err)
	}

	for _, user := range users {
		if s.cfg.Username == user.Name {
			fmt.Print("* ", user.Name, " (current)\n")
			continue
		}
		fmt.Print("* ", user.Name, "\n")
	}

	return nil
}

// feed handlers

func handlerAgg(s *state, cmd command) error {
	//fetch feed
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Print("err: ", err, "\n")
		return err
	}
	fmt.Print("feed: \n", feed, "\n")
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("error: not enough arguments\n ")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("error: error fetching user\nerr: %w\n", err)
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	var feed database.CreateFeedParams
	feed.Name = name
	feed.Url = url
	feed.UserID = user.ID
	feeds, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("error: error creating feed\nerr: %w\n", err)
	}
	

	var feed_follow database.CreateFeedFollowParams

	feed_follow.UserID = user.ID
	feed_follow.FeedID = feeds.ID
	fmt.Print("user_id",user.ID,"feed_id",feeds.ID,"\n")
	rows,err :=	s.db.CreateFeedFollow(context.Background(),feed_follow)
	if err != nil {
		return fmt.Errorf("error: error creating feed follow record\nerr: %w\n", err)
	}
	fmt.Print("rows: ",rows,"\n")
	fmt.Print("feed: \n", feeds,"\n")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error: error fetching feeds\nerr: %w\n", err)
	}

	for _, feed := range feeds {
		name, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error: error fetching user\nerr: %w\n", err)
		}
		fmt.Print("Name: ", feed.Name, "\nUrl: ", feed.Url, "\nuser", name, "\n")
	}
	return nil
}

// feed follows commands

func handlerFollow(s *state, cmd command) error {
	url := cmd.arguments[0]
	user, err := s.db.GetUser(context.Background(), s.cfg.Username)

	if err != nil {
		return fmt.Errorf("error: error fetching user\nerr: %w\n", err)
	}
	feed_id , err := s.db.GetFeedByUrl(context.Background(),url)
	if err != nil {
		return fmt.Errorf("error: error fetching feed\nerr: %w\n", err)
	}
	var feed_follow database.CreateFeedFollowParams

	feed_follow.UserID = user.ID
	feed_follow.FeedID = feed_id
	rows,err:=	s.db.CreateFeedFollow(context.Background(),feed_follow)
	if err != nil {
		return fmt.Errorf("error: error creating feed follow record\nerr: %w\n", err)
	}
	fmt.Print("succesfully follow url: ",url," feed_name: ",rows.FeedName," user_name: ",rows.UserName,"\n")
	return nil
}



func handlerFollowing(s *state, cmd command) error {
    user, err := s.db.GetUser(context.Background(), s.cfg.Username)
    if err != nil {
        return fmt.Errorf("error: error fetching user\nerr: %w", err)
    }

    // Fetch feeds the user follows
    feeds, err := s.db.GetFeedFollow(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf("error: error fetching followed feeds\nerr: %w", err)
    }

    if len(feeds) == 0 {
        fmt.Println("You are not following any feeds.")
        return nil
    }

    // Print the followed feeds
    for _, feed := range feeds {
        fmt.Println(feed.Name)
    }
    return nil
}
