package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "io/ioutil"
    "path/filepath"
    "io"
)

var random io.Reader

func Encrypt(msg string, certificate string) string {
    random = rand.Reader
    var pub * rsa.PublicKey

    file, _ := filepath.Abs(certificate)
    public, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatalf("read key file: %s", err)
    }

    pub_block, _ := pem.Decode(public)
    pubInterface, parseErr := x509.ParsePKIXPublicKey(pub_block.Bytes)

    if parseErr != nil {
        log.Println("Load public key error")
        panic(parseErr)
    }
    pub = pubInterface.( * rsa.PublicKey)

    encryptedData, encryptErr := rsa.EncryptPKCS1v15(random, pub, [] byte(msg))
    if encryptErr != nil {
        panic(encryptErr)
    }

    return string(encryptedData)
}

func Decrypt(encrypted_msg string, certificate string) string {
    random = rand.Reader
    var pri * rsa.PrivateKey   

    file, _ := filepath.Abs(certificate)
    private, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatalf("read key file: %s", err)
    }

    priv_block, _ := pem.Decode(private)
    privateKeyBlock := priv_block

    pri, parseErr := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

    if parseErr != nil {
        log.Println("Load private key error")
        panic(parseErr)
    }

    decryptedData, decryptErr := rsa.DecryptPKCS1v15(random, pri, [] byte(encrypted_msg))
    if decryptErr != nil {
        log.Println(decryptErr)
    }

    return string(decryptedData)
}
