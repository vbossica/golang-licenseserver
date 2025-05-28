package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateRsaKeyPair(privateKeyPath, publicKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey

	err = writePrivateKeyToFile(privateKey, privateKeyPath)
	if err != nil {
		return fmt.Errorf("error writing private key to file: %s", err)
	}
	err = writePublicKeyToFile(publicKey, publicKeyPath)
	if err != nil {
		return fmt.Errorf("error writing public key to file: %s", err)
	}

	return nil
}

func writePrivateKeyToFile(privatekey *rsa.PrivateKey, filename string) error {
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePemFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error when create %s: %s \n", filename, err)
	}
	err = pem.Encode(privatePemFile, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("error when encode private pem: %s \n", err)
	}
	return nil
}

func writePublicKeyToFile(publickey *rsa.PublicKey, filename string) error {
	var publicKeyBytes, err = x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return fmt.Errorf("error when dumping publickey: %s \n", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPemFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error when create %s: %s \n", filename, err)
	}
	err = pem.Encode(publicPemFile, publicKeyBlock)
	if err != nil {
		return fmt.Errorf("error when encode public pem: %s \n", err)
	}
	return nil
}
