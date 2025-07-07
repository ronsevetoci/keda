package scalers

import (
	"context"
	"fmt"

	"github.com/kedacore/keda/v2/pkg/scalers/oci"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	external_metrics "k8s.io/metrics/pkg/apis/external_metrics"
)

type ociQueueScaler struct {
	metadata       *oci.QueueMetadata
	client         *oci.QueueClient
	scalerMetadata *ScalerMetadata
}

func NewOCIQueueScaler(config *ScalerConfig) (Scaler, error) {
	meta, err := oci.ParseMetadata(config)
	if err != nil {
		return nil, fmt.Errorf("error parsing metadata: %s", err)
	}

	client, err := oci.NewQueueClient(meta)
	if err != nil {
		return nil, fmt.Errorf("error creating queue client: %s", err)
	}

	return &ociQueueScaler{
		metadata:       meta,
		client:         client,
		scalerMetadata: &config.ScalerMetadata,
	}, nil
}

func (s *ociQueueScaler) IsActive(ctx context.Context) (bool, error) {
	stats, err := s.client.Client.GetStats(ctx, s.metadata.QueueID)
	if err != nil {
		return false, err
	}

	total := *stats.Queue.VisibleMessages + *stats.Queue.InFlightMessages
	return total > 0, nil
}

func (s *ociQueueScaler) GetMetricSpecForScaling(context.Context) []external_metrics.ExternalMetricValue {
	return []external_metrics.ExternalMetricValue{{
		MetricName: oci.MetricName,
		Value:      *resource.NewQuantity(int64(s.metadata.QueueLength), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}}
}

func (s *ociQueueScaler) GetMetrics(ctx context.Context, metricName string, metricSelector labels.Selector) ([]external_metrics.ExternalMetricValue, error) {
	stats, err := s.client.Client.GetStats(ctx, s.metadata.QueueID)
	if err != nil {
		return nil, err
	}

	total := *stats.Queue.VisibleMessages + *stats.Queue.InFlightMessages
	metric := external_metrics.ExternalMetricValue{
		MetricName: oci.MetricName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}

	return []external_metrics.ExternalMetricValue{metric}, nil
}

func (s *ociQueueScaler) Close(context.Context) error {
	return nil
}
