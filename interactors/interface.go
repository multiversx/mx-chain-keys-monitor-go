package interactors

import "context"

// HTTPClientWrapper defines what an HTTP client wrapper should implement
type HTTPClientWrapper interface {
	GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
	IsInterfaceNil() bool
}
