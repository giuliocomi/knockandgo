package network

// TODO: SINGLETON

import (
    "log"
    "net"
    "strconv"
    "bytes"
	
    "github.com/giuliocomi/knockandgo/utility"
    "github.com/giuliocomi/knockandgo/crypto"
)

type udp_server struct {
    server_port int
    knockable_ports[] int
    max_forwarding_ports int
    certpath string
    timeout int
}

func NewUdpServer(server_port int, kports[] int, max_forwarding_ports int, certpath string, timeout int) udp_server {
    s := udp_server {server_port, kports, max_forwarding_ports, certpath, timeout}
    return s
}

func(s * udp_server) Run() {

    var instantiated_forwarding_ports int
    // listen to incoming udp packets
    log.Println("Whitelisted knockable ports:", s.knockable_ports)
    log.Println(s.server_port)
    pc, err := net.ListenPacket("udp", "0.0.0.0" + ":" + strconv.Itoa(s.server_port))
    if err != nil {
        log.Fatal(err)
    }
    defer pc.Close()

    for {
        buffer := make([] byte, 1024) 
        _,
        addr,
        err := pc.ReadFrom(buffer) 
        if err != nil {
            log.Println("Buffer error")
        }
        json_marshalled := crypto.Decrypt(string(bytes.Trim(buffer, "\x00")), s.certpath + "private.pem")
        json_unmarshalled := Decode_message([] byte(json_marshalled))
        kport := json_unmarshalled.Knock_port 

        client_timeout := json_unmarshalled.Timeout
        //pick a random port to start the tcp_wrapper on
        log.Println(json_unmarshalled)
        rort := utility.RandomPort()
        log.Println(rort)
        if (utility.ContainsPort(s.knockable_ports, kport)) && (instantiated_forwarding_ports < s.max_forwarding_ports) {
            //check if target port is open
            port_open := utility.CheckConnection("127.0.0.1", kport)
            if !port_open {
                pc.WriteTo(Encode_message(NewMessage(0, 0, func() int {
                    if s.timeout < client_timeout {
                        return s.timeout
                    } else {
                        return client_timeout
                    }
                }(), false)), addr)
                continue
            } else {
                forwarder := NewTcpForwarder("0.0.0.0", rort, kport, func() int {
                    if s.timeout < client_timeout {
                        return s.timeout
                    } else {
                        return client_timeout
                    }
                }())
                instantiated_forwarding_ports++
                go forwarder.Listen()

                //tcp forwarder port created
                log.Println(NewMessage(0, rort, 9999, true))
                pc.WriteTo(Encode_message(NewMessage(0, rort, func() int {
                    if s.timeout < client_timeout {
                        return s.timeout
                    } else {
                        return client_timeout
                    }
                }(), true)), addr) // return true: no error and port correctly opened
            }
        } else {
            if instantiated_forwarding_ports >= s.max_forwarding_ports {
                log.Println("Reached maximum number of available forwarding ports")
            } else {
                log.Println("Port is not whitelisted to be forwarded")
            }
            pc.WriteTo(Encode_message(NewMessage(kport, rort, 0, false)), addr)
        }
    }
}
