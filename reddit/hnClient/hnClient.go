package hnClient

import (
	"fmt"
	"github.com/caser/gophernews"
	"github.com/zucchinidev/hands-on-concurrency-go/reddit/story"
	"sync"
)

var client *gophernews.Client

func init() {
	// hacker news API allows use it without authentication
	client = gophernews.NewClient()
}

func hnStoryDetail(id int, c chan<- story.Story, wg *sync.WaitGroup) {
	defer wg.Done()
	s, err := client.GetStory(id)
	if err != nil {
		// some ids are comments, we don't care about that
		return
	}
	c <- story.Story{
		Title:  s.Title,
		URL:    s.URL,
		Author: s.By,
		Source: "HackersNews",
	}

}

func HackerNewsStories(c chan <- story.Story) {
	defer close(c)
	changes, err := client.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, id := range changes.Items {
		wg.Add(1)
		go hnStoryDetail(id, c, &wg)
	}
	wg.Wait()
}
