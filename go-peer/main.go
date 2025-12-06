package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/fagbenjaenoch/pitupi-chat/chat"
)

const BROADCAST_DELAY = 5 // seconds

var randId string = strconv.Itoa(rand.IntN(1000))
var port int

type Peer struct {
	Id      string
	Address string
}

var peersDiscovered map[string]Peer
var peerBroadcastPort int = 9999
var msgBroadcastPort int = 9998

func main() {
	flag.IntVar(&port, "port", 9000, "Port to run the peer on")
	flag.Parse()

	myAddr := fmt.Sprintf("%s:%d", getMyIpV4Address(), port)

	fmt.Printf("Here's your ip address: %s\n", myAddr)
	fmt.Printf("Here's your id: %s\n", randId)

	peersDiscovered = make(map[string]Peer)

	go broadcastPeer()
	go listenToPeerBroadcasts(broadcastPort)
	go listenForMessages(myAddr)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if utf8.RuneCountInString(line) == 0 {
			continue
		}

		parser := chat.NewParser()
		msg := parser.Parse(line)

		switch msg.Kind() {
		case "command":
			execCommand(msg.GetParts())
		case "mention":
			execMention(msg.Value())
		default:
			broadcastMessage(msg.Value())
		}

		fmt.Println("you:", line)

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stdout, "could not read standard input", err)
		}
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(struct {
	// 		Status  string `json:"status"`
	// 		Message string `json:"message"`
	// 	}{
	// 		Status:  "ok",
	// 		Message: "hello from go server",
	// 	})
	// })
	//
	// log.Println("server starting...")
	// log.Fatal(http.ListenAndServe(":0", nil))
}

func broadcastPeer() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: peerBroadcastPort,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		_, err = fmt.Fprintf(conn, "%s", randId)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * BROADCAST_DELAY)
	}
}

func createReusablePort() net.ListenConfig {
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

	return lc
}

func listenToPeerBroadcasts(port int) {
	lc := createReusablePort()

	conn, err := lc.ListenPacket(context.Background(), "udp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, src, _ := conn.ReadFrom(buf)
		peerId := string(buf[:n])

		peersDiscovered[peerId] = Peer{
			Id:      peerId,
			Address: src.String(),
		}
	}
}

func listenForMessages(address string) {
	conn, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleNewConnection(conn)
	}
}

func handleNewConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Hello from the server"))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Message from peer: %s", buf[:n])
}

func getMyIpV4Address() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("could not get peer's ip address")
	}

	var found []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				found = append(found, ipnet.IP.String())
			}
		}
	}

	return found[1] // for some reason, the second ip from the addresses found is the real ip
}

func broadcastMessage(message string) {
	conn, err := net.Dial("tcp", string(net.IPv4bcast)+":9998")
	if err != nil {
		fmt.Println("error occurred while sending message")
		return
	}

	conn.Write([]byte(message))
}

func listenToGeneralBroadcasts(port int) {
	lc := createReusablePort()

	conn, err := lc.ListenPacket(context.Background(), "udp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("could not listen to general broadcasts")
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, _ := conn.ReadFrom(buf)
		msg := buf[:n]

		fmt.Printf("#general: %s", msg)
	}
}

func execCommand(parts []string) {
	switch parts[0] {
	case "ls":
		listAllPeers()
	}
}

func execMention(m string) {}

func listAllPeers() {
	green := "\x1b[32m"
	reset := "\x1b[0m"

	fmt.Println(green + "All Peers" + reset)
	for _, peer := range peersDiscovered {
		if randId == peer.Id {
			fmt.Printf("%s (you)\n", peer.Id)
			continue
		}
		fmt.Println(peer.Id)
	}
}
