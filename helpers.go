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


func  scrapeFeeds(s *state)(error) {

	row,err:=s.db.GetNextFeedToFetch(context.Background())
	if err!=nil {
		return fmt.Errorf("error: error fetching next feed\n err: %w\n",err)
	}
	var data database.MarkFeedFetchedParams
	data.ID = row.ID
	data.LastFetchedAt = sql.NullTime{
		Time: time.Now(), // Assign the current time
		Valid: true,         // Set Valid to true to indicate it's not NULL
	}
	data.UpdatedAt = time.Now()
	err = s.db.MarkFeedFetched(context.Background(),data)
	if err!=nil {
		return fmt.Errorf("error: error marking feed\n err: %w\n",err)
	}
	feed,err:=fetchFeed(context.Background(),row.Url)
	if err!=nil {
		return fmt.Errorf("error: error fetching feed\n err: %w\n",err)
	}
	for _,val := range feed.Channel.Item {
		fmt.Print(val.Title,"\n")
	}
return nil
}