package servicechecker

import (
	"SlowSloth/common"
	"SlowSloth/pkg/statusprinter"
	"net/http"
	"time"
)

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

		req.Header.Set("User-Agent", common.GetRandomUserAgent())

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
