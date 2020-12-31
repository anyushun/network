package clientsets

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"net-tools/pkg/counter"
	"net-tools/pkg/utils"
)

type tcpClientSetConfig struct {
	parallel        int
	addr            string
	packageSize     int
	packageInterval time.Duration
}

var tcsc = &tcpClientSetConfig{}

func init() {
	flag.IntVar(&tcsc.parallel, "tcp-parallel", 1, "how many tcp connections to server")
	flag.StringVar(&tcsc.addr, "tcp-address", ":9090", "tcp server address connect to")
	flag.IntVar(&tcsc.packageSize, "tcp-package-size", 64, "package send to server in bytes")
	flag.DurationVar(&tcsc.packageInterval, "tcp-package-interval", time.Second, "interval of packages")
}

// TCPClientSet impl
type TCPClientSet struct {
	connections *counter.Counter
	packages    *counter.Counter
	bytes       *counter.Counter
	ctx         context.Context
	cancel      context.CancelFunc
	once        *sync.Once
}

// Protocol of TCPClientSet
func (tcs *TCPClientSet) Protocol() string {
	return "tcp"
}

func (tcs *TCPClientSet) handleConn() {
	conn, err := net.Dial("tcp", tcsc.addr)
	if err != nil {
		log.Println(err)
		return
	}
	tcs.connections.Add(1)

	for {
		select {
		case <-tcs.ctx.Done():
			return
		default:
			_, err := conn.Write(utils.RandonBytes(tcsc.packageSize))
			if err != nil {
				tcs.connections.Sub(1)
				conn.Close()
				return
			}
			tcs.packages.Add(1)
			tcs.bytes.Add(float64(tcsc.packageSize))
			time.Sleep(tcsc.packageInterval)
		}
	}
}

// Start TCPClientSet
func (tcs *TCPClientSet) Start() error {
	tcs.connections = &counter.Counter{}
	tcs.packages = &counter.Counter{}
	tcs.bytes = &counter.Counter{}
	tcs.once = &sync.Once{}

	tcs.ctx, tcs.cancel = context.WithCancel(context.Background())
	for i := 0; i < tcsc.parallel; i++ {
		go tcs.handleConn()
	}
	return nil
}

// Stop TCPClientSet
func (tcs *TCPClientSet) Stop() error {
	tcs.cancel()
	return nil
}

// Report TCPClientSet
func (tcs *TCPClientSet) Report() string {
	r := fmt.Sprintf("| %-15.0f | %-15.0f | %-15.0f |", tcs.connections.Get(), tcs.packages.Get(), tcs.bytes.Get())
	tcs.once.Do(func() {
		r = fmt.Sprintf("| %-15s | %-15s | %-15s |", "connections", "packages", "bytes")
	})
	return r
}
