#pragma once

extern "C" {
    int Initialize(char* backend, char* namespaceName);
    int LoadKPIs(char* path);
    int RegisterMetrics();
    int IncrementMetric(char* metricName, char** labels, int labelCount);
    int DecrementMetric(char* metricName, char** labels, int labelCount);
    int AddToMetric(char* metricName, double value, char** labels, int labelCount);
    int SetMetric(char* metricName, double value, char** labels, int labelCount);
    char** ListMetrics();
    void FreeStringArray(char** arr, int length);
    int PushMetrics(char* gatewayURL, char* jobName);
}
