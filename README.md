# knockandgo

knockandgo is a port knocking solution developed in Golang.

knockandgo permits to conceal the presence of specific services from massive Internet scanners or intranet attackers.
It is best suited for network administrators that want to access remotely some network services but also are in need of hiding their presence. 

Notably, knockandgo is cross-platform and does not require any elevated privileges (root/administrator) to perform its duty.

### Features
- Cross-platform (for both clients and server instances)
- Lightweight
- Easy to setup and use
- Configurable timeouts
- Does not rely on monitoring logs data
- Does not need to run in kernelspace
- Does not require root/administrator privileges to accomplish its task
- IP spoofing mitigation acquired through both IP whitelisting and the native presence of 'non-guessable' sequence numbers in TCP packets
- Reply attack mitigation by the presence of a timestamp in knock requests
- Integrity check to mitigate tampering attempts
- Does never expose the true service port but open a random port and then forward the traffic between the client and that service
- Hard to fingerprint thanks to the encrypted traffic and a UDP random port to listen on

### How it works

1) The UDP server and the UDP clients need to have their private certificate (client|server)\_private.pem and the public certificate (client|server)\_public.pem of the other to successfully exchange messages.
2) The UDP client sends a crafted message containing some fields such as the port to access.
3) The UDP server does authentication and authorization checks and then, if everything is ok, instantiates a TCP forwarding server that listens on a random port. This forwarding port and the timeout of the connection are sent in a message response.
4) The client outputs the information regarding the success, the forwarding port and the timeout choosed by the server and then exits.
5) Now it is possible to reach and use the target service. After the first successfull TCP connection, the TCP forwarding server stops.


### Prerequisites

knockandgo requires Golang

### Installation
#### Get the tool
```
go get https://github.com/giuliocomi/knockandgo
```
#### Start the tool
Run main.go and choose what instance (server or client) you want to start

```
Options:

  -a string
        server address (IP) (default "localhost")
  -c string
        the path to the PEM certificate (public in case of client, private in case of server (default "./certs/")
  -f int
        number of maximum tcp wrappers to instantiate (default 5)
  -i string
        the source IP to whitelist for the tcp forwarding connections (default "127.0.0.1")
  -k int
        port to open via knock technique (default 80)
  -m string
        modality of operation: server (s) or client (c) (default "c")
  -s int
        server port (default 8080)
  -t int
        timeout in seconds (default 86400)
```

### Usage

- Server
```
go run main.go -m s -s 1337 22 8080 8081
```
- Clients
```
go run main.go -m c -a 120.23.21.212 -s 1337 -i 5.223.30.120 -t 300 -k 22
```

Note: clients and server exchange messages encrypted with RSA. Therefore, to correctly send and read their content it is necessary that the 4 \*.pem files are generated in advanced and shared between clients and the server.

#### Examples

(1) Simple demostration of the output of both server (on the left) and client (on the right) instances

![alt text](https://imgur.com/h0WZ62C.png)

(2) The SSH service is available only for localhost connection on port 22, a windows client access the SSH service after the 'knock'

[![asciicast](https://asciinema.org/a/a6UMXFvBjwxsQPxLTUk3031RU.png)](https://asciinema.org/a/a6UMXFvBjwxsQPxLTUk3031RU)

![alt text](https://imgur.com/e6Aus85.png)

![alt text](https://imgur.com/tvPRRR0.png)

### Roadmap
* [ ] Integration with firewall (iptables for Linux, windows firewall API for Windows)
* [ ] Add support for IPv6

## Issues
Spot a bug? Please create an issue here on GitHub (https://github.com/giuliocomi/knockandgo/issues)

## License
This project is licensed under the Apache License 2.0
