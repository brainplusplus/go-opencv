//go:build darwin || linux || freebsd || netbsd

package runtime

import "github.com/ebitengine/purego"

func openLibrary(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

func closeLibrary(lib uintptr) {
	// Skip Dlclose — OpenCV creates global state (TBB thread pools, IPP, etc.)
	// that crashes on shared library unload. Let the OS clean up on process exit.
	// _ = purego.Dlclose(lib)
}
