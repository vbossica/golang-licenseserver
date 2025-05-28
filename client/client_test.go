package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidFeature(t *testing.T) {
	licenseClient := createLicenseClient(t)

	valid, err := licenseClient.VerifyFeature("feature-1")
	assert.Nil(t, err)
	assert.True(t, valid)
}

func TestInvalidFeature(t *testing.T) {
	licenseClient := createLicenseClient(t)

	valid, err := licenseClient.VerifyFeature("invalidFeature")
	assert.NotNil(t, err)
	assert.False(t, valid)
}

func createLicenseClient(t *testing.T) *LicenseClient {
	licenseClient := &LicenseClient{}
	err := licenseClient.Init("testdata/public.pem", "testdata/license.txt")
	assert.Nil(t, err)

	return licenseClient
}
