package protocols

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net-tools/pkg/counter"
	"net/http"
	"sync"
)

type httpConfig struct {
	addr               string
	responseStatusCode int
	responseBody       string
}

var hc = &httpConfig{}

func init() {
	flag.StringVar(&hc.addr, "http-address", "0.0.0.0:8080", "address http server listen on")
	flag.IntVar(&hc.responseStatusCode, "http-code", http.StatusOK, "http response code")
	flag.StringVar(&hc.responseBody, "http-body", "hello world", "http response body")
}

// HTTPServer struct
type HTTPServer struct {
	addr         string
	requests     *counter.Counter
	sendBytes    *counter.Counter
	receiveBytes *counter.Counter
	server       *http.Server
	once         *sync.Once
}

// Protocol of HTTPServer
func (hs *HTTPServer) Protocol() string {
	return "http"
}

// Start HTTPServer
func (hs *HTTPServer) Start() error {
	hs.addr = hc.addr
	hs.requests = &counter.Counter{}
	hs.sendBytes = &counter.Counter{}
	hs.receiveBytes = &counter.Counter{}
	hs.once = &sync.Once{}

	log.Println(hs.Protocol(), " listen on ", hs.addr)
	hs.server = &http.Server{
		Addr: hc.addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			hs.requests.Add(1)
			hs.sendBytes.Add(float64(len(hc.responseBody)))
			hs.receiveBytes.Add(float64(len(body)))
			w.WriteHeader(hc.responseStatusCode)
			w.Write([]byte(hc.responseBody))
		}),
	}
	go hs.server.ListenAndServe()
	return nil
}

// Stop HTTPServer
func (hs *HTTPServer) Stop() error {
	log.Println(hs.Protocol(), " stoped")
	return hs.server.Close()
}

// Report HTTPServer status
func (hs *HTTPServer) Report() string {
	r := fmt.Sprintf("| %-15.0f | %-15.0f | %-15.0f |", hs.requests.Get(), hs.sendBytes.Get(), hs.receiveBytes.Get())
	hs.once.Do(func() {
		r = fmt.Sprintf("| %-15s | %-15s | %-15s |", "requests", "send bytes", "receive bytes")
	})
	return r
}
