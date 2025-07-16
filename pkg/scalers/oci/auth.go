package oci

import (
	"github.com/oracle/oci-go-sdk/v65/common"
)

// GetOCISigner returns a signer based on instance principal auth
func GetOCISigner() (common.ConfigurationProvider, error) {
	return common.DefaultConfigProvider(), nil
}
