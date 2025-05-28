package server

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseServer struct {
	privateKey *rsa.PrivateKey
}

func (ls *LicenseServer) Init(privateKeyFilename string) error {
	privateKeyFile, err := os.ReadFile(privateKeyFilename)
	if err != nil {
		log.Fatal("Error reading private key:", err)
		return err
	}
	ls.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		log.Fatal("Error parsing private key:", err)
		return err
	}
	return nil
}

func (ls *LicenseServer) GenerateLicense(features []string, durationMonth int) (string, error) {
	claims := jwt.MapClaims{
		"features":   features,
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().AddDate(0, durationMonth, 0).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(ls.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
