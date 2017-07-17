package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var crawler = struct {
	urls map[string]bool
	mux sync.Mutex
	left int
	done chan bool
}{urls: make(map[string]bool), left: 1, done: make(chan bool)}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {

	crawler.mux.Lock()
	crawler.urls[url] = true
	crawler.mux.Unlock()

	defer func() {
		crawler.mux.Lock()
		crawler.left--
		if crawler.left == 0 {
			crawler.done <- true
		}
		crawler.mux.Unlock()
	}()
	
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		crawler.mux.Lock()
		_, ok := crawler.urls[u]
		if !ok {
			crawler.left++
			go Crawl(u, depth-1, fetcher)
		}
		crawler.mux.Unlock()
	}
	
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
	// Waiting for crawl tasks done
	<-crawler.done
	fmt.Println(crawler)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
