package monitor

type pollingHandler interface {
	StartProcessingLoop() error
	IsRunning() bool
	Close() error
	IsInterfaceNil() bool
}
