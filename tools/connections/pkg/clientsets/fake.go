package clientsets

import (
	"fmt"
	"log"
)

// FakeClientSet is an example of Client
type FakeClientSet struct{}

// Protocol of FakeClientSet
func (fcs *FakeClientSet) Protocol() string {
	return "fake"
}

// Start FakeClientSet
func (fcs *FakeClientSet) Start() error {
	log.Println("start client: ", fcs.Protocol())
	return nil
}

// Stop FakeClientSet
func (fcs *FakeClientSet) Stop() error {
	log.Println("stop client: ", fcs.Protocol())
	return nil
}

// Report FakeClientSet status
func (fcs *FakeClientSet) Report() string {
	return fmt.Sprintf("report client: %s", fcs.Protocol())
}
