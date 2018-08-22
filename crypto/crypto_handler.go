package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"path/filepath"
)

var random io.Reader

func Encrypt(msg string, certificate string) (string, error) {
	random = rand.Reader
	var pub *rsa.PublicKey

	file, _ := filepath.Abs(certificate)
	public, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	pub_block, _ := pem.Decode(public)
	pubInterface, parseErr := x509.ParsePKIXPublicKey(pub_block.Bytes)
	if parseErr != nil {
		return "", err
	}
	pub = pubInterface.(*rsa.PublicKey)

	encryptedData, encryptErr := rsa.EncryptPKCS1v15(random, pub, []byte(msg))
	if encryptErr != nil {
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
		return "", err
	}

	priv_block, _ := pem.Decode(private)
	privateKeyBlock := priv_block

	pri, parseErr := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if parseErr != nil {
		return "", err
	}

	decryptedData, decryptErr := rsa.DecryptPKCS1v15(random, pri, []byte(encrypted_msg))
	if decryptErr != nil {
		return "", err
	}
	return string(decryptedData), nil
}
