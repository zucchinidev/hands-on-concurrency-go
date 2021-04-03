package redditClient

import (
	"fmt"
	"github.com/jzelinskie/geddit"
	"github.com/zucchinidev/hands-on-concurrency-go/searchEngine/story"
	"os"
)

var redditSession *geddit.LoginSession

func init() {
	var err error
	redditSession, err = geddit.NewLoginSession("g_d_bot", "K417k4FTua52", "gdAgent v0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RedditStories(c chan<- story.Story) {
	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions) // more recently stories
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range submissions {
		c <- story.Story{
			Title:  s.Title,
			URL:    s.URL,
			Author: s.Author,
			Source: "Reddit /r/programming",
		}
	}
}
