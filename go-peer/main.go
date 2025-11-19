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
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

const BROADCAST_DELAY = 5 // seconds

var randId string = strconv.Itoa(rand.IntN(1000))
var port int

type Peer struct {
	Id      string
	Address string
}

var peersDiscovered map[string]Peer

func main() {
	flag.IntVar(&port, "port", 9000, "Port to run the peer on")
	flag.Parse()

	myAddr := fmt.Sprintf("%s:%d", getMyIpV4Address(), port)

	fmt.Printf("Here's your ip address: %s\n", myAddr)
	fmt.Printf("Here's your id: %s\n", randId)

	peersDiscovered = make(map[string]Peer)

	go broadcastPeer()
	go listenToPeerBroadcasts()
	go listenForMessages(myAddr)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if utf8.RuneCountInString(line) == 0 {
			continue
		}

		m := strings.SplitN(line, " ", 1)

		switch m[0][:1] { // first letter of the message
		case "@": // mentions
			sendTo := m[0][1:]
			broadcastMessage(strings.Join(m[1:], ""), sendTo)
		case "!": // user commands
			fmt.Println("you just typed a command")
			execCommand(m[0][1:])
		default:
			fmt.Println("info: message will be broadcasted to all peers")
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
		Port: 9999,
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
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}

	// conn, err := net.DialTimeout("tcp", "0.0.0.0:0", time.Second*1)
	// if err != nil {
	// 	if opErr, ok := err.(*net.OpError); ok && opErr.Addr != nil {
	// 		return opErr.Addr.String()
	// 	}
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	//
	// localAddr := conn.LocalAddr().(*net.TCPAddr).IP.String()
	// return localAddr
	return ""
}

func broadcastMessage(message, id string) {
	recipient, ok := peersDiscovered[id]
	if !ok {
		fmt.Println("could not find peer")
		return
	}
	conn, err := net.Dial("tcp", recipient.Address)
	if err != nil {
		fmt.Println(recipient.Address)
		fmt.Println("error occurred while sending message")
		return
	}

	conn.Write([]byte(message))
}

func execCommand(command string) {

}
