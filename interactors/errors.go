package interactors

import "errors"

var (
	errNilHTTPClientWrapper = errors.New("nil HTTP client wrapper")
	errReturnCodeIsNotOk    = errors.New("HTTP return code is not OK")
	errInvalidResponse      = errors.New("invalid response, nil statistics map")
	errContextClosing       = errors.New("context closing")
)
