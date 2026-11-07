package feed

import (
	"encoding/xml"
	"net/http"
	"io"
	"context"
	"html"
)

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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var feed RSSFeed
	err = xml.Unmarshal(res, &feed)
	if err != nil {
		return &RSSFeed{}, err		
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	return &feed, nil
}