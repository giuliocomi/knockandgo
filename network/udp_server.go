package network

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/giuliocomi/knockandgo/crypto"
	"github.com/giuliocomi/knockandgo/utility"
)

type udp_server struct {
	server_port          int
	knockable_ports      []int
	max_forwarding_ports int
	certpath             string
	timeout              int
}

var singleton *udp_server 
var once sync.Once
var Instantiated_forwarding_ports int //the forwarding ports currently in state 'open'
func (s *udp_server) SetUdpServer(s_port int, kports []int, max_f_ports int, cpath string, tout int) {
	s.server_port = s_port
	s.knockable_ports = kports
	s.max_forwarding_ports = max_f_ports
	s.certpath = cpath
	s.timeout = tout
}

func GetUDPServer() *udp_server {
	once.Do(func() {
		singleton = &udp_server{}
	    })
	return singleton
}

func (s *udp_server) Run() {
	
	// listen to incoming udp packets
	log.Println("Whitelisted knockable ports:", s.knockable_ports)
	log.Println(s.server_port)
	pc, err := net.ListenPacket("udp", "0.0.0.0"+":"+strconv.Itoa(s.server_port))
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buffer := make([]byte, 1024)
		_, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Println("Buffer error")
		}
		json_marshalled, errd := crypto.Decrypt(string(bytes.Trim(buffer, "\x00")), s.certpath+"server_private.pem")
		if errd != nil {
			log.Println("Error during the decryption of the message received")
			continue
		}
		json_unmarshalled := Decode_message([]byte(json_marshalled))
		kport := json_unmarshalled.Knock_port
		timestamp := json_unmarshalled.Timestamp
		ip_to_whitelist := json_unmarshalled.Ip_to_whitelist
		client_timeout := json_unmarshalled.Timeout
		//pick a random port to start the tcp_wrapper on
		rort := utility.RandomPort()
		is_expired := utility.IsExpired(timestamp)
		if (utility.ContainsPort(s.knockable_ports, kport)) && (Instantiated_forwarding_ports < s.max_forwarding_ports) && utility.IsValidIP4(ip_to_whitelist) && !is_expired {
			//check if target port is open
			port_open := utility.CheckConnection("127.0.0.1", kport)
			if !port_open {
				SendResponse(string(Encode_message(NewMessage(0, 0, ip_to_whitelist, func() int {
					if s.timeout < client_timeout {
						return s.timeout
					} else {
						return client_timeout
					}
				}(), false, time.Now().Unix()))), pc, s.certpath, addr)
			} else {
				forwarder := NewTcpForwarder("0.0.0.0", rort, kport, ip_to_whitelist, func() int {
					if s.timeout < client_timeout {
						return s.timeout
					} else {
						return client_timeout
					}
				}())
				Instantiated_forwarding_ports++
				go forwarder.Listen()
				//tcp forwarder port created

				SendResponse(string(Encode_message(NewMessage(kport, rort, ip_to_whitelist, func() int {
					if s.timeout < client_timeout {
						return s.timeout
					} else {
						return client_timeout
					}
				}(), true, time.Now().Unix()))), pc, s.certpath, addr) // result true: no error and port correctly opened
			}
		} else {
			if Instantiated_forwarding_ports >= s.max_forwarding_ports {
				log.Println("Reached maximum number of available forwarding ports")
			} else if !utility.IsValidIP4(ip_to_whitelist) {
				log.Println("The IP to whitelist is not in a correct IPv4 format")
			} else if utility.IsExpired(expiration) {
				log.Println("The message time validity is expired")
			} else {	
				log.Println("Port is not whitelisted to be forwarded")
			}
			SendResponse(string(Encode_message(NewMessage(kport, rort, "", 0, false, time.Now().Unix()))), pc, s.certpath, addr)
		}
	}
}

func SendResponse(encoded_response string, pc net.PacketConn, certpath string, addr net.Addr) {
	encrypted_response, err := crypto.Encrypt(encoded_response, certpath+"client_public.pem")
	if err != nil {
		log.Println("Error during encryption of the response", err)
		return
	}
	pc.WriteTo([]byte(encrypted_response), addr)
}

