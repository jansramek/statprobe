package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 1, "concurrency level")

	var statusCode int
	flag.IntVar(&statusCode, "s", 200, "status code that will be matched")

	var debug bool
	flag.BoolVar(&debug, "d", false, "toggle debug mode - prints non-matched status codes")

	flag.Parse()

	urlsChannel := make(chan string)
	outputChannel := make(chan string)

	// Start status code check routines
	var jobWG sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		jobWG.Add(1)
		go checkStatus(urlsChannel, outputChannel, &jobWG, statusCode, debug)
	}

	// Start outputChannel worker for printing results
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		for o := range outputChannel {
			fmt.Println(o)
		}
		outputWG.Done()
	}()

	// Close the outputChannel channel when co-routines are done
	go func() {
		jobWG.Wait()
		close(outputChannel)
	}()

	// Read input data from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urlsChannel <- scanner.Text()
	}
	close(urlsChannel)

	if scanner.Err() != nil {
		println("Error when reading an input URLs")
	}

	outputWG.Wait()
}

func checkStatus(urls chan string, o chan string, group *sync.WaitGroup, statusCode int, debug bool)  {
	for url := range urls {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Get(url)
		if (resp.StatusCode == statusCode) {
			o <- url
		} else if (debug) {
			o <- "[" + fmt.Sprintf("%d", resp.StatusCode) + "] " + url
		}
		if (err != nil) {
			o <- "[Error] " + err.Error()
		}
	}

	group.Done()
}
