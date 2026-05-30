package opencv

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"testing"
)

// dllPath returns the native backend path for the current platform if it exists.
func dllPath() string {
	var candidates []string
	switch runtime.GOOS {
	case "windows":
		candidates = []string{"dist/goopencv.dll"}
	case "darwin":
		candidates = []string{"dist/goopencv.dylib"}
	case "linux":
		if runtime.GOARCH == "arm64" {
			candidates = []string{"dist/goopencv-linux-arm64.so", "dist/goopencv.so"}
		} else {
			candidates = []string{"dist/goopencv.so", "dist/goopencv-linux-arm64.so"}
		}
	default:
		return ""
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func TestNewWithDLL(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found — build with build-tools/build-goopencv.bat first")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New(WithDLL(%q)): %v", path, err)
	}
	defer r.Close()
}

func TestDLLNewMat(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(100, 200, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	if m == nil {
		t.Fatal("NewMat returned nil")
	}
}

func TestDLLMatRowsColsType(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(100, 200, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	rows, err := m.Rows()
	if err != nil {
		t.Fatalf("Rows: %v", err)
	}
	if rows != 100 {
		t.Errorf("Rows = %d, want 100", rows)
	}

	cols, err := m.Cols()
	if err != nil {
		t.Fatalf("Cols: %v", err)
	}
	if cols != 200 {
		t.Errorf("Cols = %d, want 200", cols)
	}

	typ, err := m.Type()
	if err != nil {
		t.Fatalf("Type: %v", err)
	}
	// CV_8UC3 = 16
	if int(typ) != 16 {
		t.Errorf("Type = %d, want 16 (CV_8UC3)", typ)
	}
}

func TestDLLMatClone(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(100, 200, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	clone, err := m.Clone()
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	if clone == nil {
		t.Fatal("Clone returned nil")
	}
	defer clone.Close()

	// Verify clone has same dimensions
	rows, _ := clone.Rows()
	cols, _ := clone.Cols()
	if rows != 100 || cols != 200 {
		t.Errorf("Clone dims = %dx%d, want 100x200", rows, cols)
	}
}

func TestDLLCvtColor(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(50, 50, CV8UC4)
	if err != nil {
		t.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	if err := r.CvtColor(src, dst, ColorBGR2BGRA); err != nil {
		t.Fatalf("CvtColor: %v", err)
	}

	// After BGR->BGRA conversion, type should still be CV_8UC4 on dst
	typ, err := dst.Type()
	if err != nil {
		t.Fatalf("dst.Type(): %v", err)
	}
	if int(typ) != 24 {
		t.Errorf("dst type = %d, want 24 (CV_8UC4)", typ)
	}
}

func TestDLLResize(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.NewMat(50, 50, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	if err := r.Resize(src, dst, Size{Width: 50, Height: 50}); err != nil {
		t.Fatalf("Resize: %v", err)
	}

	// After resize, dst dimensions should be 50x50
	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 {
		t.Errorf("dst.Rows = %d, want 50", rows)
	}
	if cols != 50 {
		t.Errorf("dst.Cols = %d, want 50", cols)
	}
}

// TestDLLMatDoubleClose verifies that Mat.Close() can be called safely.
func TestDLLMatDoubleClose(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}

	if err := m.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := m.Delete(); err == nil {
		t.Errorf("second Delete = %v, want error", err)
	}
}

func TestDLLReadImage(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Create a temp PNG image (4x4 red pixels)
	tmpDir := t.TempDir()
	imgPath := tmpDir + "\\test.png"
	createTestPNG(t, imgPath, 4, 4)

	m, err := r.IMRead(imgPath)
	if err != nil {
		t.Fatalf("IMRead: %v", err)
	}
	defer m.Close()

	// Verify dimensions
	rows, _ := m.Rows()
	cols, _ := m.Cols()
	if rows != 4 || cols != 4 {
		t.Errorf("IMRead dims = %dx%d, want 4x4", rows, cols)
	}

	// Verify RGBA type (CV_8UC4 = 24)
	typ, _ := m.Type()
	if int(typ) != 16 { // CV_8UC3 = 0 + 16 = 16
		t.Errorf("IMRead type = %d, want 16 (CV_8UC3)", typ)
	}
}

func TestDLLBlur(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.Blur(src, Size{Width: 5, Height: 5})
	if err != nil {
		t.Fatalf("Blur: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("Blur dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLGaussianBlur(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.GaussianBlur(src, Size{Width: 5, Height: 5}, 0)
	if err != nil {
		t.Fatalf("GaussianBlur: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("GaussianBlur dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLThreshold(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.Threshold(src, 128, 255, ThresholdBinary)
	if err != nil {
		t.Fatalf("Threshold: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("Threshold dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLDrawing(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	img, err := r.NewMat(200, 200, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	// Test with Go color.RGBA style
	if err := r.Rectangle(img, Rect{X: 10, Y: 10, Width: 50, Height: 50}, color.RGBA{255, 0, 0, 255}, 2); err != nil {
		t.Fatalf("Rectangle with RGBA: %v", err)
	}

	// Test with OpenCV BGR style
	if err := r.Circle(img, Point{X: 100, Y: 100}, 30, BGRColor{B: 0, G: 255, R: 0, A: 255}, 2); err != nil {
		t.Fatalf("Circle with BGR: %v", err)
	}

	// Test Line
	if err := r.Line(img, Point{X: 0, Y: 0}, Point{X: 199, Y: 199}, color.RGBA{0, 0, 255, 255}, 1); err != nil {
		t.Fatalf("Line: %v", err)
	}

	// Verify img not empty
	empty, err := img.Empty()
	if err != nil {
		t.Fatalf("Empty: %v", err)
	}
	if empty {
		t.Error("img should not be empty after drawing")
	}
}

func TestDLLMatEmptyElemSize(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Non-empty mat
	m, err := r.NewMat(10, 10, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	empty, err := m.Empty()
	if err != nil {
		t.Fatalf("Empty: %v", err)
	}
	if empty {
		t.Error("NewMat(10,10,CV8UC3) should not be empty")
	}

	es, err := m.ElemSize()
	if err != nil {
		t.Fatalf("ElemSize: %v", err)
	}
	if es != 3 { // CV_8UC3 = 3 bytes per element
		t.Errorf("ElemSize = %d, want 3", es)
	}
}

func TestDLLMatCopyToChannelsStep(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(10, 20, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	// Channels
	ch, err := m.Channels()
	if err != nil {
		t.Fatalf("Channels: %v", err)
	}
	if ch != 3 {
		t.Errorf("Channels = %d, want 3", ch)
	}

	// Step
	step, err := m.Step()
	if err != nil {
		t.Fatalf("Step: %v", err)
	}
	if step != 20*3 { // cols * channels
		t.Errorf("Step = %d, want %d", step, 20*3)
	}

	// CopyTo
	buf := make([]byte, 10*step)
	n, err := m.CopyTo(buf)
	if err != nil {
		t.Fatalf("CopyTo: %v", err)
	}
	if n != 10*step {
		t.Errorf("CopyTo returned %d, want %d", n, 10*step)
	}
}

func TestDLLSaveImage(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Create a small RGBA Mat
	mat, err := r.NewMat(10, 10, CV8UC4)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer mat.Close()

	// Save as PNG
	tmpDir := t.TempDir()
	outPath := tmpDir + "\\out.png"
	if err := r.IMWrite(outPath, mat); err != nil {
		t.Fatalf("IMWrite PNG: %v", err)
	}

	// Verify file exists and is readable
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("output dims = %dx%d, want 10x10", bounds.Dx(), bounds.Dy())
	}

	// Save as JPEG
	outJPG := tmpDir + "\\out.jpg"
	if err := r.IMWrite(outJPG, mat); err != nil {
		t.Fatalf("IMWrite JPEG: %v", err)
	}

	// Verify JPEG
	f2, err := os.Open(outJPG)
	if err != nil {
		t.Fatalf("open JPEG output: %v", err)
	}
	defer f2.Close()

	img2, _, err := image.Decode(f2)
	if err != nil {
		t.Fatalf("decode JPEG output: %v", err)
	}
	if img2.Bounds().Dx() != 10 || img2.Bounds().Dy() != 10 {
		t.Errorf("JPEG output dims = %dx%d, want 10x10", img2.Bounds().Dx(), img2.Bounds().Dy())
	}
}

func TestDLLIMReadModels(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Create temp PNG
	tmpDir := t.TempDir()
	imgPath := tmpDir + "\\test.png"
	createTestPNG(t, imgPath, 8, 6)

	// Default: BGR
	bgr, err := r.IMRead(imgPath)
	if err != nil {
		t.Fatalf("IMRead default: %v", err)
	}
	defer bgr.Close()
	typ, _ := bgr.Type()
	if typ != CV8UC3 {
		t.Errorf("IMRead default type = %d, want CV8UC3", typ)
	}

	// RGB
	rgb, err := r.IMRead(imgPath, RGB)
	if err != nil {
		t.Fatalf("IMRead RGB: %v", err)
	}
	defer rgb.Close()
	typ, _ = rgb.Type()
	if typ != CV8UC3 {
		t.Errorf("IMRead RGB type = %d, want CV8UC3", typ)
	}

	// RGBA
	rgba, err := r.IMRead(imgPath, RGBA)
	if err != nil {
		t.Fatalf("IMRead RGBA: %v", err)
	}
	defer rgba.Close()
	typ, _ = rgba.Type()
	if typ != CV8UC4 {
		t.Errorf("IMRead RGBA type = %d, want CV8UC4", typ)
	}

	// Gray
	gray, err := r.IMRead(imgPath, Gray)
	if err != nil {
		t.Fatalf("IMRead Gray: %v", err)
	}
	defer gray.Close()
	typ, _ = gray.Type()
	if typ != CV8UC1 {
		t.Errorf("IMRead Gray type = %d, want CV8UC1", typ)
	}
}

func TestDLLIMReadBytesModels(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Prepare encoded PNG bytes
	tmpDir := t.TempDir()
	imgPath := tmpDir + "\\bytes.png"
	createTestPNG(t, imgPath, 9, 7)
	data, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	cases := []struct {
		name  string
		model ColorModel
		want  MatType
	}{
		{"BGR", BGR, CV8UC3},
		{"RGB", RGB, CV8UC3},
		{"RGBA", RGBA, CV8UC4},
		{"Gray", Gray, CV8UC1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.IMReadBytes(data, tc.model)
			if err != nil {
				t.Fatalf("IMReadBytes(%v): %v", tc.model, err)
			}
			defer m.Close()

			typ, err := m.Type()
			if err != nil {
				t.Fatalf("Type: %v", err)
			}
			if typ != tc.want {
				t.Fatalf("type = %d, want %d", typ, tc.want)
			}

			cm, err := m.ColorModel()
			if err != nil {
				t.Fatalf("ColorModel: %v", err)
			}
			if cm != tc.model {
				t.Fatalf("model = %v, want %v", cm, tc.model)
			}
		})
	}
}

func TestDLLMatColorModelPropagation(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	img, err := r.NewMat(50, 60, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	m, err := img.ColorModel()
	if err != nil {
		t.Fatalf("ColorModel: %v", err)
	}
	if m != BGR {
		t.Fatalf("new CV8UC3 model = %v, want BGR", m)
	}

	// Preserve model ops
	blur, err := r.Blur(img, Size{Width: 3, Height: 3})
	if err != nil {
		t.Fatalf("Blur: %v", err)
	}
	defer blur.Close()
	bm, _ := blur.ColorModel()
	if bm != BGR {
		t.Errorf("Blur model = %v, want BGR", bm)
	}

	resized, err := r.NewMat(25, 30, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat resized: %v", err)
	}
	defer resized.Close()
	if err := r.Resize(img, resized, Size{Width: 30, Height: 25}); err != nil {
		t.Fatalf("Resize: %v", err)
	}
	rm, _ := resized.ColorModel()
	if rm != BGR {
		t.Errorf("Resize model = %v, want BGR", rm)
	}

	// Deterministic Gray output
	gray, err := r.NewMat(50, 60, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat gray: %v", err)
	}
	defer gray.Close()
	if err := r.CvtColor(img, gray, ColorBGR2Gray); err != nil {
		t.Fatalf("CvtColor: %v", err)
	}
	gm, _ := gray.ColorModel()
	if gm != Gray {
		t.Errorf("Gray model = %v, want Gray", gm)
	}
}

func TestStrictColorValidation(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	prev := StrictColorValidation()
	defer SetStrictColorValidation(prev)
	SetStrictColorValidation(true)

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	graySrc, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat gray: %v", err)
	}
	defer graySrc.Close()

	dst, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	// Invalid in strict mode: Gray source with *2Gray conversion group
	err = r.CvtColor(graySrc, dst, ColorBGR2Gray)
	if err == nil {
		t.Fatalf("expected strict validation error, got nil")
	}

	// Valid case should pass
	bgrSrc, err := r.NewMat(10, 10, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat bgr: %v", err)
	}
	defer bgrSrc.Close()
	err = r.CvtColor(bgrSrc, dst, ColorBGR2Gray)
	if err != nil {
		t.Fatalf("expected valid conversion, got err: %v", err)
	}
}

func TestStrictColorValidationMatrix(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	prev := StrictColorValidation()
	defer SetStrictColorValidation(prev)
	SetStrictColorValidation(true)

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	tests := []struct {
		name    string
		typ     MatType
		code    ColorConversionCode
		wantErr bool
	}{
		// Gray source
		{"Gray_to_BGR_valid", CV8UC1, ColorGray2BGR, false},
		{"Gray_to_BGRA_valid", CV8UC1, ColorGray2BGRA, false},
		{"Gray_to_Gray_invalid", CV8UC1, ColorBGR2Gray, true},

		// BGR source
		{"BGR_to_Gray_valid", CV8UC3, ColorBGR2Gray, false},
		{"BGR_to_BGRA_valid", CV8UC3, ColorBGR2BGRA, false},
		{"BGR_to_BGRA_invalid_group", CV8UC3, ColorBGRA2BGR, true},
		{"BGR_to_Gray2_invalid", CV8UC3, ColorGray2BGR, true},

		// RGBA source
		{"RGBA_to_BGR_valid", CV8UC4, ColorRGBA2BGR, false},
		{"RGBA_to_BGRA_invalid_group", CV8UC4, ColorBGR2BGRA, true},
		{"RGBA_to_Gray_valid", CV8UC4, ColorRGBA2Gray, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src, err := r.NewMat(10, 10, tc.typ)
			if err != nil {
				t.Fatalf("NewMat src: %v", err)
			}
			defer src.Close()

			dst, err := r.NewMat(10, 10, CV8UC4)
			if err != nil {
				t.Fatalf("NewMat dst: %v", err)
			}
			defer dst.Close()

			err = r.CvtColor(src, dst, tc.code)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
		})
	}

	// strict off parity path: valid conversion must pass.
	SetStrictColorValidation(false)
	src, err := r.NewMat(10, 10, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src parity: %v", err)
	}
	defer src.Close()
	dst, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat dst parity: %v", err)
	}
	defer dst.Close()
	if err := r.CvtColor(src, dst, ColorBGR2Gray); err != nil {
		t.Fatalf("strict off valid conversion should pass, got: %v", err)
	}
}

func TestIMWriteUnknownRequiresOverride(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	m, err := r.NewMat(8, 8, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()
	m.model = Unknown

	out := t.TempDir() + "\\unknown.png"
	if err := r.IMWrite(out, m); err == nil {
		t.Fatalf("expected IMWrite error for Unknown model without override")
	}

	if err := r.IMWrite(out, m, BGR); err != nil {
		t.Fatalf("expected IMWrite success with explicit override, got: %v", err)
	}
}

func TestIMWriteModelOverride(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Create BGR mat, then force RGBA write model.
	m, err := r.NewMat(16, 12, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	tmpDir := t.TempDir()
	outPath := tmpDir + "\\override.png"
	if err := r.IMWrite(outPath, m, RGBA); err != nil {
		t.Fatalf("IMWrite override RGBA: %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open written file: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode written file: %v", err)
	}
	if img.Bounds().Dx() != 12 || img.Bounds().Dy() != 16 {
		t.Errorf("override output dims = %dx%d, want 12x16", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestUnknownColorModelPropagation(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Use a MatType that does not map to known model in NewMat() metadata rule.
	// CV16U (depth-only) keeps model Unknown by design.
	m, err := r.NewMat(8, 8, CV16U)
	if err != nil {
		t.Fatalf("NewMat CV16U: %v", err)
	}
	defer m.Close()

	model, err := m.ColorModel()
	if err != nil {
		t.Fatalf("ColorModel: %v", err)
	}
	if model != Unknown {
		t.Fatalf("model = %v, want Unknown", model)
	}

	known, err := m.IsColorKnown()
	if err != nil {
		t.Fatalf("IsColorKnown: %v", err)
	}
	if known {
		t.Fatalf("IsColorKnown = true, want false")
	}
}

func TestUnknownPropagationChain(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()
	src.model = Unknown

	blur, err := r.Blur(src, Size{Width: 3, Height: 3})
	if err != nil {
		t.Fatalf("Blur: %v", err)
	}
	defer blur.Close()
	bm, _ := blur.ColorModel()
	if bm != Unknown {
		t.Fatalf("Blur model = %v, want Unknown", bm)
	}

	resized, err := r.NewMat(32, 32, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat resized: %v", err)
	}
	defer resized.Close()
	if err := r.Resize(blur, resized, Size{Width: 32, Height: 32}); err != nil {
		t.Fatalf("Resize: %v", err)
	}
	rm, _ := resized.ColorModel()
	if rm != Unknown {
		t.Fatalf("Resize model = %v, want Unknown", rm)
	}

	gray, err := r.NewMat(32, 32, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat gray: %v", err)
	}
	defer gray.Close()
	if err := r.CvtColor(resized, gray, ColorBGR2Gray); err != nil {
		t.Fatalf("CvtColor: %v", err)
	}
	gm, _ := gray.ColorModel()
	if gm != Gray {
		t.Fatalf("CvtColor output model = %v, want Gray", gm)
	}
}

func TestStrictColorValidationUnknownPassthrough(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	prev := StrictColorValidation()
	defer SetStrictColorValidation(prev)
	SetStrictColorValidation(true)

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Unknown model source: use valid 3-channel Mat, then mark metadata Unknown.
	src, err := r.NewMat(10, 10, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()
	src.model = Unknown

	dst, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	// Strict validation should not block Unknown model source
	err = r.CvtColor(src, dst, ColorBGR2Gray)
	if err != nil {
		t.Fatalf("strict unknown passthrough expected success, got: %v", err)
	}

	m, err := dst.ColorModel()
	if err != nil {
		t.Fatalf("ColorModel: %v", err)
	}
	if m != Gray {
		t.Fatalf("dst model = %v, want Gray", m)
	}
}

func TestCvtColorOutputModelDeterministic(t *testing.T) {
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}

	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// BGR -> RGBA
	bgr, err := r.NewMat(12, 12, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat bgr: %v", err)
	}
	defer bgr.Close()

	rgba, err := r.NewMat(12, 12, CV8UC4)
	if err != nil {
		t.Fatalf("NewMat rgba: %v", err)
	}
	defer rgba.Close()

	if err := r.CvtColor(bgr, rgba, ColorBGR2BGRA); err != nil {
		t.Fatalf("CvtColor BGR2BGRA: %v", err)
	}
	m, _ := rgba.ColorModel()
	if m != RGBA {
		t.Fatalf("rgba model = %v, want RGBA", m)
	}

	// RGBA -> BGR
	bgr2, err := r.NewMat(12, 12, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat bgr2: %v", err)
	}
	defer bgr2.Close()

	if err := r.CvtColor(rgba, bgr2, ColorRGBA2BGR); err != nil {
		t.Fatalf("CvtColor RGBA2BGR: %v", err)
	}
	m, _ = bgr2.ColorModel()
	if m != BGR {
		t.Fatalf("bgr2 model = %v, want BGR", m)
	}
}

// createTestPNG writes a simple solid-color PNG to disk.
func createTestPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Fill with red
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			off := y*img.Stride + x*4
			img.Pix[off+0] = 255 // R
			img.Pix[off+1] = 0   // G
			img.Pix[off+2] = 0   // B
			img.Pix[off+3] = 255 // A
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create test image: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode test image: %v", err)
	}
}
