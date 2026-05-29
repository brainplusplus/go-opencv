package opencv

import "fmt"

type Option func(*config) error

type config struct {
	dll string
}

// WithDLL configures a native shared library (goopencv.dll / .so / .dylib).
func WithDLL(path string) Option {
	return func(c *config) error {
		if path == "" {
			return fmt.Errorf("opencv: empty dll path: %w", ErrBackendUnavailable)
		}
		c.dll = path
		return nil
	}
}
