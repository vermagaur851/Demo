package main

/*
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import (
	"amantya_metrics/metricsInterface"
	"amantya_metrics/metrics_wrapper"
	"amantya_metrics/models"
	"fmt"
	"log"
	"strings"
	"unsafe"
)

var framework *metrics_wrapper.MetricsFramework

//export Initialize
func Initialize(backendType *C.char, namespace *C.char) C.int {
	if backendType != nil {
		fmt.Println("Backend:", C.GoString(backendType))
	}
	if namespace != nil {
		fmt.Println("Namespace:", C.GoString(namespace))
	}

	options := make(map[string]interface{})
	if namespace != nil {
		options["namespace"] = C.GoString(namespace)
	}

	f, err := metrics_wrapper.MetricsType(metrics_wrapper.BackendType(C.GoString(backendType)), options)
	if err != nil {
		fmt.Println("Failed to create MetricsFramework:", err)
		return -1
	}

	fmt.Println("MetricsFramework initialized successfully")
	framework = f
	return 0
}

//export LoadKPIs
func LoadKPIs(filePath *C.char) C.int {
	if framework == nil {
		fmt.Println("LoadKPIs called before successful Initialize")
		return -1
	}
	err := framework.LoadKPIs(C.GoString(filePath))
	if err != nil {
		fmt.Println("Failed to load KPIs:", err)
		return -1
	}
	fmt.Println("KPIs loaded from:", C.GoString(filePath))
	return 0
}

//export RegisterMetrics
func RegisterMetrics() C.int {
	err := framework.RegisterMetrics()
	if err != nil {
		return -1
	}
	return 0
}

//export IncrementMetric
func IncrementMetric(metricName *C.char, labels **C.char, count C.int) C.int {
	goLabels := make(map[string]string)

	// Convert C array of strings to Go map
	if labels != nil && count > 0 {
		cLabels := (*[1 << 30]*C.char)(unsafe.Pointer(labels))[:count:count]
		for i := 0; i < int(count); i += 2 {
			if i+1 >= int(count) {
				break
			}
			key := C.GoString(cLabels[i])
			value := C.GoString(cLabels[i+1])
			goLabels[key] = value
		}
	}

	// Normalize metric name
	normalizedName := normalizeMetricName(C.GoString(metricName))

	// Get the metric
	metric, err := framework.GetMetric(normalizedName)
	if err != nil {
		log.Printf("Metric not found: %s (normalized from: %s)", normalizedName, C.GoString(metricName))
		return -1
	}

	// Validate labels against KPI requirements
	if err := validateLabels(normalizedName, goLabels); err != nil {
		log.Printf("Label validation failed for %s: %v", normalizedName, err)
		return -1
	}

	if err := metric.Inc(goLabels); err != nil {
		log.Printf("Increment failed for %s: %v", normalizedName, err)
		return -1
	}

	return 0
}

func validateLabels(metricName string, labels map[string]string) error {
	// Get KPIs through the public method
	kpis := framework.GetKPIs()

	// Find the KPI definition
	var kpi *models.KPI
	for _, k := range kpis {
		if normalizeMetricName(k.DisplayName) == metricName {
			kpi = &k
			break
		}
	}
	if kpi == nil {
		return fmt.Errorf("KPI definition not found")
	}

	// Check required labels
	for _, requiredLabel := range kpi.Object {
		if _, exists := labels[requiredLabel]; !exists {
			return fmt.Errorf("missing required label: %s", requiredLabel)
		}
	}

	return nil
}

func normalizeMetricName(displayName string) string {
	name := strings.ToLower(displayName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	return name
}

//export DecrementMetric
func DecrementMetric(metricName *C.char, labels **C.char, count C.int) C.int {
	goLabels := make(map[string]string)
	cLabels := (*[1 << 30]*C.char)(unsafe.Pointer(labels))[:count:count]

	for i := 0; i < int(count); i += 2 {
		if i+1 >= int(count) {
			break
		}
		key := C.GoString(cLabels[i])
		value := C.GoString(cLabels[i+1])
		goLabels[key] = value
	}

	if err := framework.DecrementMetric(C.GoString(metricName), goLabels); err != nil {
		return -1
	}

	return 0
}

//export AddToMetric
func AddToMetric(metricName *C.char, value C.double, labels **C.char, count C.int) C.int {
	goLabels := make(map[string]string)
	cLabels := (*[1 << 30]*C.char)(unsafe.Pointer(labels))[:count:count]

	for i := 0; i < int(count); i += 2 {
		if i+1 >= int(count) {
			break
		}
		key := C.GoString(cLabels[i])
		value := C.GoString(cLabels[i+1])
		goLabels[key] = value
	}

	if err := framework.AddToMetric(C.GoString(metricName), float64(value), goLabels); err != nil {
		return -1
	}

	return 0
}

//export SetMetric
func SetMetric(metricName *C.char, value C.double, labels **C.char, count C.int) C.int {
	// Convert C inputs to Go types
	name := C.GoString(metricName)
	goLabels := make(map[string]string)

	// Process labels if provided
	if labels != nil && count > 0 {
		cLabels := (*[1 << 30]*C.char)(unsafe.Pointer(labels))[:count:count]
		for i := 0; i < int(count); i += 2 {
			if i+1 >= int(count) {
				break
			}
			key := C.GoString(cLabels[i])
			value := C.GoString(cLabels[i+1])
			goLabels[key] = value
		}
	}

	// Get the metric
	metric, err := framework.GetMetric(name)
	if err != nil {
		log.Printf("SetMetric failed: %v", err)
		return -1
	}

	// Check if metric supports Set operation
	if metric.GetMetricType() == metricsInterface.CounterType {
		log.Printf("Set operation not supported for Counter metric: %s", name)
		return -1
	}

	// Perform the set operation
	if err := metric.Set(float64(value), goLabels); err != nil {
		log.Printf("SetMetric failed: %v", err)
		return -1
	}

	return 0
}

//export PushMetrics
func PushMetrics(gatewayURL *C.char, jobName *C.char) C.int {
	if err := framework.PushMetrics(C.GoString(gatewayURL), C.GoString(jobName)); err != nil {
		return -1
	}
	return 0
}

//export ListMetrics
func ListMetrics() **C.char {
	metrics := framework.ListMetrics()

	// Allocate array of C strings
	cArray := C.malloc(C.size_t(len(metrics)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	a := (*[1 << 30]*C.char)(cArray)
	for i, metric := range metrics {
		a[i] = C.CString(metric)
	}

	return (**C.char)(cArray)
}

//export FreeStringArray
func FreeStringArray(array **C.char, length C.int) {
	cArray := (*[1 << 30]*C.char)(unsafe.Pointer(array))[:length:length]

	for i := 0; i < int(length); i++ {
		C.free(unsafe.Pointer(cArray[i]))
	}

	C.free(unsafe.Pointer(array))
}

func main() {}
