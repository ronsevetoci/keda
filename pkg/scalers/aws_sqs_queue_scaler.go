package scalers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/kedacore/keda/v2/pkg/scalers/scalersconfig"
	"github.com/kedacore/keda/v2/pkg/util"
	"github.com/oracle/oci-go-sdk/v65/queue"
	"github.com/ronsevetoci/keda/v2/pkg/scalers/oci"
)

const ociMetricType = v2.ExternalMetricSourceType

type ociQueueScaler struct {
	metricType  v2.MetricTargetType
	metadata    *oci.OCITriggerMetadata
	queueClient *oci.QueueClient
	logger      logr.Logger
}

func NewOCIQueueScaler(ctx context.Context, config *scalersconfig.ScalerConfig) (Scaler, error) {
	metricType, err := GetMetricTargetType(config)
	if err != nil {
		return nil, err
	}

	meta, err := oci.ParseOCITriggerMetadata(config)
	if err != nil {
		return nil, fmt.Errorf("error parsing OCI metadata: %w", err)
	}

	client, err := oci.NewQueueClient(meta)
	if err != nil {
		return nil, fmt.Errorf("error creating OCI client: %w", err)
	}

	logger := InitializeLogger(config, "oci_queue_scaler")

	return &ociQueueScaler{
		metricType:  metricType,
		metadata:    meta,
		queueClient: client,
		logger:      logger,
	}, nil
}

func (s *ociQueueScaler) Close(context.Context) error {
	return nil
}

func (s *ociQueueScaler) GetMetricSpecForScaling(_ context.Context) []v2.MetricSpec {
	externalMetric := &v2.ExternalMetricSource{
		Metric: v2.MetricIdentifier{
			Name: util.GenerateMetricNameWithIndex(s.metadata.TriggerIndex, "oci-queue"),
		},
		Target: GetMetricTarget(s.metricType, s.metadata.QueueLength),
	}
	return []v2.MetricSpec{
		{
			External: externalMetric,
			Type:     ociMetricType,
		},
	}
}

func (s *ociQueueScaler) GetMetricsAndActivity(ctx context.Context, metricName string) ([]external_metrics.ExternalMetricValue, bool, error) {
	stats, err := s.queueClient.Client.GetStats(ctx, queue.GetStatsRequest{QueueId: &s.metadata.QueueID})
	if err != nil {
		s.logger.Error(err, "failed to fetch stats")
		return nil, false, err
	}

	visible := *stats.Queue.VisibleMessages
	inflight := *stats.Queue.InFlightMessages
	total := visible + inflight

	metric := GenerateMetricInMili(metricName, float64(total))
	return []external_metrics.ExternalMetricValue{metric}, total > int32(s.metadata.ActivationQueueLength), nil
}
