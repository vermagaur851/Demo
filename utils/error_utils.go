package utils

import (
	"errors"
	"fmt"
)

type MetricError struct {
	Operation string
	Metric    string
	Err       error
}

func (e *MetricError) Error() string {
	return fmt.Sprintf("metric error: %s on %s: %v", e.Operation, e.Metric, e.Err)
}

func (e *MetricError) Unwrap() error {
	return e.Err
}

func NewMetricError(operation, metric string, err error) *MetricError {
	return &MetricError{
		Operation: operation,
		Metric:    metric,
		Err:       err,
	}
}

func IsMetricNotFound(err error) bool {
	return errors.Is(err, errors.New("metric not found"))
}

func IsInvalidOperation(err error) bool {
	return errors.Is(err, errors.New("invalid operation for metric type"))
}
