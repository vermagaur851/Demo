package datadogbackend

import (
	"amantya_metrics/metricsInterface"
	"fmt"
	"log"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type DataDogMetric struct {
	client *statsd.Client
	name   string
}

func (dm *DataDogMetric) Inc(labels map[string]string) error {
	tags := convertLabelsToTags(labels)
	log.Printf("Datadog: Incrementing metric %s with labels %v", dm.name, labels)

	return dm.client.Incr(dm.name, tags, 1)
}

func (dm *DataDogMetric) Dec(labels map[string]string) error {
	tags := convertLabelsToTags(labels)
	return dm.client.Decr(dm.name, tags, 1)
}

func (dm *DataDogMetric) Add(value float64, labels map[string]string) error {
	tags := convertLabelsToTags(labels)
	return dm.client.Count(dm.name, int64(value), tags, 1)
}

func (dm *DataDogMetric) Set(value float64, labels map[string]string) error {
	tags := convertLabelsToTags(labels)
	log.Printf("Datadog: Setting metric '%s' = %f with labels %v", dm.name, value, labels)
	return dm.client.Gauge(dm.name, value, tags, 1)
}

func (dm *DataDogMetric) Observe(value float64, labels map[string]string) error {
	return metricsInterface.ErrInvalidOperation
}

func (dm *DataDogMetric) GetMetricType() metricsInterface.MetricType {
	// DataDog doesn't strictly distinguish between counter and gauge in the client
	// We'll treat all metrics as counters unless they use Set
	return metricsInterface.CounterType
}

func convertLabelsToTags(labels map[string]string) []string {
	tags := make([]string, 0, len(labels))
	for k, v := range labels {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}
	return tags
}

type DataDogBackend struct {
	client    *statsd.Client
	namespace string
}

func NewDataDogBackend(namespace string) (*DataDogBackend, error) {
	// You can configure the address based on your environment
	// Default is "localhost:8125" for UDP
	client, err := statsd.New("localhost:8125",
		statsd.WithNamespace(namespace),
		statsd.WithTags([]string{fmt.Sprintf("service:%s", os.Getenv("DD_SERVICE"))}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create DataDog client: %w", err)
	}

	return &DataDogBackend{
		client:    client,
		namespace: namespace,
	}, nil
}

func (db *DataDogBackend) NewCounter(name, help string, labels []string) (metricsInterface.Metric, error) {
	return &DataDogMetric{
		client: db.client,
		name:   name,
	}, nil
}

func (db *DataDogBackend) NewGauge(name, help string, labels []string) (metricsInterface.Metric, error) {
	return &DataDogMetric{
		client: db.client,
		name:   name,
	}, nil
}

func (db *DataDogBackend) NewHistogram(name, help string, labels []string, buckets []float64) (metricsInterface.Metric, error) {
	return nil, metricsInterface.ErrBackendNotSupported
}

func (db *DataDogBackend) PushToGateway(gatewayURL, jobName string) error {
	// DataDog doesn't use Pushgateway, metrics are sent directly
	// Flush any buffered metrics
	return db.client.Flush()
}
