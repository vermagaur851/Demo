# Amantya Metrics Framework

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

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/your-repository.git
   cd your-repository

2. **Install Go dependencies:**
    ```bash
    go mod tidy
    ```

3. **Run Prometheus Pushgateway (if not already running):**

    ```bash
    docker run -d -p 9091:9091 prom/pushgateway
    ```

## **Run the Application**

    ```bash
    go run main.go
    ```

This will:

- Load KPIs from kpi.json

- Register and initialize metrics

- Perform operations (add, set, increment, etc.)

- Push metrics to Prometheus PushGateway (http://localhost:9091)

## Sample KPI Configuration (kpi.json)
```bash
[
  {
    "name": "registration_success_rate_single_slice",
    "type": "gauge",
    "help": "Success rate of registrations per network slice",
    "labels": ["NetworkSlice"]
  },
  {
    "name": "registered_subscribers_udm",
    "type": "counter",
    "help": "Total registered subscribers in UDM",
    "labels": ["Network"]
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

## Push Output
```bash 
    http://localhost:9091/metrics 
```

## Author

Your Name – GitHub

## License

This project is licensed under the MIT License.