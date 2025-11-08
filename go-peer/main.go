package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"syscall"
	"time"
)

func main() {
	go broadcastPeer()
	go listenToPeerBroadcasts()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "ok",
			Message: "hello from go server",
		})
	})

	log.Println("server starting...")
	log.Fatal(http.ListenAndServe(":0", nil))
}

func broadcastPeer() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9999,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for range 10 {
		_, err = conn.Write([]byte("Hello"))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func listenToPeerBroadcasts() {
	var lc net.ListenConfig
	lc.Control = func(network, address string, c syscall.RawConn) error {
		var err error
		c.Control(func(fd uintptr) {
			err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if err != nil {
				return
			}

			err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		})

		return err
	}

	conn, err := lc.ListenPacket(context.Background(), "udp", ":9999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, src, _ := conn.ReadFrom(buf)
		fmt.Printf("from %v: %s \n", src, string(buf[:n]))
	}
}
