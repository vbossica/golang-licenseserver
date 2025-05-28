package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vbossica/golang-licenseserver/client"
	"github.com/vbossica/golang-licenseserver/server"
)

func main() {
	err := server.GenerateRsaKeyPair("server/testdata/private.pem", "client/testdata/public.pem")
	if err != nil {
		fmt.Printf("Cannot generate RSA key pair: %s\n", err)
		os.Exit(1)
	}

	// Create a license server
	licenseServer := &server.LicenseServer{}
	licenseServer.Init("server/testdata/private.pem")

	// generate the license for features and duration (in months)
	features := []string{"feature-1", "feature-2"}
	license, err := licenseServer.GenerateLicense(features, 12)
	if err != nil {
		log.Fatal("Error creating license:", err)
	}

	fmt.Println("License to be saved and used on the device:", license)
	// save the license to a file
	licenseFile, err := os.Create("client/testdata/license.txt")
	if err != nil {
		log.Fatal("Error creating license file:", err)
	}
	defer licenseFile.Close()

	_, err = licenseFile.WriteString(license)
	if err != nil {
		log.Fatal("Error writing license to file:", err)
	}
	fmt.Println("License saved to license.txt")

	licenseClient := &client.LicenseClient{}
	err = licenseClient.Init("client/testdata/public.pem", "client/testdata/license.txt")
	if err != nil {
		log.Fatal("Error writing instantiating the LicenseClient:", err)
	}

	// verify the claim on the device
	isValid, _ := licenseClient.VerifyLicense("feature-1")
	fmt.Println("Is claim valid:", isValid)

	isValid, _ = licenseClient.VerifyLicense("feature-3")
	fmt.Println("Is claim valid:", isValid)
}
