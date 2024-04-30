package mock

import (
	"context"
	"errors"
)

// HTTPClientWrapperStub -
type HTTPClientWrapperStub struct {
	GetHTTPHandler  func(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTPHandler func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
}

// GetHTTP -
func (stub *HTTPClientWrapperStub) GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error) {
	if stub.GetHTTPHandler != nil {
		return stub.GetHTTPHandler(ctx, endpoint)
	}

	return nil, 0, errors.New("not implemented")
}

// PostHTTP -
func (stub *HTTPClientWrapperStub) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
	if stub.PostHTTPHandler != nil {
		return stub.PostHTTPHandler(ctx, endpoint, data)
	}

	return nil, 0, errors.New("not implemented")
}

// IsInterfaceNil -
func (stub *HTTPClientWrapperStub) IsInterfaceNil() bool {
	return stub == nil
}
