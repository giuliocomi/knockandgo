package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

var random io.Reader

func Encrypt(msg string, certificate string) (string, error) {
	random = rand.Reader
	var pub *rsa.PublicKey

	file, _ := filepath.Abs(certificate)
	public, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("read key file: %v", err)
		return "", err
	}

	pub_block, _ := pem.Decode(public)
	pubInterface, parseErr := x509.ParsePKIXPublicKey(pub_block.Bytes)
	if parseErr != nil {
		log.Println(parseErr)
		return "", err
	}
	pub = pubInterface.(*rsa.PublicKey)

	encryptedData, encryptErr := rsa.EncryptPKCS1v15(random, pub, []byte(msg))
	if encryptErr != nil {
		log.Println(encryptErr)
		return "", err
	}
	return string(encryptedData), nil
}

func Decrypt(encrypted_msg string, certificate string) (string, error) {
	random = rand.Reader
	var pri *rsa.PrivateKey

	file, _ := filepath.Abs(certificate)
	private, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("read key file: %v", err)
		return "", err
	}

	priv_block, _ := pem.Decode(private)
	privateKeyBlock := priv_block

	pri, parseErr := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if parseErr != nil {
		log.Println(parseErr)
		return "", err
	}

	decryptedData, decryptErr := rsa.DecryptPKCS1v15(random, pri, []byte(encrypted_msg))
	if decryptErr != nil {
		log.Println(decryptErr)
		return "", err
	}
	return string(decryptedData), nil
}
