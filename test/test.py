from ctypes import CDLL, c_char_p, c_double, c_int, POINTER
import sys

# Load shared library
lib = CDLL("./build/libamantyametrics.so")

# Function signatures
MetricWithLabels = [c_char_p, POINTER(c_char_p), c_int]
MetricWithValue = [c_char_p, c_double, POINTER(c_char_p), c_int]

lib.Initialize.argtypes         = [c_char_p, c_char_p]
lib.LoadKPIs.argtypes           = [c_char_p]
lib.RegisterMetrics.argtypes    = []

lib.IncrementMetric.argtypes    = MetricWithLabels
lib.DecrementMetric.argtypes    = MetricWithLabels
lib.AddToMetric.argtypes        = MetricWithValue
lib.SetMetric.argtypes          = MetricWithValue

lib.PushMetrics.argtypes        = [c_char_p, c_char_p]
lib.ListMetrics.restype         = POINTER(c_char_p)
lib.FreeStringArray.argtypes    = [POINTER(c_char_p), c_int]

# Helper function to create C-compatible label array
def make_labels(labels_dict):
    flat = []
    for k, v in labels_dict.items():
        flat.extend([k.encode(), v.encode()])
    arr = (c_char_p * (len(flat) + 1))()
    arr[:len(flat)] = flat
    arr[len(flat)] = None
    return arr, len(flat)

# Initialization
print("Initialize called")
if lib.Initialize(b"prometheus", b"testns") != 0:
    print("Failed to initialize metrics framework")
    sys.exit(1)
print("MetricsFramework initialized successfully")

# Load KPIs
print("Loading KPIs...")
if lib.LoadKPIs(b"models/kpi.json") != 0:
    print("Failed to load KPIs")
    sys.exit(1)
print("KPIs loaded successfully")

# Register
print("Registering metrics...")
if lib.RegisterMetrics() != 0:
    print("Failed to register metrics")
    sys.exit(1)
print("Metrics registered successfully")

# Labels
labels, label_count = make_labels({"Network": "test_network", "NetworkSlice": "test_slice"})

# Increment metric
print("Incrementing 'mean_registered_subscribers_amf'...")
if lib.IncrementMetric(b"mean_registered_subscribers_amf", labels, label_count) != 0:
    print("Increment failed")
else:
    print("Incremented successfully")

# Decrement (should fail)
print("Trying to decrement (should fail)...")
if lib.DecrementMetric(b"mean_registered_subscribers_amf", labels, label_count) != 0:
    print("Expected: Decrement not allowed")
else:
    print("Warning: Decrement succeeded")

# Add to counter
print("Adding 5.5 to counter...")
if lib.AddToMetric(b"mean_registered_subscribers_amf", 5.5, labels, label_count) != 0:
    print("AddToMetric failed")
else:
    print("5.5 added successfully")

# Set gauge
gauge_labels, gauge_count = make_labels({"NetworkSlice": "test_slice"})
print("Setting gauge value...")
if lib.SetMetric(b"registration_success_rate_single_slice", 95.5, gauge_labels, gauge_count) != 0:
    print("Failed to set gauge")
else:
    print("Gauge set successfully")

# Set on counter (should fail)
counter_labels, counter_count = make_labels({"Network": "test_network"})
print("Trying to set counter (should fail)...")
if lib.SetMetric(b"mean_registered_subscribers_amf", 42.0, counter_labels, counter_count) == 0:
    print("Warning: Set succeeded on counter")
else:
    print("Properly rejected Set on counter")

# List metrics
print("Listing registered metrics:")
metrics = lib.ListMetrics()
if not metrics:
    print("Failed to list metrics")
else:
    idx = 0
    while metrics[idx]:
        print(" -", metrics[idx].decode())
        idx += 1
    lib.FreeStringArray(metrics, idx)
    print("Freed listed metric names")

# Push to PushGateway
print("Pushing metrics...")
if lib.PushMetrics(b"http://localhost:9091", b"test_job") != 0:
    print("Push failed (non-critical)")
else:
    print("Metrics pushed successfully")

print("All metric operations completed!")
