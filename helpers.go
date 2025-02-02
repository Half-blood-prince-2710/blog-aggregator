package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/half-blood-prince-2710/blog-aggregator/internal/database"
	"github.com/lib/pq"
)

// function for fetching feeds
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


// scrape feeds function


func scrapeFeeds(s *state) error {
	// Get the next feed to fetch from the database
	row, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error: fetching next feed\n err: %w\n", err)
	}

	// Mark the feed as fetched
	data := database.MarkFeedFetchedParams{
		ID: row.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true, // Mark as valid to indicate it's not NULL
		},
		UpdatedAt: time.Now(),
	}

	err = s.db.MarkFeedFetched(context.Background(), data)
	if err != nil {
		return fmt.Errorf("error: marking feed as fetched\n err: %w\n", err)
	}

	// Fetch the RSS feed data
	feed, err := fetchFeed(context.Background(), row.Url)
	if err != nil {
		return fmt.Errorf("error: fetching feed\n err: %w\n", err)
	}

	// Iterate through feed items and save them to the database
	for _, item := range feed.Channel.Item {
		// Parse the "published_at" timestamp correctly
		publishedAt, err := parsePublishedTime(item.PubDate)
		if err != nil {
			fmt.Printf("warning: failed to parse published_at for post %s\n", item.Title)
			continue // Skip this post if the timestamp can't be parsed
		}

		// Prepare the post data for insertion
		postData := database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      row.ID,
		}

		// Insert post into the database
		_, err = s.db.CreatePost(context.Background(), postData)
		if err != nil {
			// Ignore duplicate URL errors, log others
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				fmt.Printf("info: duplicate post skipped: %s\n", item.Title)
			} else {
				fmt.Printf("error: failed to save post %s\n err: %v\n", item.Title, err)
			}
			continue
		}

		fmt.Printf("success: saved post %s\n", item.Title)
	}

	return nil
}

// Parses different date formats from RSS feeds
func parsePublishedTime(pubDate string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,    // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,   // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,    // "2006-01-02T15:04:05Z07:00"
		"Mon, 02 Jan 2006 15:04:05 -0700", // Some RSS variations
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, pubDate)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse published date: %s", pubDate)
}
