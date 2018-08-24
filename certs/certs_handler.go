package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

func HandleCert(certpath string) {
	fmt.Println("Generating the key pairs...")
	//generate the key pair for the clients
	generateCerts(certpath, "client_public.pem", "client_private.pem")
	
	//generate the key pair for the server
	generateCerts(certpath, "server_public.pem", "server_private.pem")
	fmt.Println(certpath, "Finished successfully...")
}

func generateCerts(certpath, public, private string) {
// Create the keys
	priv, pub := GenerateRsaKeyPair()

	// Export the keys to pem string
	priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
	pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

	// Import the keys from pem string
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

	// Export the newly imported keys
	priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

	err := ioutil.WriteFile(certpath+public, []byte(pub_parsed_pem), 0400)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(certpath+private, []byte(priv_parsed_pem), 0400)
	if err != nil {
		panic(err)
	}

	// Check that the exported/imported keys match the original keys
	if priv_pem != priv_parsed_pem || pub_pem != pub_parsed_pem {
		panic("Failure: Export and Import did not result in same Keys")
	}
}

func CheckCerts(certpath, modality string) bool {
	var file string
	var errf, errc error

	switch modality {
	case string("s"):
		file, errf = filepath.Abs(certpath+"server_private.pem")
		_, errc = ioutil.ReadFile(file)
		if errf != nil || errc != nil {
			return false
		}
		file, errf = filepath.Abs(certpath+"client_public.pem")
		_, errc = ioutil.ReadFile(file)
		if errf != nil || errc != nil {
			return false
		}
		return true
	case string("c"):
		file, errf = filepath.Abs(certpath+"client_private.pem")
		_, errc = ioutil.ReadFile(file)
		if errf != nil || errc != nil {
			return false
		}
		file, errf = filepath.Abs(certpath+"server_public.pem")
		_, errc = ioutil.ReadFile(file)
		if errf != nil || errc != nil {
			return false
		}
		return true
	default:
		return false	
	}
}
