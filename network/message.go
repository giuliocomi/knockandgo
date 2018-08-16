package network

import (
    "encoding/json"
    "log"
)

type message struct {
    Knock_port int `json:"Knock_port"` //set only in client requests
    Forward_port int `json:"Forward_port"` //set only in server responses
    Timeout int `json:"Timeout"` //used in request as a graceful proposal by client
    Result bool `json:"Result"` // used in the response from the server
}

func NewMessage(knock_port, forward_port, timeout int, result bool) message {
    m := message {knock_port, forward_port, timeout, result}
    return m
}

func Encode_message(msg message)[] byte {
    json_bytes, err_e := json.Marshal(msg)
    if err_e != nil {
        log.Println(err_e)
    }
    return json_bytes
}

func Decode_message(json_marshalled[] byte) message {
    m := message{}
    err_d := json.Unmarshal(json_marshalled, &m)
    if err_d != nil {
        log.Println(err_d)
    }
    return m
}
