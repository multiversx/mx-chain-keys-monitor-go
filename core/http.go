package core

// IsHttpStatusCodeSuccess returns true if the provided status code is 2xx
func IsHttpStatusCodeSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
