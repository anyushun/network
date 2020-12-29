package clientsets

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"net-tools/pkg/counter"
)

type httpClientSetConfig struct {
	parallel        int
	requestInterval time.Duration
	requestTimeout  time.Duration
	requestURL      string
	requestMethod   string
	requestBody     string
}

var hcsc = &httpClientSetConfig{}

func init() {
	flag.IntVar(&hcsc.parallel, "http-parallel", 1, "http connect parallel")
	flag.DurationVar(&hcsc.requestInterval, "http-request-interval", time.Second, "http request interval")
	flag.DurationVar(&hcsc.requestTimeout, "http-request-timeout", time.Second, "http request timeout")
	flag.StringVar(&hcsc.requestURL, "http-request-url", "http://127.0.0.1:8080/", "http request url")
	flag.StringVar(&hcsc.requestMethod, "http-request-method", http.MethodPost, "http request method")
	flag.StringVar(&hcsc.requestBody, "http-request-body", "hello, http server", "http request body")
}

// HTTPClientSet impl
type HTTPClientSet struct {
	connections  *counter.Counter
	requests     *counter.Counter
	sendBytes    *counter.Counter
	receiveBytes *counter.Counter
	ctx          context.Context
	cancel       context.CancelFunc
	once         *sync.Once
}

// Protocol of HTTPClientSet
func (hcs *HTTPClientSet) Protocol() string {
	return "http"
}

func (hcs *HTTPClientSet) handleClient() {
	client := &http.Client{
		Timeout: hcsc.requestTimeout,
	}
	hcs.connections.Add(1)

	for {
		select {
		case <-hcs.ctx.Done():
			return
		default:
			request, err := http.NewRequest(hcsc.requestMethod, hcsc.requestURL, strings.NewReader(hcsc.requestBody))
			if err != nil {
				log.Panicln(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				log.Println(err)
				hcs.connections.Sub(1)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			resp.Body.Close()
			hcs.requests.Add(1)
			hcs.sendBytes.Add(float64(len(hcsc.requestBody)))
			hcs.receiveBytes.Add(float64(len(body)))
		}
	}
}

// Start HTTPClientSet
func (hcs *HTTPClientSet) Start() error {
	hcs.connections = &counter.Counter{}
	hcs.requests = &counter.Counter{}
	hcs.sendBytes = &counter.Counter{}
	hcs.receiveBytes = &counter.Counter{}
	hcs.once = &sync.Once{}
	hcs.ctx, hcs.cancel = context.WithCancel(context.Background())

	for i := 0; i < hcsc.parallel; i++ {
		go hcs.handleClient()
	}

	return nil
}

// Stop HTTPClientSet
func (hcs *HTTPClientSet) Stop() error {
	hcs.cancel()
	return nil
}

// Report HTTPClientSet
func (hcs *HTTPClientSet) Report() string {
	r := fmt.Sprintf("| %-15.0f | %-15.0f | %-15.0f | %-15.0f |", hcs.connections.Get(), hcs.requests.Get(), hcs.sendBytes.Get(), hcs.receiveBytes.Get())
	hcs.once.Do(func() {
		r = fmt.Sprintf("| %-15s | %-15s | %-15s | %-15s |", "connections", "requests", "send bytes", "receive bytes")
	})
	return r
}
