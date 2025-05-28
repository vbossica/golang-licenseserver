package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vbossica/golang-licenseserver/client"
)

func main() {
	verifyFeature := flag.Bool("verifyFeature", false, "Verify that a feature is supported in the license")
	publicKeyPath := flag.String("publicKeyPath", "public.pem", "Path to the public key file")
	licensePath := flag.String("licensePath", "license.txt", "Path to the license file")
	feature := flag.String("feature", "", "Features to validate in the license")

	flag.Parse()

	if *verifyFeature {
		licenseClient := &client.LicenseClient{}
		err := licenseClient.Init(*publicKeyPath, *licensePath)
		if err != nil {
			log.Fatal("Error writing instantiating the LicenseClient:", err)
		}

		isValid, err := licenseClient.VerifyLicense(*feature)
		if err != nil {
			fmt.Println("Error verifying license:", err)
			os.Exit(1)
		}

		fmt.Println("Is claim valid:", isValid)
	}
}
