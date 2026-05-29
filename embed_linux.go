//go:build linux && amd64

package opencv

import (
	_ "embed"
)

//go:embed dist/goopencv.so
var embedLibData []byte

func embedLibName() string {
	return "goopencv.so"
}
