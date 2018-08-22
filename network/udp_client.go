package network

import (
	"log"
	"net"
	"strconv"
	"bytes"

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
	if encrypted_msg, erre := crypto.Encrypt(string(json_marshalled), c.certpath+"server_public.pem"); erre != nil {
		log.Println("Error encrypting the message", erre)
		return
	} else {

		conn.Write([]byte(encrypted_msg))
	}

	//read which forwarding port has been picked
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	json_resp_marshalled, errd := crypto.Decrypt(string(bytes.Trim(buffer, "\x00")), c.certpath+"client_private.pem")
	if errd != nil {
		log.Println("Error decrypting the response from the server")
		return
	}
	json_resp_unmarshalled := Decode_message([]byte(json_resp_marshalled))
	log.Println(json_resp_unmarshalled)
}
