//go:build !(windows && amd64) && !(linux && amd64) && !(linux && arm64) && !darwin

package opencv

import (
	"fmt"
	"runtime"
)

var embedLibData []byte = nil

func embedLibName() string {
	return ""
}

func init() {
	// This package is selected on unsupported platforms — provide clear error
	_ = fmt.Sprintf("go-opencv: unsupported platform %s/%s. Supported: windows/amd64, linux/amd64, linux/arm64, darwin/amd64, darwin/arm64", runtime.GOOS, runtime.GOARCH)
}
