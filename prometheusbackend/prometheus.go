package prometheusbackend

import (
	"amantya_metrics/metricsInterface"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PrometheusCounter struct {
	counter *prometheus.CounterVec
}

func (pc *PrometheusCounter) Inc(labels map[string]string) error {
	pc.counter.With(labels).Inc()
	return nil
}

func (pc *PrometheusCounter) Dec(labels map[string]string) error {
	return metricsInterface.ErrInvalidOperation
}

func (pc *PrometheusCounter) Add(value float64, labels map[string]string) error {
	pc.counter.With(labels).Add(value)
	return nil
}

func (pc *PrometheusCounter) Set(value float64, labels map[string]string) error {
	return metricsInterface.ErrInvalidOperation
}

func (pc *PrometheusCounter) Observe(value float64, labels map[string]string) error {
	return metricsInterface.ErrInvalidOperation
}

func (pc *PrometheusCounter) GetMetricType() metricsInterface.MetricType {
	return metricsInterface.CounterType
}

type PrometheusGauge struct {
	gauge *prometheus.GaugeVec
}

func (pg *PrometheusGauge) Inc(labels map[string]string) error {
	pg.gauge.With(labels).Inc()
	return nil
}

func (pg *PrometheusGauge) Dec(labels map[string]string) error {
	pg.gauge.With(labels).Dec()
	return nil
}

func (pg *PrometheusGauge) Add(value float64, labels map[string]string) error {
	pg.gauge.With(labels).Add(value)
	return nil
}

func (pg *PrometheusGauge) Set(value float64, labels map[string]string) error {
	pg.gauge.With(labels).Set(value)
	return nil
}

func (pg *PrometheusGauge) Observe(value float64, labels map[string]string) error {
	return metricsInterface.ErrInvalidOperation
}

func (pg *PrometheusGauge) GetMetricType() metricsInterface.MetricType {
	return metricsInterface.GaugeType
}

type PrometheusBackend struct {
	registry *prometheus.Registry
}

func NewPrometheusBackend() *PrometheusBackend {
	return &PrometheusBackend{
		registry: prometheus.NewRegistry(),
	}
}

func (pb *PrometheusBackend) NewCounter(name, help string, labels []string) (metricsInterface.Metric, error) {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	if err := pb.registry.Register(counter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			existing := are.ExistingCollector.(*prometheus.CounterVec)
			return &PrometheusCounter{counter: existing}, nil
		}
		return nil, err
	}

	return &PrometheusCounter{counter: counter}, nil
}

func (pb *PrometheusBackend) NewGauge(name, help string, labels []string) (metricsInterface.Metric, error) {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	if err := pb.registry.Register(gauge); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			existing := are.ExistingCollector.(*prometheus.GaugeVec)
			return &PrometheusGauge{gauge: existing}, nil
		}
		return nil, err
	}

	return &PrometheusGauge{gauge: gauge}, nil
}

func (pb *PrometheusBackend) NewHistogram(name, help string, labels []string, buckets []float64) (metricsInterface.Metric, error) {
	return nil, metricsInterface.ErrBackendNotSupported
}

func (pb *PrometheusBackend) PushToGateway(gatewayURL, jobName string) error {
	pusher := push.New(gatewayURL, jobName).Gatherer(pb.registry)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	pusher.Client(client)

	if err := pusher.Add(); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	return nil
}
