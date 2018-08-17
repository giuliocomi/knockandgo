package network

import (
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/giuliocomi/knockandgo/crypto"
	"github.com/giuliocomi/knockandgo/utility"
)

type udp_client struct {
	server_address string
	server_port    int
	knock_port     int
	certpath       string
	timeout        int
}

func NewUdpClient(server_address string, server_port int, knock_port int, certpath string, timeout int) udp_client {
	c := udp_client{server_address, server_port, knock_port, certpath, timeout}
	return c
}

func (c *udp_client) Run() {
	//connect
	if !utility.IsValidIP4(c.server_address) {
		panic("Invalid IP v4")
	}
	conn, err := net.Dial("udp", c.server_address+":"+strconv.Itoa(c.server_port))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	//craft message
	msg := NewMessage(c.knock_port, 0, c.timeout, false)
	json_marshalled := Encode_message((msg))
	//send knock message to server
	encrypted_msg := crypto.Encrypt(string(json_marshalled), c.certpath+"public.pem")
	conn.Write([]byte(encrypted_msg))

	//read which forwarding port has been picked
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	json_unmarshalled := Decode_message([]byte(bytes.Trim(buffer, "\x00")))
	fport := json_unmarshalled.Forward_port

	//check if port is reachable
	port_open := utility.CheckConnection(c.server_address, fport)
	if port_open {
		log.Println(string(buffer))
	} else {
		log.Println("The forwarding port seems closed")
	}
}
