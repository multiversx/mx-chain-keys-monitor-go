package notifiers

import "errors"

var (
	errNilLogger         = errors.New("nil logger")
	errReturnCodeIsNotOk = errors.New("HTTP return code is not OK")
)
