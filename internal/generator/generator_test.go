package generator

import (
	"strings"
	"testing"

	"github.com/brainplusplus/go-opencv/internal/contract"
)

func TestRenderABIMarkdownIncludesStableSymbols(t *testing.T) {
	out := RenderABIMarkdown(contract.Modules)
	for _, want := range []string{"goopencv_mat_new", "goopencv_mat_rows", "goopencv_imgproc_cvt_color"} {
		if !strings.Contains(out, want) {
			t.Fatalf("ABI markdown missing %q", want)
		}
	}
}

func TestRenderCPPSkeletonIncludesKnownStubs(t *testing.T) {
	out := RenderCPPSkeleton(contract.Modules)
	for _, want := range []string{"extern \"C\"", "uint64_t goopencv_mat_new", "int32_t goopencv_imgproc_resize"} {
		if !strings.Contains(out, want) {
			t.Fatalf("C++ skeleton missing %q", want)
		}
	}
}
