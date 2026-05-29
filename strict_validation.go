package opencv

import (
	"sync/atomic"
)

// strictColorValidation controls optional runtime model/code validation.
// Default: off (parity-safe with OpenCV behavior).
var strictColorValidation atomic.Bool

// SetStrictColorValidation enables/disables strict color model validation.
// Off by default to preserve OpenCV parity behavior.
func SetStrictColorValidation(enabled bool) {
	strictColorValidation.Store(enabled)
}

// StrictColorValidation reports whether strict validation is enabled.
func StrictColorValidation() bool {
	return strictColorValidation.Load()
}
