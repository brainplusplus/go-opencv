package opencv

import "errors"

var (
	// ErrBackendUnavailable means no native DLL was configured.
	ErrBackendUnavailable = errors.New("opencv: backend unavailable")

	// ErrClosed means the runtime or object was already closed.
	ErrClosed = errors.New("opencv: closed")

	// ErrInvalidMat means a Mat handle is nil, closed, or owned by another runtime.
	ErrInvalidMat = errors.New("opencv: invalid mat")

	// ErrNotImplemented marks API surface reserved for OpenCV.js parity but not wired yet.
	ErrNotImplemented = errors.New("opencv: not implemented")
)
