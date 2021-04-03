package main

import (
	"fmt"
	"github.com/zucchinidev/hands-on-concurrency-go/searchEngine/hnClient"
	"github.com/zucchinidev/hands-on-concurrency-go/searchEngine/redditClient"
	"github.com/zucchinidev/hands-on-concurrency-go/searchEngine/story"
	"net/http"
	"strings"
	"time"
)

var stories []story.Story

func searchStories(query string) []story.Story {
	var ff []story.Story

	for _, s := range stories {
		if strings.Contains(strings.ToUpper(s.Title), strings.ToUpper(query)) {
			ff = append(ff, s)
		}
	}
	return ff
}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	// searchhn.com/search?q=making
	if query == "" {
		http.Error(w, "Search parameter q is required to search.", http.StatusNotAcceptable)
		return
	}

	_, _ = w.Write([]byte("<html><body>"))
	ss := searchStories(query)
	if len(ss) == 0 {
		_, _ = w.Write([]byte(fmt.Sprintf("No results for query '%s'<br>", query)))
	} else {
		for _, s := range ss {
			_, _ = w.Write([]byte(createLink(s)))
		}
	}
	_, _ = w.Write([]byte("<a href='../'>Back</a>"))
	_, _ = w.Write([]byte("</html></body>"))
}

func topTen(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("<html><body>"))
	form := "<form action='search' method='get'>Search: <input type='text' name='q'/> <input type='submit'/></form>"
	_, _ = w.Write([]byte(form))
	for i := 0; i < 10 && len(stories)-1 >= i; i++ {
		s := stories[i]
		_, _ = w.Write([]byte(createLink(s)))
	}
	_, _ = w.Write([]byte("</html></body>"))
}

func createLink(s story.Story) string {
	return fmt.Sprintf("<a href='%s'>'%s'</a><br> %s on %s <br><br>", s.URL, s.Title, s.Author, s.Source)
}

func toList(c <-chan story.Story) {
	for {
		s := <-c
		if !isRepeated(s) {
			stories = append(stories, s)
		}
	}
}

func isRepeated(s story.Story) bool {
	for _, st := range stories {
		if st.Title == s.Title {
			return true
		}
	}
	return false
}

func main() {

	toListCh := make(chan story.Story, 8)
	go toList(toListCh)

	go func() {
		for {
			fromHN := make(chan story.Story, 8)
			fromReddit := make(chan story.Story, 8)
			go hnClient.HackerNewsStories(fromHN)
			go redditClient.RedditStories(fromReddit)

			hnOpen, redditOpen := true, true
			var s story.Story

			for hnOpen || redditOpen {
				select {
				case s, hnOpen = <-fromHN:
					if hnOpen {
						toListCh <- s
					}

				case s, redditOpen = <-fromReddit:
					if redditOpen {
						toListCh <- s
					}
				}
			}

			// Now we'll report that we're finished and
			// wait a bit before getting new stories.
			fmt.Println("Done fetching new stories.")
			time.Sleep(10 * time.Second)
		}
	}()

	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}
