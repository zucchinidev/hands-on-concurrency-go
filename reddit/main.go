package main

import (
	"fmt"
	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
	"os"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {
	// hacker news API allows use it without authentication
	hackerNewsClient = gophernews.NewClient()
	var err error
	redditSession, err = geddit.NewLoginSession("g_d_bot", "K417k4FTua52", "gdAgent v0")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Story struct {
	title  string
	url    string
	author string
	source string
}

func newHackerNewsStories() []Story {
	var stories []Story
	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, id := range changes.Items {
		story, err := hackerNewsClient.GetStory(id)
		if err != nil {
			// some ids are comments, we don't care about that
			continue
		}
		stories = append(stories, Story{
			title:  story.Title,
			url:    story.URL,
			author: story.By,
			source: "HackersNews",
		})
	}
	return stories
}

func newRedditStories() []Story {
	var stories []Story
	sort := geddit.PopularitySort(geddit.NewSubmissions) // more recently stories
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, s := range submissions {
		stories = append(stories, Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		})
	}
	return stories
}

func main() {
	hnStories := newHackerNewsStories()
	redditStories := newRedditStories()

	var stories []Story

	if hnStories != nil {
		stories = append(stories, hnStories...)
	}

	if redditStories != nil {
		stories = append(stories, redditStories...)
	}

	f, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	for _, story := range stories {
		fmt.Fprintf(f, stringTmpl(), story.title, story.url, story.author, story.source)
	}

	for _, story := range stories {
		fmt.Printf(stringTmpl(), story.title, story.url, story.author, story.source)
	}
}

func stringTmpl() string {
	return `%s: %s
by %s on %s

`
}
