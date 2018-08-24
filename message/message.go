package message

import (
	"encoding/json"
	"log"
	"time"
)

type message struct {
	Knock_port   int  `json:"Knock_port"`   //set only in client requests
	Forward_port int  `json:"Forward_port"` //set only in server responses
	Ip_to_whitelist string `json:"Ip_to_whitelist"` //the IP to whitelist
	Timeout      int  `json:"Timeout"`      //used in request as a graceful proposal by client
	Result       bool `json:"Result"`       //used in the response from the server
	Timestamp      int64 `json:"Timestamp"` 	//used to restring the time period for which a message is valid
}

func NewMessage(knock_port, forward_port int, Ip_to_whitelist string, timeout int, result bool, timestamp int64) message {
	m := message{knock_port, forward_port, Ip_to_whitelist, timeout, result, timestamp}
	return m
}

func Encode_message(msg message) []byte {
	json_bytes, err_e := json.Marshal(msg)
	if err_e != nil {
		log.Println(err_e)
	}
	return json_bytes
}

func Decode_message(json_marshalled []byte) (message, error) {
	m := message{}
	err_d := json.Unmarshal(json_marshalled, &m)
	if err_d != nil {
		return m, err_d
	}
	return m, nil
}

func IsExpired(timestamp int64) bool {
	const delay = 10 //10 seconds choosed as an arbitrary value to prevent reply attacks and at the same time to allow legitimate client packets to reach the server before they expire
	tnow := time.Now().Unix()
	if (timestamp > tnow + delay) || (timestamp + delay < tnow) {
		return true
	}
	return false
}
