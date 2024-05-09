package mock

import "context"

// BLSKEysFetcherStub -
type BLSKEysFetcherStub struct {
	GetAllBLSKeysHandler func(ctx context.Context, sender string) ([]string, error)
}

// GetAllBLSKeys -
func (stub *BLSKEysFetcherStub) GetAllBLSKeys(ctx context.Context, sender string) ([]string, error) {
	if stub.GetAllBLSKeysHandler != nil {
		return stub.GetAllBLSKeysHandler(ctx, sender)
	}

	return make([]string, 0), nil
}

// IsInterfaceNil -
func (stub *BLSKEysFetcherStub) IsInterfaceNil() bool {
	return stub == nil
}
