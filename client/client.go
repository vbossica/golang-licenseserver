package client

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseClient struct {
	publicKey *rsa.PublicKey
	license   string
}

func (lc *LicenseClient) Init(publicKeyFilename string, licenseFilename string) error {
	// read the public key from the file
	publicKeyFile, err := os.ReadFile(publicKeyFilename)
	if err != nil {
		return fmt.Errorf("error reading public key: %w", err)
	}
	lc.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}

	// read the license from the file
	licenseBytes, err := os.ReadFile(licenseFilename)
	if err != nil {
		return fmt.Errorf("error reading license file: %w", err)
	}
	lc.license = string(licenseBytes)

	return nil
}

func (lc *LicenseClient) VerifyLicense(feature string) (bool, error) {
	token, err := jwt.Parse(lc.license, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return lc.publicKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims:", claims)

		if features, ok := claims["features"].([]interface{}); ok {
			for _, f := range features {
				if f == feature {
					return true, nil
				}
			}
		}
		return false, fmt.Errorf("feature not found in the claims")
	}
	return false, fmt.Errorf("invalid token")
}
