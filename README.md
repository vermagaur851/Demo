# Amantya Metrics Framework (Version v1.0.0)

The **Amantya Metrics Framework** is a modular, backend-agnostic metric management system built in **Go**. It simplifies the registration, management, and interaction with various types of metrics such as **Counters** and **Gauges**, and supports seamless integration with observability backends like **Prometheus** and **Datadog**.

This example demonstrates using the framework with Prometheus, loading KPIs from a JSON file, operating on metrics, and pushing them to a Prometheus **Pushgateway**.

---

## Project Overview

- **Backend Agnostic**: Works with Prometheus, Datadog, and other metric backends.
- **Dynamic KPI Integration**: Load KPIs from a JSON file.
- **Metric Operations**: Supports `increment`, `decrement`, `add`, and `set` operations.
- **Default Initialization**: No need to pre-register metrics manually.
- **PushGateway Support**: Pushes metrics to Prometheus PushGateway for scraping.
- **Metric Listing**: Easily list all registered metrics.

---

## Features

- **Backend Support**: Prometheus, Datadog  
- **KPI JSON Loading**  
- **Metric APIs**: Increment, Decrement, Add, Set  
- **Push to Prometheus Pushgateway**  
- **Dynamic Metric Listing**  
- **Plug-and-Play Architecture**

---

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/dl/) (v1.18 or later)
- [Prometheus Pushgateway](https://github.com/prometheus/pushgateway)

> You should also have a basic understanding of Go and Prometheus.

---

### Installation

**Clone the repository:**
   ```path
   https://bitbucket.org/amantyatech/cg_5gc/src/integration/5gCN/common/goCommon/src/amantya_metrics
   ```
 
**Configure Prometheus:**
    If you are using Prometheus for metrics, make sure it's running and properly configured to collect data from your project. You can follow the Prometheus setup guide if it's not set up yet.

This will:

- Load KPIs from kpi.json

- Register and initialize metrics

- Perform operations (add, set, increment, etc.)

- Push metrics to Prometheus PushGateway (http://localhost:9091)

## Prject structure 
``` bash
AMANTYA_METRICS/
├── build/
├── datadogbackend/
├── lang_wrapper/
├── metricsInterface/
├── metricsregistry/
├── models/
├── prometheusbackend/
├── service/
├── utils/              
├── Makefile
├── go.mod / go.sum
└── README.md
```

## Sample KPI Configuration (kpi.json)
```bash
[
  {
    "name": "mean_registered_subscribers_amf",
    "displayName": "Mean Registered Subscribers (AMF)",
    "description": "Average number of registered subscribers in the network and network slice handled by AMF.",
    "formula": "Sum of registered subscribers / Observation period",
    "unit": "Count",
    "type": "Mean",
    "object": ["NetworkSlice", "Network"],
    "prometheus_type": "Counter",
    "nf_type": "AMF",
    "increment": true,
    "decrement": false
  },
  {
    "name": "registration_success_rate_single_slice",
    "displayName": "Registration Success Rate (Single Slice)",
    "description": "Success rate of user equipment (UE) registration for a single network slice.",
    "formula": "(Number of successful registrations / Total number of registrations) * 100",
    "unit": "%",
    "type": "SuccessRate",
    "object": ["NetworkSlice"],
    "prometheus_type": "Gauge",
    "nf_type": "AMF",
    "increment": true,
    "decrement": true
  }
]
```

## Example Metric Operations

```bash
framework.IncrementMetric("mean_registered_subscribers_amf", map[string]string{"NetworkSlice": "slice1"})
framework.AddToMetric("registered_subscribers_udm", 5, map[string]string{"Network": "net1"})
framework.SetMetric("registration_success_rate_single_slice", 95.5, map[string]string{"NetworkSlice": "slice1"})
framework.DecrementMetric("registration_success_rate_single_slice", map[string]string{"NetworkSlice": "slice1"})
```

## 🧰 API Reference
- RegisterMetrics()
    Registers all metrics defined in your KPI configuration.
    Example:
    framework.RegisterMetrics()

- GetMetric(name string)
    Retrieves a previously registered metric by its name. Returns the metric object (or an error if not found).
    Example:
    metric, err := framework.GetMetric("api_requests")

- ListMetrics() []string
    Returns a slice of all registered metric names.
    Example:
    metrics := framework.ListMetrics()

- UnregisterMetric(name string)
    Removes/unregisters the metric with the given name.
    Example:
    framework.UnregisterMetric("old_metric")

- PushMetrics(gatewayURL, jobName string)
    Pushes all current metric values to a Prometheus PushGateway at the specified URL under the given job name.
    Example:
    framework.PushMetrics("http://pushgateway:9091", "my_job")


## Push Output
```bash 
    http://localhost:9091/metrics 
```

## Author

Amantya

## License

This project is licensed under the MIT License.