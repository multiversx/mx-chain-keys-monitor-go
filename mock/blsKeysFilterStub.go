package mock

// BLSKeysFilterStub -
type BLSKeysFilterStub struct {
	ShouldNotifyCalled func(blsKey string) bool
}

// ShouldNotify -
func (stub *BLSKeysFilterStub) ShouldNotify(blsKey string) bool {
	if stub.ShouldNotifyCalled != nil {
		return stub.ShouldNotifyCalled(blsKey)
	}

	return true
}

// IsInterfaceNil -
func (stub *BLSKeysFilterStub) IsInterfaceNil() bool {
	return stub == nil
}
