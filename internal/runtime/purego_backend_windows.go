//go:build windows

package runtime

import (
	"syscall"
)

func openLibrary(path string) (uintptr, error) {
	h, err := syscall.LoadLibrary(path)
	return uintptr(h), err
}

func closeLibrary(lib uintptr) {
	// Skip FreeLibrary on Windows — OpenCV creates global state
	// (thread pools, TBB, etc.) that crashes on DLL unload.
	// Let the OS clean up when the process exits.
	// _ = syscall.FreeLibrary(syscall.Handle(lib))
}
