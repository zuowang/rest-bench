package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func doRequests(url string, request string, resultCh chan<- bool) {
	for {
		resp, err := http.Post(url, "application/json", strings.NewReader(request))
		if err != nil {
			resultCh <- false
			fmt.Printf("error: %v\n", err)
		} else if resp.StatusCode != 200 {
			resultCh <- false
			fmt.Printf("status: %v\n", resp.Status)
		} else {
			resultCh <- true
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

func main() {
	var url = flag.String("url", "", "Target URL")
	var parallel = flag.Int("parallel", 1, "parallel closed-loop operations")
	var request = flag.String("request", "", "Request content")
	flag.Parse()

	reportPeriod := 1 * time.Second

	results := make(chan bool, *parallel*10)
	for p := 0; p < *parallel; p++ {
		go doRequests(*url, *request, results)
	}

	start := time.Now()
	success := 0
	fail := 0
	for {
		res := <-results
		if res {
			success++
		} else {
			fail++
		}

		elapsed := time.Since(start)
		if elapsed >= reportPeriod {
			fmt.Printf("success: %v fail: %v\n", float64(success)/elapsed.Seconds(), float64(fail)/elapsed.Seconds())
			success = 0
			fail = 0
			start = time.Now()
		}
	}
}
