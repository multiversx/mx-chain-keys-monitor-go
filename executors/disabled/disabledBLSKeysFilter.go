package disabled

type disabledBLSKeysFilter struct{}

// NewDisabledBLSKeysFilter will create a new instance of type disabledBLSKeysFilter
func NewDisabledBLSKeysFilter() *disabledBLSKeysFilter {
	return &disabledBLSKeysFilter{}
}

// ShouldNotify returns true regardless of the key provided
func (disabled *disabledBLSKeysFilter) ShouldNotify(_ string) bool {
	return true
}

// IsInterfaceNil returns true if there is no value under the interface
func (disabled *disabledBLSKeysFilter) IsInterfaceNil() bool {
	return disabled == nil
}
