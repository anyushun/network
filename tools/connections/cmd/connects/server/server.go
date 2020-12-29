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

	protocols "net-tools/pkg/servers"
)

// Server interface
type Server interface {
	// Protocol of the server
	Protocol() string

	// Start the server
	Start() error

	// Stop the server
	Stop() error

	// Report the server status
	Report() string
}

// ServerEngine management the server to start
type ServerEngine struct {
	// ReportInterval interval of report
	ReportInterval time.Duration

	// servers of the server want to start, key: protocol, value: Server{}
	servers sync.Map

	// startedServers
	startedServers sync.Map

	lock sync.Mutex
}

// Load all servers
func (se *ServerEngine) Load() error {
	fs := &protocols.FakeServer{}
	se.servers.Store(fs.Protocol(), fs)

	ts := &protocols.TCPServer{}
	se.servers.Store(ts.Protocol(), ts)

	hs := &protocols.HTTPServer{}
	se.servers.Store(hs.Protocol(), hs)
	return nil
}

// Report Server status
func (se *ServerEngine) Report(protocol string) error {
	se.lock.Lock()
	defer se.lock.Unlock()
	server, found := se.findServer(protocol)
	if !found {
		return fmt.Errorf("%s incorrect server", protocol)
	}
	log.Println(server.Report())
	return nil
}

// ReportAll Server status
func (se *ServerEngine) ReportAll() {
	se.startedServers.Range(func(key, value interface{}) bool {
		protocol, ok := key.(string)
		if !ok {
			return true
		}
		se.Report(protocol)
		return true
	})
}

func (se *ServerEngine) findServer(protocol string) (Server, bool) {
	s, ok := se.servers.Load(protocol)
	if !ok {
		return nil, false
	}
	server, ok := s.(Server)
	if !ok {
		return nil, false
	}
	return server, true
}

// Start a server
func (se *ServerEngine) Start(protocol string) error {
	se.lock.Lock()
	defer se.lock.Unlock()
	server, found := se.findServer(protocol)
	if !found {
		return fmt.Errorf("%s incorrect server", protocol)
	}
	err := server.Start()
	if err == nil {
		se.startedServers.Store(server.Protocol(), server)
	}
	return err
}

// Stop Server by protocol
func (se *ServerEngine) Stop(protocol string) error {
	se.lock.Lock()
	defer se.lock.Unlock()
	server, found := se.findServer(protocol)
	if !found {
		return fmt.Errorf("%s incorrect server", protocol)
	}
	err := server.Stop()
	if err == nil {
		se.startedServers.Delete(protocol)
	}
	return err
}

// StopAll Server
func (se *ServerEngine) StopAll() error {
	se.startedServers.Range(func(key, value interface{}) bool {
		protocol, ok := key.(string)
		if !ok {
			return true
		}
		se.Stop(protocol)
		return true
	})
	return nil
}

// NewServerEngine return a ServerEngine
func NewServerEngine() *ServerEngine {
	se := &ServerEngine{}
	se.Load()
	return se
}

func main() {
	var protocol string
	var reportInterval time.Duration

	flag.StringVar(&protocol, "protocol", "fake", "start server protocol")
	flag.DurationVar(&reportInterval, "report-interval", time.Second, "interval of report")
	flag.Parse()

	se := NewServerEngine()
	log.Println("start server: ", protocol)
	se.Start(protocol)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case sig := <-ch:
			log.Println("server stop: ", sig)
			se.StopAll()
			return
		default:
			se.ReportAll()
			time.Sleep(reportInterval)
		}
	}
}
