package metricsregistry

import (
	"amantya_metrics/metricsInterface"
	"sync"
)

type Registry struct {
	metrics map[string]metricsInterface.Metric
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: make(map[string]metricsInterface.Metric),
	}
}

func (r *Registry) Register(name string, metric metricsInterface.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.metrics[name]; exists {
		return metricsInterface.ErrMetricAlreadyRegistered
	}

	r.metrics[name] = metric
	return nil
}

func (r *Registry) Get(name string) (metricsInterface.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if metric, exists := r.metrics[name]; exists {
		return metric, nil
	}

	return nil, metricsInterface.ErrMetricNotFound
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.metrics))
	for name := range r.metrics {
		names = append(names, name)
	}

	return names
}

func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.metrics[name]; !exists {
		return metricsInterface.ErrMetricNotFound
	}

	delete(r.metrics, name)
	return nil
}
