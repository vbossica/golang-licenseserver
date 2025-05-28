package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vbossica/golang-licenseserver/client"
	"github.com/vbossica/golang-licenseserver/server"
)

func getSerialNumber() (string, error) {
	return "device-123", nil
}

func createLicense(privateKey *rsa.PrivateKey, deviceId string, features []string) (string, error) {
	claims := jwt.MapClaims{
		"device_id":  deviceId,
		"features":   features,
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().Add(time.Hour * 24 * 365).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func verifyLicense(publicKey *rsa.PublicKey, tokenString string, deviceId string, feature string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims:", claims)

		if claims["device_id"] == deviceId {
			if features, ok := claims["features"].([]interface{}); ok {
				for _, f := range features {
					if f == feature {
						return true, nil
					}
				}
			}
			return false, fmt.Errorf("feature not found in the claims")
		} else {
			return false, fmt.Errorf("device ID does not match")
		}
	} else {
		fmt.Println("Invalid token!")
		return false, fmt.Errorf("invalid token")
	}
}

func main() {
	privateKey, publicKey, err := server.GenerateRsaKeyPair()
	if err != nil {
		fmt.Printf("Cannot generate RSA key pair: %s\n", err)
		os.Exit(1)
	}

	err = server.WritePrivateKeyToFile(privateKey, "server/testdata/private.pem")
	if err != nil {
		fmt.Printf("Error writing private key to file: %s\n", err)
		os.Exit(1)
	}

	err = server.WritePublicKeyToFile(publicKey, "client/testdata/public.pem")
	if err != nil {
		fmt.Printf("Error writing public key to file: %s\n", err)
		os.Exit(1)
	}

	// Load the private key from file
	privateKeyFile, err := os.ReadFile("server/testdata/private.pem")
	if err != nil {
		log.Fatal("Error reading private key:", err)
	}
	readPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		log.Fatal("Error parsing private key:", err)
	}

	// Load the public key from file
	publicKeyFile, err := os.ReadFile("client/testdata/public.pem")
	if err != nil {
		log.Fatal("Error reading public key:", err)
	}
	readPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		log.Fatal("Error parsing public key:", err)
	}

	// create and sign the claim
	deviceId, err := getSerialNumber()
	if err != nil {
		log.Fatal("Error getting device ID:", err)
	}
	log.Println("Device ID:", deviceId)

	features := []string{"feature-1", "feature-2"}

	license, err := createLicense(readPrivateKey, deviceId, features)
	if err != nil {
		log.Fatal("Error creating license:", err)
	}

	fmt.Println("License to be saved and used on the device:", license)
	// save the license to a file
	licenseFile, err := os.Create("license.txt")
	if err != nil {
		log.Fatal("Error creating license file:", err)
	}
	defer licenseFile.Close()

	_, err = licenseFile.WriteString(license)
	if err != nil {
		log.Fatal("Error writing license to file:", err)
	}
	fmt.Println("License saved to license.txt")

	// read the license from the file
	licenseBytes, err := os.ReadFile("license.txt")
	if err != nil {
		log.Fatal("Error reading license file:", err)
	}
	license = string(licenseBytes)

	// verify the claim on the device
	isValid, _ := verifyLicense(readPublicKey, license, deviceId, "feature-1")
	fmt.Println("Is claim valid:", isValid)

	isValid, _ = verifyLicense(readPublicKey, license, deviceId, "feature-3")
	fmt.Println("Is claim valid:", isValid)

	isValid, _ = verifyLicense(readPublicKey, license, "device-234", "feature-1")
	fmt.Println("Is claim valid:", isValid)

	isValid, _ = verifyLicense(readPublicKey, license+"tempered", deviceId, "feature-1")
	fmt.Println("Is claim valid:", isValid)
}
