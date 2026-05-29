//go:build windows && amd64

package opencv

import (
	_ "embed"
)

//go:embed dist/goopencv.dll
var embedLibData []byte

func embedLibName() string {
	return "goopencv.dll"
}
