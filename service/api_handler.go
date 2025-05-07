package service

import (
	"amantya_metrics/metrics_wrapper"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type MetricRequest struct {
	Name   string            `json:"name"`
	Value  float64           `json:"value,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type APIHandler struct {
	framework *metrics_wrapper.MetricsFramework
}

func NewAPIHandler(framework *metrics_wrapper.MetricsFramework) *APIHandler {
	return &APIHandler{framework: framework}
}

// Consistent normalization function used everywhere
func normalizeMetricName(displayName string) string {
	name := strings.ToLower(displayName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	return name
}

func (h *APIHandler) RegisterMetrics(w http.ResponseWriter, r *http.Request) {
	// Unregister existing metrics
	for _, name := range h.framework.ListMetrics() {
		if err := h.framework.UnregisterMetric(name); err != nil {
			log.Printf("Failed to unregister metric %s: %v", name, err)
		}
	}

	// Register fresh metrics
	if err := h.framework.RegisterMetrics(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) IncrementMetric(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricName := normalizeMetricName(req.Name)
	log.Printf("Attempting to increment metric: %s", metricName)

	metric, err := h.framework.GetMetric(metricName)
	if err != nil {
		log.Printf("Metric not found: %s (searched as: %s)", req.Name, metricName)
		http.Error(w, "metric not found: "+metricName, http.StatusNotFound)
		return
	}

	if err := metric.Inc(req.Labels); err != nil {
		log.Printf("Increment failed for %s: %v", metricName, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Successfully incremented metric: %s", metricName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) DecrementMetric(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricName := normalizeMetricName(req.Name)
	log.Printf("Decrementing metric: %s", metricName)

	metric, err := h.framework.GetMetric(metricName)
	if err != nil {
		http.Error(w, "metric not found: "+metricName, http.StatusNotFound)
		return
	}

	if err := metric.Dec(req.Labels); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) AddToMetric(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Value  float64           `json:"value"`
		Labels map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricName := normalizeMetricName(req.Name)
	log.Printf("Adding to metric: %s", metricName)

	metric, err := h.framework.GetMetric(metricName)
	if err != nil {
		http.Error(w, "metric not found: "+metricName, http.StatusNotFound)
		return
	}

	if err := metric.Add(req.Value, req.Labels); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) SetMetric(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Value  float64           `json:"value"`
		Labels map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricName := normalizeMetricName(req.Name)
	log.Printf("Setting metric: %s", metricName)

	metric, err := h.framework.GetMetric(metricName)
	if err != nil {
		http.Error(w, "metric not found: "+metricName, http.StatusNotFound)
		return
	}

	if err := metric.Set(req.Value, req.Labels); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) PushMetrics(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GatewayURL string `json:"gateway_url"`
		JobName    string `json:"job_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.framework.PushMetrics(req.GatewayURL, req.JobName); err != nil {
		log.Printf("Push to gateway failed: %v", err)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "Metrics collected but push failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.framework.ListMetrics()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"count":   len(metrics),
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (h *APIHandler) DebugMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.framework.ListMetrics()
	details := make(map[string]interface{})

	for _, name := range metrics {
		metric, err := h.framework.GetMetric(name)
		if err != nil {
			details[name] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
			continue
		}

		details[name] = map[string]interface{}{
			"type":   metric.GetMetricType(),
			"status": "ok",
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(details)
}
