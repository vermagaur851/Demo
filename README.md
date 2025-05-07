## Project Overview
- The Amantya Metrics Framework is a versatile, backend-agnostic metric management system designed to provide seamless integration with popular monitoring tools such as Prometheus and Datadog. It supports the registration, management, and interaction with various types of metrics, including Counters and Gauges. This framework is ideal for systems and applications that require real-time monitoring, logging, and metric collection.

## Features
- **Backend Support :**: Choose from a variety of backend options like Prometheus and Datadog for metric storage and collection.

- **KPI Integration :**: Easily load and register Key Performance Indicators (KPIs) from external JSON files, allowing dynamic metric handling.

- **Metric Operations :**: Perform operations such as incrementing, decrementing, and adding to metrics with intuitive API methods.

- **Automatic Metric Initialization :** Ensure metrics are initialized with default values, ready for use in production systems.

- **Push Metrics to Gateway :** Push metrics to external monitoring systems (e.g., Prometheus PushGateway) for real-time visualization and alerting.

- **Metric Registration & List Management :** Register and list custom metrics dynamically for easy tracking and management.

## Getting Started
- Follow these instructions to get the project up and running on your local machine.

### Prerequisites

Before you begin, make sure you have the following software installed:

- [Go](https://golang.org/dl/) (for Go-based services)
- [Prometheus](https://prometheus.io/docs/introduction/overview/) (for monitoring)
  
You should also have a basic understanding of how to work with Go and Prometheus.

### Installation

Clone the repository:
To install and run the project locally, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/your-repository.git

2. **Navigate to the project folder:**
    ```bash
    cd your-repository

3. **Install dependencies:**
    ```bash
    go mod tidy

4. **Configure Prometheus:**
    If you are using Prometheus for metrics, make sure it's running and properly configured to collect data from your project. You can follow the Prometheus setup guide if it's not set up yet.

### Run
    ```
    go run main.go
    ```

## Project Structure

\`\`\`
AMANTYA_METRICS/
├── build/
├── datadogbackend/
├── lang_wrapper/
├── metricsInterface/
├── metricsregistry/
├── models/
├── prometheusbackend/
├── service/
├── uitls/
├── Makefile
├── README.md
\`\`\`

## License

[MIT](LICENSE)

## Author

Your Name - [GitHub](https://github.com/yourusername)
