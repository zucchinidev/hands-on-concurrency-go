package main

import (
	"fmt"
	"github.com/zucchinidev/hands-on-concurrency-go/searchEngine/story"
	"net/http"
	"strings"
)

var stories []story.Story

func init() {
	stories = append(stories,
		story.Story{
			Title:  "The Raspberry Pi can boot off NVMe SSDs now",
			URL:    "https://www.jeffgeerling.com/blog/2021/raspberry-pi-can-boot-nvme-ssds-now",
			Author: "geerlingguy",
			Source: "HackersNews",
		},

		story.Story{
			Title:  "Practical Cryptography for Developers",
			URL:    "https://cryptobook.nakov.com/",
			Author: "r_singh",
			Source: "HackersNews",
		},

		story.Story{
			Title:  "ARK’s Price Target for Tesla in 2025",
			URL:    "https://ark-invest.com/articles/analyst-research/tesla-price-target-2/",
			Author: "mg",
			Source: "HackersNews",
		},

		story.Story{
			Title:  "Making honey without bees and milk without cows",
			URL:    "https://www.bbc.com/news/business-56154143",
			Author: "elorant",
			Source: "HackersNews",
		},

		story.Story{
			Title:  "Richard Stallman is coming back to the board of the FSF",
			URL:    "http://techrights.org/2021/03/21/richard-stallman-is-coming-back-to-the-board-of-the-free-software-foundation-founded-by-himself-35-years-ago/",
			Author: "wrycoder",
			Source: "HackersNews",
		},
	)
}

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
	}

	w.Write([]byte("<html><body>"))
	ss := searchStories(query)
	if len(ss) == 0 {
		w.Write([]byte(fmt.Sprintf("No results for query '%s'<br>", query)))
	} else {
		for _, s := range ss {
			w.Write([]byte(createLink(s)))
		}
	}
	w.Write([]byte("<a href='../>Back</a>"))
	w.Write([]byte("</html></body>"))
}

func topTen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body>"))
	form := "<form action='search' method='get'>Search: <input type='submit'></form>"
	w.Write([]byte(form))
	for i := 0; i < 10 && len(stories)-1 >= i; i++ {
		s := stories[i]
		w.Write([]byte(createLink(s)))
	}
	w.Write([]byte("</html></body>"))
}

func createLink(s story.Story) string  {
	return fmt.Sprintf("<a href='%s'>'%s'</a><br> %s on %s <br><br>", s.URL, s.Title, s.Author, s.Source)
}

func main() {
	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}
