package message

import (
	"encoding/json"
	"log"
	"time"
	"errors"
	"crypto/sha256"
	"bytes"
)

type message struct {
	Knock_port   		int  `json:"Knock_port"`   	//set only in client requests
	Forwarding_port 	int  `json:"Forward_port"` 	//set only in server responses
	Ip_to_whitelist 	string `json:"Ip_to_whitelist"` //the IP to whitelist
	Timeout     		int  `json:"Timeout"`     	//used in request as a graceful proposal by client
	Result       		bool `json:"Result"`       	//used in the response from the server
	Timestamp       	int64 `json:"Timestamp"` 	//used to restring the time period for which a message is valid
	Checksum 		[32]byte  `json:"Checksum"`	//used to verify the integrity of the message fields
}

func NewMessage(k, f int, i string, t int, r bool, o int64) message {
	m := message{k, f, i, t, r, o,  sha256.Sum256([]byte(string(k)+string(i)+string(t)+string(o)))}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(m)
	m.Checksum = sha256.Sum256([]byte(buf.Bytes()))
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
	//check the integrity of the message received
	mi := NewMessage(m.Knock_port, m.Forwarding_port, m.Ip_to_whitelist, m.Timeout, m.Result, m.Timestamp)
	if m.Checksum != mi.Checksum {
		log.Println(m.Checksum)
		log.Println(mi.Checksum)
		return m, errors.New("The integrity check has failed. The message has been tampered.")		
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
