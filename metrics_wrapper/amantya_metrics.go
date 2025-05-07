package metrics_wrapper

import (
	"amantya_metrics/metricsInterface"
	"amantya_metrics/metricsregistry"
	"amantya_metrics/models"
	"amantya_metrics/prometheusbackend"
	"fmt"
	"log"
	"strings"
)

type BackendType string

const (
	PrometheusBackend BackendType = "prometheus"
	DataDogBackend    BackendType = "datadog"
)

type MetricsFramework struct {
	registry *metricsregistry.Registry
	backend  metricsInterface.Backend
	kpIs     []models.KPI
}

func MetricsType(backendType BackendType, options map[string]interface{}) (*MetricsFramework, error) {
	var backend metricsInterface.Backend
	var err error

	switch backendType {
	case PrometheusBackend:
		backend = prometheusbackend.NewPrometheusBackend() // Initialize properly
	// case DataDogBackend:
	//     namespace, ok := options["namespace"].(string)
	//     if !ok {
	//         namespace = "amantya"
	//     }
	//     backend, err = datadogbackend.NewDataDogBackend(namespace)
	default:
		return nil, fmt.Errorf("%w: %s", metricsInterface.ErrBackendNotSupported, backendType)
	}

	if err != nil {
		return nil, err
	}

	return &MetricsFramework{
		registry: metricsregistry.NewRegistry(),
		backend:  backend,
	}, nil
}

func (mf *MetricsFramework) LoadKPIs(filePath string) error {
	kpIs, err := models.LoadKPIsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load KPIs: %w", err)
	}

	mf.kpIs = kpIs
	return nil
}
func (m *MetricsFramework) Backend() interface{} {
	return m.backend
}

func (m *MetricsFramework) GetKPIs() []models.KPI {
	return m.kpIs
}

func normalizeMetricName(displayName string) string {
	name := strings.ToLower(displayName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
		name = "g" + name
	}
	return name
}
func (mf *MetricsFramework) RegisterMetrics() error {
	for _, kpi := range mf.kpIs {
		// Ensure consistent naming
		metricName := normalizeMetricName(kpi.DisplayName)

		// Check if metric already exists
		if _, err := mf.registry.Get(metricName); err == nil {
			continue // Skip if already registered
		}

		var metric metricsInterface.Metric
		var err error

		switch kpi.PrometheusType {
		case "Counter":
			metric, err = mf.backend.NewCounter(metricName, kpi.Description, kpi.Object)
		case "Gauge":
			metric, err = mf.backend.NewGauge(metricName, kpi.Description, kpi.Object)
		default:
			return fmt.Errorf("unsupported metric type: %s", kpi.PrometheusType)
		}

		if err != nil {
			return fmt.Errorf("failed to create metric %s: %w", metricName, err)
		}

		if err := mf.registry.Register(metricName, metric); err != nil {
			return fmt.Errorf("failed to register metric %s: %w", metricName, err)
		}
		log.Printf("Successfully registered metric: %s", metricName)
	}
	return nil
}

func (mf *MetricsFramework) GetMetric(name string) (metricsInterface.Metric, error) {
	return mf.registry.Get(name)
}

func (mf *MetricsFramework) IncrementMetric(name string, labels map[string]string) error {
	metric, err := mf.registry.Get(name)
	if err != nil {
		return err
	}

	return metric.Inc(labels)
}

func (mf *MetricsFramework) DecrementMetric(name string, labels map[string]string) error {
	metric, err := mf.registry.Get(name)
	if err != nil {
		return err
	}

	return metric.Dec(labels)
}

func (mf *MetricsFramework) AddToMetric(name string, value float64, labels map[string]string) error {
	metric, err := mf.registry.Get(name)
	if err != nil {
		return err
	}

	return metric.Add(value, labels)
}

func (mf *MetricsFramework) SetMetric(name string, value float64, labels map[string]string) error {
	metric, err := mf.registry.Get(name)
	if err != nil {
		return err
	}

	return metric.Set(value, labels)
}

func (mf *MetricsFramework) PushMetrics(gatewayURL, jobName string) error {
	return mf.backend.PushToGateway(gatewayURL, jobName)
}

func (mf *MetricsFramework) ListMetrics() []string {
	return mf.registry.List()
}

func (mf *MetricsFramework) UnregisterMetric(name string) error {
	return mf.registry.Unregister(name)
}

func PushWithDefaults(mf *MetricsFramework, gatewayURL, jobName string) error {
	if err := mf.InitializeDefaults(); err != nil {
		return fmt.Errorf("failed to initialize defaults: %w", err)
	}
	return mf.PushMetrics(gatewayURL, jobName)
}

func createDefaultLabels(objects []string) map[string]string {
	labels := make(map[string]string)
	for _, obj := range objects {
		labels[obj] = "default_" + obj
	}
	return labels
}

// InitializeDefaults sets zero values for all registered metrics
func (mf *MetricsFramework) InitializeDefaults() error {
	for _, kpi := range mf.kpIs {
		metricName := normalizeMetricName(kpi.DisplayName)
		labels := createDefaultLabels(kpi.Object)

		switch kpi.PrometheusType {
		case "Counter":
			if err := mf.AddToMetric(metricName, 0, labels); err != nil {
				return fmt.Errorf("failed to initialize counter %s: %w", metricName, err)
			}
		case "Gauge":
			if err := mf.SetMetric(metricName, 0, labels); err != nil {
				return fmt.Errorf("failed to initialize gauge %s: %w", metricName, err)
			}
		default:
			log.Printf("Skipping initialization for unknown metric type %s (%s)",
				kpi.PrometheusType, metricName)
		}
	}
	return nil
}
