package metricsInterface

import "errors"

type MetricType string

const (
	CounterType MetricType = "Counter"
	GaugeType   MetricType = "Gauge"
)

type Backend interface {
	NewCounter(name, help string, labels []string) (Metric, error)
	NewGauge(name, help string, labels []string) (Metric, error)
	NewHistogram(name, help string, labels []string, buckets []float64) (Metric, error)
	PushToGateway(gatewayURL, jobName string) error
}

type Metric interface {
	Inc(labels map[string]string) error
	Dec(labels map[string]string) error
	Add(value float64, labels map[string]string) error
	Set(value float64, labels map[string]string) error
	Observe(value float64, labels map[string]string) error
	GetMetricType() MetricType
}

var (
	ErrMetricAlreadyRegistered = errors.New("metric already registered")
	ErrMetricNotFound          = errors.New("metric not found")
	ErrInvalidOperation        = errors.New("invalid operation for metric type")
	ErrInvalidLabel            = errors.New("invalid label provided")
	ErrBackendNotSupported     = errors.New("backend not supported")
)
