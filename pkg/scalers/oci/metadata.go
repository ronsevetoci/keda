package oci

import (
	"fmt"

	"github.com/kedacore/keda/v2/pkg/scalers/scalersconfig"
)

type OCITriggerMetadata struct {
	QueueID               string
	QueueLength           int64
	ActivationQueueLength int64
	Region                string
	Endpoint              string
	TriggerIndex          int
}

func ParseOCITriggerMetadata(config *scalersconfig.ScalerConfig) (*OCITriggerMetadata, error) {
	meta := &OCITriggerMetadata{}

	if err := config.TypedConfig(meta); err != nil {
		return nil, fmt.Errorf("error parsing OCI metadata: %w", err)
	}

	meta.TriggerIndex = config.TriggerIndex
	return meta, nil
}
