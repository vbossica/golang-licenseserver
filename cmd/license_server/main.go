package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vbossica/golang-licenseserver/server"
)

func main() {
	generateKeys := flag.Bool("generateKeys", false, "Generate public and private keys")
	privateKeyPath := flag.String("privateKeyPath", "private.pem", "Path to the private key file")
	publicKeyPath := flag.String("publicKeyPath", "public.pem", "Path to the public key file")

	generateLicense := flag.Bool("generateLicense", false, "Generate a license file")
	licensePath := flag.String("licensePath", "license.txt", "Path to the license file")
	features := flag.String("features", "", "Comma-separated list of features for the license")
	duration := flag.Int("duration", 12, "Duration of the license in months")

	flag.Parse()

	if *generateKeys {
		fmt.Println("Generate public and private keys")

		err := server.GenerateRsaKeyPair(*privateKeyPath, *publicKeyPath)
		if err != nil {
			fmt.Println("Error generating keys:", err)
			os.Exit(1)
		}
	} else if *generateLicense {
		licenseServer := &server.LicenseServer{}
		licenseServer.Init(*privateKeyPath)

		license, err := licenseServer.GenerateLicense(strings.Split(*features, ","), *duration)
		if err != nil {
			fmt.Println("Error generating license:", err)
			os.Exit(1)
		}

		// save the license to a file
		licenseFile, err := os.Create(*licensePath)
		if err != nil {
			fmt.Println("Error creating license file:", err)
			os.Exit(1)
		}
		defer licenseFile.Close()
		_, err = licenseFile.WriteString(license)
		if err != nil {
			fmt.Println("Error writing license to file:", err)
			os.Exit(1)
		}
	}
}
