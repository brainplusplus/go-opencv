package opencv

import (
	"context"
	"errors"
	"testing"
)

func TestNewAutoLoad(t *testing.T) {
	if len(embedLibData) == 0 {
		t.Skip("no embedded library for this platform")
	}

	r, err := New(context.Background())
	if err != nil {
		t.Fatalf("New() with auto-load: %v", err)
	}
	defer r.Close()
}

func TestMatDeleteAlias(t *testing.T) {
	var m Mat
	if !errors.Is(m.Delete(), ErrInvalidMat) {
		t.Fatalf("zero Mat Delete() must report ErrInvalidMat")
	}
}
