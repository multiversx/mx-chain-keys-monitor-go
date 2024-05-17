package disabled

type disabledCloser struct {
}

// NewDisabledCloser creates a new disabled closer instance
func NewDisabledCloser() *disabledCloser {
	return &disabledCloser{}
}

// Close returns nil
func (disabled *disabledCloser) Close() error {
	return nil
}
