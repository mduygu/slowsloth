package servicechecker

import (
	"SlowSloth/pkg/statusprinter"
	"math/rand"
	"net/http"
	"time"
)

// Assuming userAgents slice is available here, either by importing or redefining
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
	"Mozilla/5.0 (iPhone12,1; U; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/15E148 Safari/602.1",
	"Mozilla/5.0 (iPhone12,1; U; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/15E148 Safari/602.1",
	"Mozilla/5.0 (iPhone13,2; U; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/15E148 Safari/602.1",
	"Mozilla/5.0 (iPhone14,3; U; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/19A346 Safari/602.1",
	"Mozilla/5.0 (iPhone14,6; U; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/19E241 Safari/602.1",
	// Add more user agents here
}

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomUserAgent() string {
	return userAgents[globalRand.Intn(len(userAgents))]
}

type ServiceChecker struct {
	url           string
	statusManager *statusprinter.StatusManager
	checkInterval time.Duration
}

func NewServiceChecker(url string, statusManager *statusprinter.StatusManager, interval time.Duration) *ServiceChecker {
	return &ServiceChecker{
		url:           url,
		statusManager: statusManager,
		checkInterval: interval,
	}
}

func (sc *ServiceChecker) CheckServiceAvailability() {
	client := &http.Client{}

	for {
		req, err := http.NewRequest("GET", sc.url, nil)
		if err != nil {
			sc.statusManager.SetServiceAvailable(false)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			sc.statusManager.SetServiceAvailable(false)
		} else {
			sc.statusManager.SetServiceAvailable(resp.StatusCode == http.StatusOK)
			resp.Body.Close()
		}

		time.Sleep(sc.checkInterval)
	}
}
