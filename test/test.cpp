#include <iostream>
#include <string>
#include <vector>
#include "metrics_wrapper.hpp"

int main() {
    std::cout << "Initialize called" << std::endl;
    if (Initialize(const_cast<char*>("prometheus"), const_cast<char*>("testns")) != 0) {
        std::cerr << "Failed to initialize metrics framework" << std::endl;
        return 1;
    }
    std::cout << "MetricsFramework initialized successfully" << std::endl;

    std::cout << "Loading KPIs..." << std::endl;
    if (LoadKPIs(const_cast<char*>("models/kpi.json")) != 0) {
        std::cerr << "Failed to load KPIs" << std::endl;
        return 1;
    }
    std::cout << "KPIs loaded successfully" << std::endl;

    std::cout << "Registering metrics..." << std::endl;
    if (RegisterMetrics() != 0) {
        std::cerr << "RegisterMetrics failed" << std::endl;
        return 1;
    }
    std::cout << "Metrics registered successfully" << std::endl;

    // Prepare labels in a C++ friendly way
    std::vector<const char*> labels = {"Network", "test_network", "NetworkSlice", "test_slice", nullptr};

    std::cout << "Incrementing 'mean_registered_subscribers_amf'..." << std::endl;
    if (IncrementMetric(const_cast<char*>("mean_registered_subscribers_amf"), 
                       const_cast<char**>(labels.data()), 
                       static_cast<int>(labels.size() - 1)) != 0) {
        std::cerr << "IncrementMetric failed" << std::endl;
        return 1;
    }
    std::cout << "'mean_registered_subscribers_amf' incremented successfully" << std::endl;

    std::cout << "Trying to decrement (should fail)..." << std::endl;
    if (DecrementMetric(const_cast<char*>("mean_registered_subscribers_amf"), 
                       const_cast<char**>(labels.data()), 
                       static_cast<int>(labels.size() - 1)) != 0) {
        std::cout << "Expected: Decrement not allowed for 'mean_registered_subscribers_amf'" << std::endl;
    } else {
        std::cout << "Warning: Decrement succeeded, but shouldn't have" << std::endl;
    }

    std::cout << "Adding 5.5 to 'mean_registered_subscribers_amf'..." << std::endl;
    if (AddToMetric(const_cast<char*>("mean_registered_subscribers_amf"), 
                   5.5, 
                   const_cast<char**>(labels.data()), 
                   static_cast<int>(labels.size() - 1)) != 0) {
        std::cerr << "AddToMetric failed" << std::endl;
        return 1;
    }
    std::cout << "5.5 added to 'mean_registered_subscribers_amf'" << std::endl;

    // Gauge metric labels
    std::vector<const char*> gauge_labels = {"NetworkSlice", "test_slice", nullptr};
    if (SetMetric(const_cast<char*>("registration_success_rate_single_slice"), 
                 95.5, 
                 const_cast<char**>(gauge_labels.data()), 
                 static_cast<int>(gauge_labels.size() - 1)) != 0) {
        std::cerr << "SetMetric failed on gauge (should have worked)" << std::endl;
    } else {
        std::cout << "Gauge metric set successfully" << std::endl;
    }

    // Counter metric labels
    std::vector<const char*> counter_labels = {"Network", "test_network", nullptr};
    if (SetMetric(const_cast<char*>("mean_registered_subscribers_amf"), 
                 42.0, 
                 const_cast<char**>(counter_labels.data()), 
                 static_cast<int>(counter_labels.size() - 1)) == 0) {
        std::cout << "SetMetric succeeded on counter (should have failed)" << std::endl;
    } else {
        std::cout << "Properly rejected Set on counter metric" << std::endl;
    }

    std::cout << "Listing registered metrics:" << std::endl;
    char** metrics = ListMetrics();
    if (metrics == nullptr) {
        std::cerr << "Failed to list metrics" << std::endl;
        return 1;
    }

    int index = 0;
    while (metrics[index] != nullptr) {
        std::cout << "  - " << metrics[index] << std::endl;
        index++;
    }
    FreeStringArray(metrics, index);
    std::cout << "Metric names listed and freed" << std::endl;

    std::cout << "Pushing metrics to gateway..." << std::endl;
    if (PushMetrics(const_cast<char*>("http://localhost:9091"), 
                   const_cast<char*>("test_job")) != 0) {
        std::cerr << "PushMetrics failed (non-critical)" << std::endl;
    } else {
        std::cout << "Metrics pushed successfully" << std::endl;
    }

    std::cout << "All metric operations completed!" << std::endl;
    return 0;
}