package mock

import "context"

// ExecutorStub -
type ExecutorStub struct {
	ExecuteHandler func(ctx context.Context) error
}

// Execute -
func (stub *ExecutorStub) Execute(ctx context.Context) error {
	if stub.ExecuteHandler != nil {
		return stub.ExecuteHandler(ctx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *ExecutorStub) IsInterfaceNil() bool {
	return stub == nil
}
