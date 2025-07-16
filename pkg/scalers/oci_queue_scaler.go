package scalers

import (
	"context"
	"fmt"

	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/kedacore/keda/v2/pkg/scalers/scalersconfig"
	kedautil "github.com/kedacore/keda/v2/pkg/util"
	"github.com/ronsevetoci/keda/v2/pkg/scalers/oci"
)

const (
	ociQueueMetricType = "External"
)

type ociQueueScaler struct {
	metadata         *oci.OCITriggerMetadata
	queueClient      *oci.QueueClient
	scalerMetricName string
	metricType       v2.MetricTargetType
}

func NewOCIQueueScaler(config *scalersconfig.ScalerConfig) (Scaler, error) {
	meta, err := oci.ParseOCITriggerMetadata(config.TriggerMetadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing OCI trigger metadata: %s", err)
	}

	client, err := oci.NewQueueClient(meta)
	if err != nil {
		return nil, fmt.Errorf("error creating OCI queue client: %s", err)
	}

	metricName := kedautil.NormalizeString(fmt.Sprintf("oci-queue-%s", meta.QueueID))

	return &ociQueueScaler{
		metadata:         meta,
		queueClient:      client,
		scalerMetricName: GenerateMetricNameWithIndex(config.TriggerIndex, metricName),
		metricType:       v2.AverageValue,
	}, nil
}

func (s *ociQueueScaler) IsActive(ctx context.Context) (bool, error) {
	stats, err := s.queueClient.GetStats(ctx, s.metadata.QueueID)
	if err != nil {
		return false, err
	}

	visible := *stats.Queue.VisibleMessages
	inflight := *stats.Queue.InFlightMessages

	return (visible + inflight) > 0, nil
}

func (s *ociQueueScaler) Close(context.Context) error {
	return nil
}

func (s *ociQueueScaler) GetMetricSpecForScaling(context.Context) []v2.MetricSpec {
	externalMetric := &v2.ExternalMetricSource{
		Metric: v2.MetricIdentifier{
			Name: s.scalerMetricName,
		},
		Target: GetMetricTarget(s.metricType, s.metadata.QueueLength),
	}

	return []v2.MetricSpec{
		{
			Type:     v2.ExternalMetricSourceType,
			External: externalMetric,
		},
	}
}

func (s *ociQueueScaler) GetMetricsAndActivity(ctx context.Context, metricName string, _ labels.Selector) ([]external_metrics.ExternalMetricValue, bool, error) {
	stats, err := s.queueClient.GetStats(ctx, s.metadata.QueueID)
	if err != nil {
		return nil, false, err
	}

	visible := *stats.Queue.VisibleMessages
	inflight := *stats.Queue.InFlightMessages
	total := visible + inflight

	metric := GenerateMetricInMili(metricName, float64(total))
	return []external_metrics.ExternalMetricValue{metric}, total > 0, nil
}
