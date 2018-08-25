# knockandgo

knockandgo is a port knocking solution developed in Golang.

knockandgo permits to conceal the presence of specific services from massive Internet scanners or intranet attackers.
It is best suited for network administrators that want to access remotely some network services but also are in need of hiding their presence. 

Notably, knockandgo is cross-platform and does not require any elevated privileges (root/administrator) to perform its duty.

### Features
- Cross-platform (for both clients and server instances)
- Lightweight
- Easy to setup and use
- Does not rely on monitoring logs data
- Does not need to run in kernelspace
- Does not require root/administrator privileges to accomplish its task
- IP spoofing prevention acquired through both IP whitelisting and the native presence of 'non-guessable' sequence numbers in TCP packets
- Reply attack mitigation by the presence of expiration time in messages
- Does never expose the true service port but open a random port and then forward the traffic between the client and that service
- Configurable timeouts
- Message exchange based on asymmetric encryption
- The server listen on UDP for semplicity and for reducing the probability that automatic scanners discover its UDP port 'open' 
- Hard to fingerprint thanks to the encrypted traffic and a UDP random port to listen on

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

### Usage Examples

- Server
```
go run main.go -m s -s 1337 22 8080 8081
```
- Clients
```
go run main.go -m c -a 120.23.21.212 -s 1337 -i 5.223.30.120 -t 300 -k 22
```

Note: clients and server exchange messages encrypted with RSA. Therefore, to correctly send and read their content it is necessary that the 4 \*.pem files are generated in advanced and shared between clients and server.

#### Scenario:
a SSH service is available only for localhost connection on port 22
##### Steps
1) Server console:

[![asciicast](https://asciinema.org/a/8IZnS3bwImcZCIESgl41xrQmt.png?autoplay=1)](https://asciinema.org/a/8IZnS3bwImcZCIESgl41xrQmt)
2) Client console:

![alt text](https://imgur.com/1lat28c.png)
3) Now the SSH service is available:

![alt text](https://imgur.com/fPOhFF4.png)

### Roadmap
* [ ] Integration with firewall (iptables for Linux, windows firewall API for Windows)
    
## License

This project is licensed under the Apache License 2.0
