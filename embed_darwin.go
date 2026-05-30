//go:build darwin && (amd64 || arm64)

package opencv

import (
	_ "embed"
)

//go:embed dist/goopencv.dylib
var embedLibData []byte

func embedLibName() string {
	return "goopencv.dylib"
}
