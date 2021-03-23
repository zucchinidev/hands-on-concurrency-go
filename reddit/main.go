package main

import (
	"fmt"
	"github.com/zucchinidev/hands-on-concurrency-go/reddit/hnClient"
	"github.com/zucchinidev/hands-on-concurrency-go/reddit/redditClient"
	"github.com/zucchinidev/hands-on-concurrency-go/reddit/story"
	"os"
)

func toStdout(c <-chan story.Story)  {
	for {
		s := <-c
		fmt.Printf(stringTmpl(), s.Title, s.URL, s.Author, s.Source)
	}
}

func toFile(c <-chan story.Story, f *os.File)  {
	for {
		s := <-c
		_, _ = fmt.Fprintf(f, stringTmpl(), s.Title, s.URL, s.Author, s.Source)
	}
}

func stringTmpl() string {
	return `%s: %s
by %s on %s

`
}
func main() {
	fromHN := make(chan story.Story, 8)
	fromReddit := make(chan story.Story, 8)
	toFileCh := make(chan story.Story, 8)
	toStdoutCh := make(chan story.Story, 8)


	go hnClient.HackerNewsStories(fromHN)
	go redditClient.RedditStories(fromReddit)

	f, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	go toStdout(toStdoutCh)
	go toFile(toFileCh, f)

	hnOpen, redditOpen := true, true
	var s story.Story

	for hnOpen || redditOpen {
		select {
		case s, hnOpen = <-fromHN:
			if hnOpen {
				toFileCh <- s
				toStdoutCh <- s
			}

		case s, redditOpen = <-fromReddit:
			if redditOpen {
				toFileCh <- s
				toStdoutCh <- s
			}
		}
	}
}


