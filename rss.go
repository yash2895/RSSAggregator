package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title        string  `xml:"title"`
		Link         string  `xml:"link"`
		Description  string  `xml:"description"`
		Language     string  `xml:"language"`
		Item         []RSSItem`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
		Title        string  `xml:"title"`
		Link         string  `xml:"link"`
		Description  string  `xml:"description"`
		pubDate      string  `xml:"pubdate"`
}	

func getFeedFromURL (url string) (RSSFeed, error) {
	client := http.Client {
		Timeout: 10 * time.Second,
	}
	resp,err := client.Get(url)
	if err != nil {
		return RSSFeed{},fmt.Errorf("Not able to feed fetch %s", url)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{} , fmt.Errorf("Error reading the fetched feed: %v", err);
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data,&rssFeed);
	if err != nil {
		return RSSFeed{},err;
	}
	return rssFeed,nil;
}
