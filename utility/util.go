package utility

import (
    "log"
    "net"
    "strconv"
    "time"
    "math/rand"
    "strings"
)

//verify if the target port is in state 'open'
func CheckConnection(server_address string, fport int) bool {

    conn, err := net.Dial("tcp", server_address + ":" + strconv.Itoa(fport))
    if err != nil {
        log.Println("Connection error:", err)
        return false
    }
    defer conn.Close()
    return true
}

//pick a random port available to set the TCP forwarder to listen on 
func RandomPort() int {
    port_available := false
    var rort int
    rand.Seed(time.Now().Unix())

    for port_available == false {
        rort := rand.Intn(65535 - 1025) + 1025
        if !CheckConnection("127.0.0.1", rort) {
            return rort
        }
    }
    return rort
}

func ContainsPort(s[] int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

//validate IP v4 | credits to: https://github.com/asaskevich/govalidator/blob/f9ffefc3facfbe0caee3fea233cbb6e8208f4541/validator.go
func IsValidIP4(ipAddress string) bool {
    ip := net.ParseIP(ipAddress)
    if (ip != nil && strings.Contains(ipAddress, ".")) || ipAddress == "localhost" {
        return true
    }
    return false
}

func SliceAtoi(slice_array[] string)([] int, error) {
    slice_of_ints := make([] int, 0, len(slice_array))
    for _,
    a := range slice_array {
        integer, err := strconv.Atoi(a)
        if err != nil {
            return slice_of_ints, err
        }
        slice_of_ints = append(slice_of_ints, integer)
    }
    return slice_of_ints,
    nil
}

//this awful function has the aim to recover from unexpected condition of tcp connections errors/kills
func HandlePanic() {
    if r := recover();
    r != nil {
        log.Println("Gracefully handling the panic:", r)
        return
    }
}
