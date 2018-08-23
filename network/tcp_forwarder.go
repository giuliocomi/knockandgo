package network

import (
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/giuliocomi/knockandgo/utility"
)

type tcp_forwarder struct {
	server_address string
	server_port    int
	knock_port     int
	ip_to_whitelist string
	timeout        int
}

func NewTcpForwarder(server_address string, server_port int, knock_port int, ip_to_whitelist string, timeout int) tcp_forwarder {
	c := tcp_forwarder{server_address, server_port, knock_port, ip_to_whitelist, timeout}
	return c
}

func (t *tcp_forwarder) Listen() {
	defer utility.HandlePanic()
	timer := time.NewTimer(time.Duration(t.timeout) * time.Second)
	log.Println("timeout started")
	listener, err := net.Listen("tcp", t.server_address+":"+strconv.Itoa(t.server_port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Port forwarding server up and listening on ", t.server_address+":"+strconv.Itoa(t.server_port))


	conn, _ := listener.Accept()
	//check if the timeout has expired, if so kill the tcp wrapper and the active connection, otherwise accept the new incoming tcp connection and set a timeout for it
	select {
	case <-timer.C:
		Instantiated_forwarding_ports-- //a new slot for a tcp forwarder is available after this is freed
		listener.Close()
	default:
		conn.SetDeadline(time.Now().Add(time.Duration(t.timeout) * time.Second))
		//verify that the source Ip address of the connection is trusted
		if t.ip_to_whitelist == utility.GetStringIpFromAddr(conn) {
			handleConnection("127.0.0.1", t.knock_port, conn)
			listener.Close() //after a successfull TCP connection, close the port of this tcp wrapper
		} else {
			log.Println("The remote IP address is not in the whitelist")
		}
	}
	
}

func handleConnection(rtsh string, kp int, c net.Conn) {
	log.Println("Connection from : ", c.RemoteAddr())

	//resume if the target port goes down after the fist connection (as an example, nc -lvnp [PORT] and request from the client when the knock port is [PORT])
	defer utility.HandlePanic()

	remote, err := net.Dial("tcp", "127.0.0.1"+":"+strconv.Itoa(kp))
	if err != nil {
		panic("The net.Dial throws a panic error... we need to recover from this otherwise the UDP server itself is halted... ")
	}
	// goroutines to initiate bi-directional communication for local server with a remote client
	go io.Copy(remote, c)
	go io.Copy(c, remote)
}
