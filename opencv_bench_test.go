package opencv

import (
	"context"
	"os"
	"testing"
)

func BenchmarkCvtColor_BGR2Gray_1080p(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(1080, 1920, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(1080, 1920, CV8UC1)
	if err != nil {
		b.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.CvtColor(src, dst, ColorBGR2Gray); err != nil {
			b.Fatalf("CvtColor: %v", err)
		}
	}
}

func BenchmarkPipeline_NoConvert_1080p(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(1080, 1920, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(540, 960, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.Resize(src, dst, Size{Width: 960, Height: 540}); err != nil {
			b.Fatalf("Resize: %v", err)
		}
	}
}

func BenchmarkCvtColor_BGR2RGBA_1080p(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(1080, 1920, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(1080, 1920, CV8UC4)
	if err != nil {
		b.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.CvtColor(src, dst, ColorBGR2BGRA); err != nil {
			b.Fatalf("CvtColor: %v", err)
		}
	}
}

func BenchmarkCvtColor_BGR2RGBA_4K(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(2160, 3840, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(2160, 3840, CV8UC4)
	if err != nil {
		b.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.CvtColor(src, dst, ColorBGR2BGRA); err != nil {
			b.Fatalf("CvtColor: %v", err)
		}
	}
}

func BenchmarkIMWrite_PNG_1080p_ModelOverrideRGBA(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(1080, 1920, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	dir := b.TempDir()
	out := dir + "\\bench_out.png"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.IMWrite(out, m, RGBA); err != nil {
			b.Fatalf("IMWrite: %v", err)
		}
		_ = os.Remove(out)
	}
}

func BenchmarkIMWrite_JPEG_1080p_ModelOverrideRGBA(b *testing.B) {
	path := dllPath()
	if path == "" {
		b.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(1080, 1920, CV8UC3)
	if err != nil {
		b.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	dir := b.TempDir()
	out := dir + "\\bench_out.jpg"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.IMWrite(out, m, RGBA); err != nil {
			b.Fatalf("IMWrite JPEG: %v", err)
		}
		_ = os.Remove(out)
	}
}
