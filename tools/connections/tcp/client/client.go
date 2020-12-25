package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var connections int64 = 0
var lock sync.Mutex

func connectionsPlus(conns int64) {
	lock.Lock()
	defer lock.Unlock()
	connections += conns
}

func connectionsMinus(conns int64) {
	lock.Lock()
	defer lock.Unlock()
	connections -= conns
}

func handleConn(ctx context.Context, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	connectionsPlus(1)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := conn.Write([]byte("hello\n"))
			if err != nil {
				connectionsMinus(1)
				conn.Close()
				return
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	var parallel int
	var addr string
	var interval time.Duration
	flag.IntVar(&parallel, "parallel", 1, "how many client to connect server")
	flag.StringVar(&addr, "addr", ":9090", "server address")
	flag.DurationVar(&interval, "interval", time.Second, "send data interval")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < parallel; i++ {
		go handleConn(ctx, addr)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case sig := <-ch:
			log.Println("server stop: ", sig)
			os.Exit(1)
		default:
			log.Printf("connections: %d\n", connections)
			time.Sleep(interval)
		}
	}
}

