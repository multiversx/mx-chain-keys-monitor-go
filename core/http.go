package core

// IsStatusCodeIs2xx returns true if the provided status code is not 2xx
func IsStatusCodeIs2xx(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
