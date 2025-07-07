package oci

import (
	"github.com/oracle/oci-go-sdk/v65/queue"
)

type QueueClient struct {
	Client queue.QueueClient
}

func NewQueueClient(meta *QueueMetadata) (*QueueClient, error) {
	provider, err := GetOCIAuthProvider(meta.AuthType)
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
