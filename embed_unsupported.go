//go:build !(windows && amd64) && !(linux && amd64) && !(linux && arm64) && !(darwin && amd64) && !(darwin && arm64)

package opencv

var embedLibData []byte = nil

func embedLibName() string {
	return ""
}
