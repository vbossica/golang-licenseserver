package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLicense(t *testing.T) {
	licenseServer := createLicenseServer(t)

	features := []string{"feature-1", "feature-2"}
	duration := 12

	license, err := licenseServer.GenerateLicense(features, duration)
	assert.Nil(t, err)
	assert.NotEmpty(t, license, "Generated license should not be empty")
}

func createLicenseServer(t *testing.T) *LicenseServer {
	licenseServer := &LicenseServer{}
	err := licenseServer.Init("testdata/private.pem")
	assert.Nil(t, err)

	return licenseServer
}
