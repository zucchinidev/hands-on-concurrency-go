package main

import (
	"fmt"
	"github.com/zucchinidev/hands-on-concurrency-go/tcp-memcache/cache"
	"github.com/zucchinidev/hands-on-concurrency-go/tcp-memcache/tcp"
	"log"
)

func main() {
	cac := cache.New()
	tcpServer := tcp.New(":8080")

	if errListening := tcpServer.Listen(); errListening != nil {
		log.Fatal(errListening)
	}
	defer tcpServer.Close()

	for {
		conn, _ := tcpServer.Accept()
		_, _ = fmt.Fprintf(conn, "Welcome to CACHE 1.0\n->")
		go tcpServer.Invoke(conn, cac)
	}
}
