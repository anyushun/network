package protocols

import (
	"fmt"
)

// FakeServer struct
type FakeServer struct{}

// Protocol of the server
func (fs *FakeServer) Protocol() string {
	return "fake"
}

// Start FakeServer
func (fs *FakeServer) Start() error {
	fmt.Println("Start server")
	return nil
}

// Stop TCPServer
func (fs *FakeServer) Stop() error {
	fmt.Println("Stop server")
	return nil
}

// Report FakeServer status
func (fs *FakeServer) Report() string {
	return "Report server"
}
