package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
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
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9999,
	})
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	buf := make([]byte, 1024)
	for {
		n, src, _ := ln.ReadFromUDP(buf)
		fmt.Printf("from %v: %s \n", src, string(buf[:n]))
	}
}
