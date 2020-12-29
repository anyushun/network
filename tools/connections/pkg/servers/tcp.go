package protocols

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"net-tools/pkg/counter"
)

type tcpServerConfig struct {
	addr           string
	maxPackageSize int
}

var tsc = &tcpServerConfig{}

func init() {
	flag.StringVar(&tsc.addr, "tcp-address", "0.0.0.0:9090", "address tcp server listen")
	flag.IntVar(&tsc.maxPackageSize, "tcp-pkg-size", 4096, "tcp max package size")
}

// TCPServer struct
type TCPServer struct {
	addr        string
	connections *counter.Counter
	packages    *counter.Counter
	bytes       *counter.Counter
	listen      net.Listener
	cancel      context.CancelFunc
	once        *sync.Once
}

// Protocol of the server
func (ts *TCPServer) Protocol() string {
	return "tcp"
}

// Start TCPServer
func (ts *TCPServer) Start() error {
	ts.addr = tsc.addr
	ts.connections = &counter.Counter{}
	ts.packages = &counter.Counter{}
	ts.bytes = &counter.Counter{}
	ts.once = &sync.Once{}

	log.Println(ts.Protocol(), " listen on ", ts.addr)
	listen, err := net.Listen("tcp", ts.addr)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	ts.listen = listen

	ctx, cancel := context.WithCancel(context.Background())
	ts.cancel = cancel
	go ts.handleAccept(ctx)
	return nil
}

func (ts *TCPServer) handleAccept(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := ts.listen.Accept()
			if err != nil {
				log.Fatalln(err)
				return
			}
			go ts.handleConnection(ctx, conn)
		}
	}
}

func (ts *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	ts.connections.Add(1)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, tsc.maxPackageSize)
			n, err := conn.Read(data)
			if err != nil {
				conn.Close()
				ts.connections.Sub(1)
				return
			}
			ts.packages.Add(1)
			ts.bytes.Add(float64(n))
		}
	}
}

// Stop TCPServer
func (ts *TCPServer) Stop() error {
	log.Println(ts.Protocol(), " stoped")
	ts.cancel()
	return ts.listen.Close()
}

// Report TCPServer status
func (ts *TCPServer) Report() string {
	r := fmt.Sprintf("| %-15.0f | %-15.0f | %-15.0f |", ts.connections.Get(), ts.packages.Get(), ts.bytes.Get())
	ts.once.Do(func() {
		r = fmt.Sprintf("| %-15s | %-15s | %-15s |", "connections", "packages", "bytes")
	})
	return r
}
