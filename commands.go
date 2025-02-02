package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("error: not enough arguments\n")
	}

	timeBetweenReqs := cmd.arguments[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	// Print the interval message
	fmt.Printf("Collecting feeds every %v\n", duration)

	// Set up a ticker to scrape feeds at regular intervals
	ticker := time.NewTicker(duration)

	// Create a signal channel to handle graceful shutdown
	signalChan := make(chan os.Signal, 1)
	// Notify the signal channel for SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Ensure the loop runs immediately before waiting for the first ticker
	go func() {
		for {
			select {
			case <-ticker.C:
				// Scrape feeds periodically
				err := scrapeFeeds(s)
				if err != nil {
					log.Printf("Error scraping feeds: %v", err)
				}
			case <-signalChan:
				// Stop the ticker and exit the loop when receiving a signal
				fmt.Println("\nShutting down gracefully...")
				ticker.Stop()
				return
			}
		}
	}()

	// Block the main goroutine to keep the process running
	<-signalChan
	return nil
}


func handlerAddFeed(s *state, cmd command,user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("error: not enough arguments\n ")
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


func handlerUnfollow(s *state, cmd command, user database.User)error {
	url := cmd.arguments[0]
	feed_id,err:=s.db.GetFeedByUrl(context.Background(),url)
	if err!=nil {
		return fmt.Errorf("error: error fecthing feed it\nerr: %w", err)
	}
	var data database.DeleteFeedFollowParams
	data.UserID= user.ID
	data.FeedID = feed_id
	err = s.db.DeleteFeedFollow(context.Background(),data)
	if err!=nil {
		return fmt.Errorf("error: error unfollowing feed\nerr: %w", err)
	}
	fmt.Print("Succesfully unfollow\n")
	return nil
}


func handlerBrowse(s *state, cmd command) error {
	// Set the default limit to 2 if not provided
	limit := 2
	if len(cmd.arguments) > 0 {
		// If a limit argument is provided, parse it
		parsedLimit, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}
	 user, err := s.db.GetUser(context.Background(), s.cfg.Username)
    if err != nil {
        return fmt.Errorf("error: error fetching user\nerr: %w", err)
    }


	// Fetch the posts from the database
	data:= database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), data)
	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	// Print the posts
	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\nDescription: %s\nPublished At: %s\n\n",
			post.Title, post.Url, post.Description, post.PublishedAt)
	}

	return nil
}
