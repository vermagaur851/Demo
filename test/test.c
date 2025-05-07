#include <stdio.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "libamantyametrics.h"
#ifdef __cplusplus
}
#endif


int main() {
    printf("Initialize called\n");
    if (Initialize("prometheus", "testns") != 0) {
        printf("Failed to initialize metrics framework\n");
        return 1;
    }
    printf("MetricsFramework initialized successfully\n");

    printf("Loading KPIs...\n");
    if (LoadKPIs("models/kpi.json") != 0) {
        printf("Failed to load KPIs\n");
        return 1;
    }
    printf("KPIs loaded successfully\n");

    printf("Registering metrics...\n");
    if (RegisterMetrics() != 0) {
        printf("RegisterMetrics failed\n");
        return 1;
    }
    printf("Metrics registered successfully\n");

    // Labels for mean_registered_subscribers_amf
    const char* labels[] = {
        "Network", "test_network",
        "NetworkSlice", "test_slice",
        NULL
    };

    printf("Incrementing 'mean_registered_subscribers_amf'...\n");
    if (IncrementMetric("mean_registered_subscribers_amf",(char **) labels, 4) != 0) {
        printf("IncrementMetric failed\n");
        return 1;
    }
    printf("'mean_registered_subscribers_amf' incremented successfully\n");

    // This will be skipped, as decrement is not allowed
    printf("Trying to decrement (should fail)...\n");
    if (DecrementMetric("mean_registered_subscribers_amf", (char **)labels, 4) != 0) {
        printf("Expected: Decrement not allowed for 'mean_registered_subscribers_amf'\n");
    } else {
        printf("Warning: Decrement succeeded, but shouldn't have\n");
    }

    printf("Adding 5.5 to 'mean_registered_subscribers_amf'...\n");
    if (AddToMetric("mean_registered_subscribers_amf", 5.5,(char **) labels, 4) != 0) {
        printf("AddToMetric failed\n");
        return 1;
    }
    printf("5.5 added to 'mean_registered_subscribers_amf'\n");

     // Test setting a gauge metric (should work)
    const char* gauge_labels[] = {"NetworkSlice", "test_slice", NULL};
    if (SetMetric("registration_success_rate_single_slice", 95.5,(char **) gauge_labels, 2) != 0) {
        printf("SetMetric failed on gauge (should have worked)\n");
    } else {
        printf("Gauge metric set successfully\n");
    }

    // Test setting a counter metric (should fail)
    const char* counter_labels[] = {"Network", "test_network", NULL};
    if (SetMetric("mean_registered_subscribers_amf", 42.0,(char **) counter_labels, 2) == 0) {
        printf("SetMetric succeeded on counter (should have failed)\n");
    } else {
        printf("Properly rejected Set on counter metric\n");
    }

    printf("Listing registered metrics:\n");
    char** metrics = ListMetrics();
    if (metrics == NULL) {
        printf("Failed to list metrics\n");
        return 1;
    }

    int index = 0;
    while (metrics[index] != NULL) {
        printf("  - %s\n", metrics[index]);
        index++;
    }

    FreeStringArray(metrics, index);
    printf("Metric names listed and freed\n");

    printf("Pushing metrics to gateway...\n");
    if (PushMetrics("http://localhost:9091", "test_job") != 0) {
        printf("PushMetrics failed (non-critical)\n");
    } else {
        printf("Metrics pushed successfully\n");
    }

    printf("All metric operations completed!\n");
    return 0;
}
