//
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func handleConn(ctx context.Context, conn net.Conn) {
	log.Println(conn.RemoteAddr(), " connected")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			netData, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			netData = strings.Replace(netData, "\n", "", 1)
			log.Printf("[%s]", string(netData))
		}
	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":9090", "tcp server address")
	flag.Parse()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	log.Println("Server started. Press Ctrl-C to stop server")

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case sig := <-ch:
			log.Println("server stop: ", sig)
			os.Exit(1)
		default:
			c, err := listen.Accept()
			if err != nil {
				log.Panicln(err)
			}
			go handleConn(ctx, c)
		}
	}
}

