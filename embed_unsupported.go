//go:build !(windows && amd64) && !(linux && amd64) && !darwin

package opencv

var embedLibData []byte = nil

func embedLibName() string {
	return ""
}
