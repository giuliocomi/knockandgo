package network

import (
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
	ip_to_whitelist string
}

func NewUdpClient(server_address string, server_port int, knock_port int, certpath string, timeout int, ip_to_whitelist string) udp_client {
	c := udp_client{server_address, server_port, knock_port, certpath, timeout, ip_to_whitelist}
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
	msg := NewMessage(c.knock_port, 0, c.ip_to_whitelist, c.timeout, false)
	json_marshalled := Encode_message((msg))
	//send knock message to server
	encrypted_msg := crypto.Encrypt(string(json_marshalled), c.certpath+"public.pem")
	conn.Write([]byte(encrypted_msg))

	//read which forwarding port has been picked
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	
	log.Println(string(buffer))
}
