//go:build linux && arm64

package opencv

import (
	_ "embed"
)

//go:embed dist/goopencv-linux-arm64.so
var embedLibData []byte

func embedLibName() string {
	return "goopencv-linux-arm64.so"
}
