Step1: Create a logginglogmetric resource, operation succeeds

Step2: Update spec.metricDescriptor.metricKind and valueType fields for the resource, operation fails. We'll see no PUT request
in _http01 and resource has UpdateFailed message in _object01.
