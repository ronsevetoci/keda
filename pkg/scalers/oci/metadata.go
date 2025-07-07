package oci

import (
	"fmt"
	"strconv"

	"github.com/kedacore/keda/v2/pkg/scalers"
)

type QueueMetadata struct {
	QueueID     string
	Region      string
	QueueLength int
	Endpoint    string
	AuthType    string
}

func ParseMetadata(config *scalers.ScalerConfig) (*QueueMetadata, error) {
	meta := &QueueMetadata{}

	if val, ok := config.TriggerMetadata["queueId"]; ok && val != "" {
		meta.QueueID = val
	} else {
		return nil, fmt.Errorf("no queueId given")
	}

	if val, ok := config.TriggerMetadata["region"]; ok && val != "" {
		meta.Region = val
	} else {
		return nil, fmt.Errorf("no region given")
	}

	if val, ok := config.TriggerMetadata["queueLength"]; ok && val != "" {
		l, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid queueLength")
		}
		meta.QueueLength = l
	} else {
		meta.QueueLength = DefaultQueueLength
	}

	meta.Endpoint = config.TriggerMetadata["endpoint"]
	meta.AuthType = config.TriggerMetadata["authType"]

	return meta, nil
}
