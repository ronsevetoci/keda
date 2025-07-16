package oci

import (
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/queue"
)

type QueueClient struct {
	Client queue.QueueClient
}

func NewQueueClient(meta *OCITriggerMetadata) (*QueueClient, error) {
	provider, err := common.DefaultConfigProvider()
	if err != nil {
		return nil, err
	}

	client, err := queue.NewQueueClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	client.SetRegion(meta.Region)
	if meta.Endpoint != "" {
		client.BaseClient.Host = meta.Endpoint
	}

	return &QueueClient{Client: client}, nil
}
