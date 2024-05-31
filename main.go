package main

import (
	"SlowSloth/common"
	"SlowSloth/pkg/servicechecker"
	slowrequest "SlowSloth/pkg/slowrequester"
	"SlowSloth/pkg/statusprinter"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
}

type ChartData struct {
	Timestamps []string `json:"timestamps"`
	Values     []int    `json:"values"`
}

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
		common.LogError(fmt.Errorf("Validation error: %w", err))
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

	var connectionsData []DataPoint
	var ramUsageData []DataPoint
	var availabilityData []DataPoint

	go func() {
		defer printWg.Done()
		for {
			select {
			case <-done:
				return
			default:
				timestamp := time.Now()
				connectionsData = append(connectionsData, DataPoint{Timestamp: timestamp, Value: int(statusManager.ActiveConnections())})
				ramUsageData = append(ramUsageData, DataPoint{Timestamp: timestamp, Value: int(statusManager.TotalRAMUsage())})
				availabilityValue := 0
				if statusManager.IsServiceAvailable() {
					availabilityValue = 1
				}
				availabilityData = append(availabilityData, DataPoint{Timestamp: timestamp, Value: availabilityValue})

				// Terminal imlecini önceki 5 satırın başlangıcına taşı ve satırları temizle
				fmt.Print("\033[H\033[2J")
				fmt.Printf("Total active connections: %d\n", statusManager.ActiveConnections())
				fmt.Printf("Service availability: %s%s\n", statusManager.SetServiceColor(statusManager.IsServiceAvailable()), statusManager.ServiceAvailability())
				fmt.Printf(Reset)
				fmt.Printf("Total RAM usage: %d MB\n", statusManager.TotalRAMUsage())
				fmt.Printf("Total bandwidth usage: %d bytes\n", statusManager.TotalBandwidth())
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
	fmt.Printf("Final Status - Total active connections: %d, Service availability: %t, Total RAM usage %d MB, Total bandwidth usage %d bytes\n",
		statusManager.ActiveConnections(),
		statusManager.IsServiceAvailable(),
		statusManager.TotalRAMUsage(),
		statusManager.TotalBandwidth())

	fmt.Println("All requests completed.")

	// Save the data to an HTML file
	saveResults(connectionsData, ramUsageData, availabilityData)
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

func saveResults(connectionsData, ramUsageData, availabilityData []DataPoint) {
	tmpl, err := os.ReadFile("html/template.html")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	connectionsChartData := formatChartData(connectionsData)
	ramUsageChartData := formatChartData(ramUsageData)
	availabilityChartData := formatChartData(availabilityData)

	htmlContent := strings.Replace(string(tmpl), "{{CONNECTIONS_DATA}}", connectionsChartData, 1)
	htmlContent = strings.Replace(htmlContent, "{{RAM_USAGE_DATA}}", ramUsageChartData, 1)
	htmlContent = strings.Replace(htmlContent, "{{AVAILABILITY_DATA}}", availabilityChartData, 1)

	currentTime := time.Now().Format("2006-01-02-15-04-05")
	fileName := fmt.Sprintf("results-%s.html", currentTime)
	err = os.WriteFile(fileName, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("Error writing results file:", err)
	}
}

func formatChartData(data []DataPoint) string {
	var timestamps []string
	var values []int
	for _, dp := range data {
		timestamps = append(timestamps, dp.Timestamp.Format(time.RFC3339))
		values = append(values, dp.Value)
	}

	chartData := ChartData{
		Timestamps: timestamps,
		Values:     values,
	}

	chartDataJSON, err := json.Marshal(chartData)
	if err != nil {
		fmt.Println("Error marshalling chart data:", err)
		return "{}"
	}

	return string(chartDataJSON)
}
