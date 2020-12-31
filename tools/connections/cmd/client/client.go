package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net-tools/pkg/clientsets"
)

// ClientSet interface spec
type ClientSet interface {
	// Protocol of ClientSet used
	Protocol() string

	// Start ClientSet and test
	Start() error

	// Stop ClientSet
	Stop() error

	// Report ClientSet status
	Report() string
}

// ClientSetEngine management the ClientSet to start
type ClientSetEngine struct {
	clientSets        sync.Map
	startedClientSets sync.Map
	lock              sync.Mutex
}

// Load all ClientSet
func (cse *ClientSetEngine) Load() error {
	fc := &clientsets.FakeClientSet{}
	cse.clientSets.Store(fc.Protocol(), fc)

	tc := &clientsets.TCPClientSet{}
	cse.clientSets.Store(tc.Protocol(), tc)

	hc := &clientsets.HTTPClientSet{}
	cse.clientSets.Store(hc.Protocol(), hc)
	return nil
}

func (cse *ClientSetEngine) findClientSet(protocol string) (ClientSet, bool) {
	cs, ok := cse.clientSets.Load(protocol)
	if !ok {
		return nil, false
	}
	clientSet, ok := cs.(ClientSet)
	if !ok {
		return nil, false
	}
	return clientSet, true
}

// Report ClientSet status
func (cse *ClientSetEngine) Report(protocol string) error {
	cse.lock.Lock()
	defer cse.lock.Unlock()
	cs, found := cse.findClientSet(protocol)
	if !found {
		return fmt.Errorf("%s incorrect clientset", protocol)
	}
	log.Println(cs.Report())
	return nil
}

// ReportAll ClientSet status
func (cse *ClientSetEngine) ReportAll() {
	cse.startedClientSets.Range(func(key, value interface{}) bool {
		protocol, ok := key.(string)
		if !ok {
			return true
		}
		cse.Report(protocol)
		return true
	})
}

// Start a ClientSet
func (cse *ClientSetEngine) Start(protocol string) error {
	cse.lock.Lock()
	defer cse.lock.Unlock()
	cs, found := cse.findClientSet(protocol)
	if !found {
		return fmt.Errorf("%s incorrect clientset", protocol)
	}
	err := cs.Start()
	if err == nil {
		cse.startedClientSets.Store(cs.Protocol(), cs)
	}
	return err
}

// Stop ClientSet by protocol
func (cse *ClientSetEngine) Stop(protocol string) error {
	cse.lock.Lock()
	defer cse.lock.Unlock()
	cs, found := cse.findClientSet(protocol)
	if !found {
		return fmt.Errorf("%s incorrect clientset", protocol)
	}
	err := cs.Stop()
	if err == nil {
		cse.startedClientSets.Delete(protocol)
	}
	return err
}

// StopAll ClientSet
func (cse *ClientSetEngine) StopAll() error {
	cse.startedClientSets.Range(func(key, value interface{}) bool {
		protocol, ok := key.(string)
		if !ok {
			return true
		}
		cse.Stop(protocol)
		return true
	})
	return nil
}

// NewClientSetEngine return a ClientSetEngine
func NewClientSetEngine() *ClientSetEngine {
	cse := &ClientSetEngine{}
	cse.Load()
	return cse
}

func main() {
	var protocol string
	var reportInterval time.Duration

	flag.StringVar(&protocol, "protocol", "fake", "start clientset protocol")
	flag.DurationVar(&reportInterval, "report-interval", time.Second, "interval of report")
	flag.Parse()

	cse := NewClientSetEngine()
	log.Println("start clientset: ", protocol)
	cse.Start(protocol)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case sig := <-ch:
			log.Println("clientset stop: ", sig)
			cse.StopAll()
			return
		default:
			cse.ReportAll()
			time.Sleep(reportInterval)
		}
	}
}
