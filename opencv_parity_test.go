package opencv

import (
	"context"
	"image/color"
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Batch 1: Filtering, gradient, morphology, drawing, histogram
// ---------------------------------------------------------------------------

func TestDLLMedianBlur(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	dst, err := r.MedianBlur(src, 5)
	if err != nil {
		t.Fatalf("MedianBlur: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("MedianBlur dims = %dx%d, want 100x100", rows, cols)
	}
	m, _ := dst.ColorModel()
	if m != BGR {
		t.Errorf("MedianBlur model = %v, want BGR", m)
	}
}

func TestDLLFlip(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 80, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	// Flip horizontal
	dst, err := r.NewMat(50, 80, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat dst: %v", err)
	}
	defer dst.Close()

	if err := r.Flip(src, dst, FlipHorizontal); err != nil {
		t.Fatalf("Flip horizontal: %v", err)
	}
	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 80 {
		t.Errorf("Flip dims = %dx%d, want 50x80", rows, cols)
	}
	m, _ := dst.ColorModel()
	if m != BGR {
		t.Errorf("Flip model = %v, want BGR", m)
	}

	// Flip vertical
	if err := r.Flip(src, dst, FlipVertical); err != nil {
		t.Fatalf("Flip vertical: %v", err)
	}

	// Flip both
	if err := r.Flip(src, dst, FlipBoth); err != nil {
		t.Fatalf("Flip both: %v", err)
	}
}

func TestDLLSobel(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	dst, err := r.Sobel(src, CV8U, 1, 0, 3, 1, 0)
	if err != nil {
		t.Fatalf("Sobel: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("Sobel dims = %dx%d, want 64x64", rows, cols)
	}
}

func TestDLLLaplacian(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	dst, err := r.Laplacian(src, CV8U, 3, 1, 0)
	if err != nil {
		t.Fatalf("Laplacian: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("Laplacian dims = %dx%d, want 64x64", rows, cols)
	}
}

func TestDLLErode(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	kernel, err := r.GetStructuringElement(MorphRect, Size{Width: 3, Height: 3})
	if err != nil {
		t.Fatalf("GetStructuringElement: %v", err)
	}
	defer kernel.Close()

	dst, err := r.Erode(src, kernel, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("Erode: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("Erode dims = %dx%d, want 64x64", rows, cols)
	}
}

func TestDLLDilate(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	kernel, err := r.GetStructuringElement(MorphRect, Size{Width: 3, Height: 3})
	if err != nil {
		t.Fatalf("GetStructuringElement: %v", err)
	}
	defer kernel.Close()

	dst, err := r.Dilate(src, kernel, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("Dilate: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("Dilate dims = %dx%d, want 64x64", rows, cols)
	}
}

func TestDLLMorphologyEx(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	kernel, err := r.GetStructuringElement(MorphEllipse, Size{Width: 5, Height: 5})
	if err != nil {
		t.Fatalf("GetStructuringElement: %v", err)
	}
	defer kernel.Close()

	// Test Open
	dst, err := r.MorphologyEx(src, MorphOpen, kernel, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("MorphologyEx Open: %v", err)
	}
	dst.Close()

	// Test Close
	dst, err = r.MorphologyEx(src, MorphClose, kernel, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("MorphologyEx Close: %v", err)
	}
	dst.Close()

	// Test Gradient
	dst, err = r.MorphologyEx(src, MorphGradient, kernel, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("MorphologyEx Gradient: %v", err)
	}
	dst.Close()
}

func TestDLLErodeNilKernel(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(32, 32, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	// nil kernel = default 3x3 rect
	dst, err := r.Erode(src, nil, Point{X: -1, Y: -1}, 1)
	if err != nil {
		t.Fatalf("Erode nil kernel: %v", err)
	}
	dst.Close()
}

func TestDLLEqualizeHist(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	dst, err := r.EqualizeHist(src)
	if err != nil {
		t.Fatalf("EqualizeHist: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("EqualizeHist dims = %dx%d, want 64x64", rows, cols)
	}
	m, _ := dst.ColorModel()
	if m != Gray {
		t.Errorf("EqualizeHist model = %v, want Gray", m)
	}
}

func TestDLLPutText(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.NewMat(200, 400, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	err = r.PutText(img, "Hello OpenCV", Point{X: 50, Y: 100}, FontHersheySimplex, 1.0, color.RGBA{255, 255, 255, 255}, 2)
	if err != nil {
		t.Fatalf("PutText: %v", err)
	}
}

func TestDLLArrowedLine(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.NewMat(200, 200, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	err = r.ArrowedLine(img, Point{X: 10, Y: 10}, Point{X: 190, Y: 190}, color.RGBA{0, 255, 0, 255}, 2, 0.1)
	if err != nil {
		t.Fatalf("ArrowedLine: %v", err)
	}
}

func TestDLLTranspose(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 80, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	dst, err := r.Transpose(src)
	if err != nil {
		t.Fatalf("Transpose: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 80 || cols != 50 {
		t.Errorf("Transpose dims = %dx%d, want 80x50", rows, cols)
	}
}

// ---------------------------------------------------------------------------
// Batch 2: Contours + Hough
// ---------------------------------------------------------------------------

func TestDLLFindContours(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Create a binary image with a white rectangle
	img, err := r.NewMat(100, 100, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	// Draw a filled white rectangle
	err = r.Rectangle(img, Rect{X: 20, Y: 20, Width: 60, Height: 60}, color.Gray{Y: 255}, -1)
	if err != nil {
		t.Fatalf("Rectangle: %v", err)
	}

	contours, err := r.FindContours(img, RetrievalExternal, ChainApproxSimple)
	if err != nil {
		t.Fatalf("FindContours: %v", err)
	}

	if len(contours) == 0 {
		t.Fatal("FindContours returned 0 contours, expected at least 1")
	}

	// The contour should retain at least one vertex across platform-specific simplification.
	if len(contours[0]) == 0 {
		t.Error("First contour has 0 points, want >= 1")
	}
}

func TestDLLContourArea(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Square contour: (0,0), (10,0), (10,10), (0,10)
	contour := []Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}}
	area, err := r.ContourArea(contour)
	if err != nil {
		t.Fatalf("ContourArea: %v", err)
	}

	if math.Abs(area-100) > 1 {
		t.Errorf("ContourArea = %f, want ~100", area)
	}
}

func TestDLLArcLength(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Square contour
	contour := []Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}}
	length, err := r.ArcLength(contour, true)
	if err != nil {
		t.Fatalf("ArcLength: %v", err)
	}

	if math.Abs(length-40) > 1 {
		t.Errorf("ArcLength = %f, want ~40", length)
	}
}

func TestDLLBoundingRect(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	contour := []Point{{X: 5, Y: 5}, {X: 15, Y: 5}, {X: 15, Y: 20}, {X: 5, Y: 20}}
	rect, err := r.BoundingRect(contour)
	if err != nil {
		t.Fatalf("BoundingRect: %v", err)
	}

	if rect.X != 5 || rect.Y != 5 || rect.Width != 11 || rect.Height != 16 {
		t.Errorf("BoundingRect = %+v, want {5,5,11,16}", rect)
	}
}

func TestDLLMoments(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	contour := []Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}}
	moments, err := r.Moments(contour, false)
	if err != nil {
		t.Fatalf("Moments: %v", err)
	}

	if moments.M00 < 99 || moments.M00 > 101 {
		t.Errorf("Moments M00 = %f, want ~100", moments.M00)
	}
}

func TestDLLDrawContours(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	contours := [][]Point{
		{{X: 10, Y: 10}, {X: 90, Y: 10}, {X: 90, Y: 90}, {X: 10, Y: 90}},
	}
	err = r.DrawContours(img, contours, 0, color.RGBA{0, 255, 0, 255}, 2)
	if err != nil {
		t.Fatalf("DrawContours: %v", err)
	}
}

func TestDLLMinEnclosingCircle(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	contour := []Point{{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 20}, {X: 0, Y: 20}}
	center, radius, err := r.MinEnclosingCircle(contour)
	if err != nil {
		t.Fatalf("MinEnclosingCircle: %v", err)
	}

	if radius <= 0 {
		t.Errorf("MinEnclosingCircle radius = %f, want > 0", radius)
	}
	t.Logf("MinEnclosingCircle: center=(%d,%d) radius=%f", center.X, center.Y, radius)
}

func TestDLLHoughLinesP(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Create image with lines
	img, err := r.NewMat(100, 100, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer img.Close()

	// Draw a line
	err = r.Line(img, Point{X: 10, Y: 50}, Point{X: 90, Y: 50}, color.Gray{Y: 255}, 2)
	if err != nil {
		t.Fatalf("Line: %v", err)
	}

	lines, err := r.HoughLinesP(img, 1, math.Pi/180, 50, 10, 5)
	if err != nil {
		t.Fatalf("HoughLinesP: %v", err)
	}

	if len(lines) == 0 {
		t.Log("HoughLinesP: no lines detected (may be OK depending on threshold)")
	}
}

// ---------------------------------------------------------------------------
// Batch 3: Warp, core ops
// ---------------------------------------------------------------------------

func TestDLLWarpAffine(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat src: %v", err)
	}
	defer src.Close()

	// Get a rotation matrix
	M, err := r.GetRotationMatrix2D(Point{X: 50, Y: 50}, 45, 1.0)
	if err != nil {
		t.Fatalf("GetRotationMatrix2D: %v", err)
	}
	defer M.Close()

	dst, err := r.WarpAffine(src, M, Size{Width: 100, Height: 100})
	if err != nil {
		t.Fatalf("WarpAffine: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("WarpAffine dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLGetRotationMatrix2D(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	M, err := r.GetRotationMatrix2D(Point{X: 50, Y: 50}, 90, 1.0)
	if err != nil {
		t.Fatalf("GetRotationMatrix2D: %v", err)
	}
	defer M.Close()

	rows, _ := M.Rows()
	cols, _ := M.Cols()
	if rows != 2 || cols != 3 {
		t.Errorf("Rotation matrix dims = %dx%d, want 2x3", rows, cols)
	}
}

func TestDLLGetAffineTransform(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src := [3]Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 0, Y: 100}}
	dst := [3]Point{{X: 10, Y: 10}, {X: 110, Y: 10}, {X: 10, Y: 110}}

	M, err := r.GetAffineTransform(src, dst)
	if err != nil {
		t.Fatalf("GetAffineTransform: %v", err)
	}
	defer M.Close()

	rows, _ := M.Rows()
	cols, _ := M.Cols()
	if rows != 2 || cols != 3 {
		t.Errorf("Affine matrix dims = %dx%d, want 2x3", rows, cols)
	}
}

func TestDLLBitwiseOps(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src1, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat src1: %v", err)
	}
	defer src1.Close()

	src2, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat src2: %v", err)
	}
	defer src2.Close()

	// AND
	and, err := r.BitwiseAnd(src1, src2)
	if err != nil {
		t.Fatalf("BitwiseAnd: %v", err)
	}
	and.Close()

	// OR
	or, err := r.BitwiseOr(src1, src2)
	if err != nil {
		t.Fatalf("BitwiseOr: %v", err)
	}
	or.Close()

	// XOR
	xor, err := r.BitwiseXor(src1, src2)
	if err != nil {
		t.Fatalf("BitwiseXor: %v", err)
	}
	xor.Close()

	// NOT
	not, err := r.BitwiseNot(src1)
	if err != nil {
		t.Fatalf("BitwiseNot: %v", err)
	}
	not.Close()
}

func TestDLLArithmeticOps(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src1, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat src1: %v", err)
	}
	defer src1.Close()

	src2, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat src2: %v", err)
	}
	defer src2.Close()

	// Add
	add, err := r.Add(src1, src2)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	add.Close()

	// Subtract
	sub, err := r.Subtract(src1, src2)
	if err != nil {
		t.Fatalf("Subtract: %v", err)
	}
	sub.Close()

	// Multiply
	mul, err := r.Multiply(src1, src2, 1.0)
	if err != nil {
		t.Fatalf("Multiply: %v", err)
	}
	mul.Close()

	// Divide
	div, err := r.Divide(src1, src2, 1.0)
	if err != nil {
		t.Fatalf("Divide: %v", err)
	}
	div.Close()

	// AbsDiff
	abs, err := r.AbsDiff(src1, src2)
	if err != nil {
		t.Fatalf("AbsDiff: %v", err)
	}
	abs.Close()
}

func TestDLLMinMaxLoc(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Use Zeros() to guarantee zero-initialized mat
	src, err := r.Zeros(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("Zeros: %v", err)
	}
	defer src.Close()

	result, err := r.MinMaxLoc(src)
	if err != nil {
		t.Fatalf("MinMaxLoc: %v", err)
	}

	// Zero mat: min=max=0
	if result.MinVal != 0 || result.MaxVal != 0 {
		t.Errorf("MinMaxLoc on zero mat: min=%f max=%f, want 0,0", result.MinVal, result.MaxVal)
	}
}

func TestDLLCountNonZero(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Use Zeros() to guarantee zero-initialized mat
	src, err := r.Zeros(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("Zeros: %v", err)
	}
	defer src.Close()

	count, err := r.CountNonZero(src)
	if err != nil {
		t.Fatalf("CountNonZero: %v", err)
	}
	if count != 0 {
		t.Errorf("CountNonZero on zero mat = %d, want 0", count)
	}
}

func TestDLLNormalize(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(64, 64, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	dst, err := r.Normalize(src, 0, 255, NormMinMax)
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 64 || cols != 64 {
		t.Errorf("Normalize dims = %dx%d, want 64x64", rows, cols)
	}
}

func TestDLLSplit(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(32, 32, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	channels, err := r.Split(src)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	if len(channels) != 3 {
		t.Fatalf("Split returned %d channels, want 3", len(channels))
	}
	for i, ch := range channels {
		defer ch.Close()
		rows, _ := ch.Rows()
		cols, _ := ch.Cols()
		if rows != 32 || cols != 32 {
			t.Errorf("Split channel %d dims = %dx%d, want 32x32", i, rows, cols)
		}
		m, _ := ch.ColorModel()
		if m != Gray {
			t.Errorf("Split channel %d model = %v, want Gray", i, m)
		}
	}
}

func TestDLLMerge(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	ch1, err := r.NewMat(16, 16, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat ch1: %v", err)
	}
	defer ch1.Close()

	ch2, err := r.NewMat(16, 16, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat ch2: %v", err)
	}
	defer ch2.Close()

	ch3, err := r.NewMat(16, 16, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat ch3: %v", err)
	}
	defer ch3.Close()

	merged, err := r.Merge([]*Mat{ch1, ch2, ch3})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	defer merged.Close()

	rows, _ := merged.Rows()
	cols, _ := merged.Cols()
	if rows != 16 || cols != 16 {
		t.Errorf("Merge dims = %dx%d, want 16x16", rows, cols)
	}
	ch, _ := merged.Channels()
	if ch != 3 {
		t.Errorf("Merge channels = %d, want 3", ch)
	}
}

func TestDLLMeanStdDev(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer src.Close()

	result, err := r.MeanStdDev(src)
	if err != nil {
		t.Fatalf("MeanStdDev: %v", err)
	}
	defer result.Mean.Close()
	defer result.StdDev.Close()
}

func TestDLLZerosOnesEye(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	// Zeros
	z, err := r.Zeros(10, 10, CV8UC1)
	if err != nil {
		t.Fatalf("Zeros: %v", err)
	}
	defer z.Close()
	count, _ := r.CountNonZero(z)
	if count != 0 {
		t.Errorf("Zeros CountNonZero = %d, want 0", count)
	}

	// Ones
	o, err := r.Ones(5, 5, CV8UC1)
	if err != nil {
		t.Fatalf("Ones: %v", err)
	}
	defer o.Close()

	// Eye
	e, err := r.Eye(3, 3, CV8UC1)
	if err != nil {
		t.Fatalf("Eye: %v", err)
	}
	defer e.Close()
	rows, _ := e.Rows()
	cols, _ := e.Cols()
	if rows != 3 || cols != 3 {
		t.Errorf("Eye dims = %dx%d, want 3x3", rows, cols)
	}
}

func TestDLLMatRowColRegionReshape(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	m, err := r.NewMat(10, 20, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	// Row
	row, err := m.Row(5)
	if err != nil {
		t.Fatalf("Row: %v", err)
	}
	defer row.Close()
	rows, _ := row.Rows()
	cols, _ := row.Cols()
	if rows != 1 || cols != 20 {
		t.Errorf("Row dims = %dx%d, want 1x20", rows, cols)
	}

	// Col
	col, err := m.Col(3)
	if err != nil {
		t.Fatalf("Col: %v", err)
	}
	defer col.Close()
	rows, _ = col.Rows()
	cols, _ = col.Cols()
	if rows != 10 || cols != 1 {
		t.Errorf("Col dims = %dx%d, want 10x1", rows, cols)
	}

	// Region
	region, err := m.Region(Rect{X: 2, Y: 3, Width: 5, Height: 4})
	if err != nil {
		t.Fatalf("Region: %v", err)
	}
	defer region.Close()
	rows, _ = region.Rows()
	cols, _ = region.Cols()
	if rows != 4 || cols != 5 {
		t.Errorf("Region dims = %dx%d, want 4x5", rows, cols)
	}
	rm, _ := region.ColorModel()
	if rm != BGR {
		t.Errorf("Region model = %v, want BGR", rm)
	}

	// Reshape
	reshaped, err := m.Reshape(1, 0)
	if err != nil {
		t.Fatalf("Reshape: %v", err)
	}
	defer reshaped.Close()
	ch, _ := reshaped.Channels()
	if ch != 1 {
		t.Errorf("Reshape channels = %d, want 1", ch)
	}
}

func TestDLLMatTotal(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	m, err := r.NewMat(10, 20, CV8UC3)
	if err != nil {
		t.Fatalf("NewMat: %v", err)
	}
	defer m.Close()

	total := m.Total()
	if total != 200 {
		t.Errorf("Total = %d, want 200", total)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestRuntime(t *testing.T) *Runtime {
	t.Helper()
	path := dllPath()
	if path == "" {
		t.Skip("goopencv.dll not found")
	}
	r, err := New(context.Background(), WithDLL(path))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

// ---------------------------------------------------------------------------
// Batch 5: New imgproc APIs (BilateralFilter, InRange, Rotate, Hconcat, Vconcat, etc.)
// ---------------------------------------------------------------------------

func TestDLLBilateralFilter(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.BilateralFilter(src, 5, 50, 50)
	if err != nil { t.Fatalf("BilateralFilter: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("BilateralFilter dims = %dx%d, want 50x50", rows, cols)
	}
	m, _ := dst.ColorModel()
	if m != BGR { t.Errorf("model = %v, want BGR", m) }
}

func TestDLLInRange(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	lower := Scalar{V0: 100, V1: 0, V2: 0, V3: 0}
	upper := Scalar{V0: 255, V1: 100, V2: 100, V3: 255}
	dst, err := r.InRange(src, lower, upper)
	if err != nil { t.Fatalf("InRange: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("InRange dims = %dx%d, want 100x100", rows, cols)
	}
	ch, _ := dst.Channels()
	if ch != 1 { t.Errorf("InRange channels = %d, want 1", ch) }
	m, _ := dst.ColorModel()
	if m != Gray { t.Errorf("model = %v, want Gray", m) }
}

func TestDLLMatchTemplate(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.NewMat(100, 100, CV8UC3)
	if err != nil { t.Fatalf("NewMat img: %v", err) }
	defer img.Close()

	tmpl, err := r.NewMat(20, 20, CV8UC3)
	if err != nil { t.Fatalf("NewMat tmpl: %v", err) }
	defer tmpl.Close()

	result, err := r.MatchTemplate(img, tmpl, TMSqDiff)
	if err != nil { t.Fatalf("MatchTemplate: %v", err) }
	defer result.Close()

	rows, _ := result.Rows()
	cols, _ := result.Cols()
	if rows != 81 || cols != 81 {
		t.Errorf("MatchTemplate result = %dx%d, want 81x81", rows, cols)
	}
}

func TestDLLCalcHist(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(100, 100, CV8UC1)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	hist, err := r.CalcHist(src, 256, 0, 256)
	if err != nil { t.Fatalf("CalcHist: %v", err) }
	defer hist.Close()

	rows, _ := hist.Rows()
	if rows != 256 { t.Errorf("CalcHist rows = %d, want 256", rows) }
}

func TestDLLConnectedComponents(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.Zeros(100, 100, CV8UC1)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer src.Close()

	// Draw two white blobs
	r.Circle(src, Point{X: 25, Y: 25}, 10, color.Gray{Y: 255}, -1)
	r.Circle(src, Point{X: 75, Y: 75}, 10, color.Gray{Y: 255}, -1)

	labels, err := r.ConnectedComponents(src, 8)
	if err != nil { t.Fatalf("ConnectedComponents: %v", err) }
	defer labels.Close()

	rows, _ := labels.Rows()
	cols, _ := labels.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("ConnectedComponents dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLDistanceTransform(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.Zeros(100, 100, CV8UC1)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer src.Close()

	r.Circle(src, Point{X: 50, Y: 50}, 30, color.Gray{Y: 255}, -1)

	dst, err := r.DistanceTransform(src, DistL2, 5)
	if err != nil { t.Fatalf("DistanceTransform: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 100 || cols != 100 {
		t.Errorf("DistanceTransform dims = %dx%d, want 100x100", rows, cols)
	}
}

func TestDLLCopyMakeBorder(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.CopyMakeBorder(src, 10, 10, 10, 10, BorderConstant, Scalar{})
	if err != nil { t.Fatalf("CopyMakeBorder: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 70 || cols != 70 {
		t.Errorf("CopyMakeBorder dims = %dx%d, want 70x70", rows, cols)
	}
}

func TestDLLRotate(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(100, 200, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.Rotate(src, Rotate90Clockwise)
	if err != nil { t.Fatalf("Rotate: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 200 || cols != 100 {
		t.Errorf("Rotate90 dims = %dx%d, want 200x100", rows, cols)
	}
}

func TestDLLHconcat(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src1, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat1: %v", err) }
	defer src1.Close()

	src2, err := r.NewMat(50, 30, CV8UC3)
	if err != nil { t.Fatalf("NewMat2: %v", err) }
	defer src2.Close()

	dst, err := r.Hconcat(src1, src2)
	if err != nil { t.Fatalf("Hconcat: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 80 {
		t.Errorf("Hconcat dims = %dx%d, want 50x80", rows, cols)
	}
}

func TestDLLVconcat(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src1, err := r.NewMat(50, 80, CV8UC3)
	if err != nil { t.Fatalf("NewMat1: %v", err) }
	defer src1.Close()

	src2, err := r.NewMat(30, 80, CV8UC3)
	if err != nil { t.Fatalf("NewMat2: %v", err) }
	defer src2.Close()

	dst, err := r.Vconcat(src1, src2)
	if err != nil { t.Fatalf("Vconcat: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 80 || cols != 80 {
		t.Errorf("Vconcat dims = %dx%d, want 80x80", rows, cols)
	}
}

func TestDLLLUT(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC1)
	if err != nil { t.Fatalf("NewMat src: %v", err) }
	defer src.Close()

	// Create identity LUT (256 entries, CV8UC1)
	lut, err := r.NewMat(1, 256, CV8UC1)
	if err != nil { t.Fatalf("NewMat lut: %v", err) }
	defer lut.Close()
	for i := 0; i < 256; i++ {
		lut.SetByte(0, int32(i), 0, uint8(i))
	}

	dst, err := r.LUT(src, lut)
	if err != nil { t.Fatalf("LUT: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("LUT dims = %dx%d, want 50x50", rows, cols)
	}
}

func TestDLLIntegral(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC1)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	sum, err := r.Integral(src)
	if err != nil { t.Fatalf("Integral: %v", err) }
	defer sum.Close()

	rows, _ := sum.Rows()
	cols, _ := sum.Cols()
	if rows != 51 || cols != 51 {
		t.Errorf("Integral dims = %dx%d, want 51x51", rows, cols)
	}
}

func TestDLLGetPerspectiveTransform(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src := [4]Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 100}, {X: 0, Y: 100}}
	dst := [4]Point{{X: 10, Y: 10}, {X: 90, Y: 10}, {X: 90, Y: 90}, {X: 10, Y: 90}}

	M, err := r.GetPerspectiveTransform(src, dst)
	if err != nil { t.Fatalf("GetPerspectiveTransform: %v", err) }
	defer M.Close()

	rows, _ := M.Rows()
	cols, _ := M.Cols()
	if rows != 3 || cols != 3 {
		t.Errorf("GetPerspectiveTransform dims = %dx%d, want 3x3", rows, cols)
	}
}

func TestDLLFillConvexPoly(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.Zeros(100, 100, CV8UC3)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer img.Close()

	pts := []Point{{X: 50, Y: 10}, {X: 90, Y: 90}, {X: 10, Y: 90}}
	err = r.FillConvexPoly(img, pts, color.RGBA{0, 255, 0, 255}, Line8)
	if err != nil { t.Fatalf("FillConvexPoly: %v", err) }
}

// ---------------------------------------------------------------------------
// Batch 6: ConvertModel
// ---------------------------------------------------------------------------

func TestDLLConvertModel(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	img, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer img.Close()

	// BGR -> Gray
	gray, err := r.ConvertModel(img, Gray)
	if err != nil { t.Fatalf("ConvertModel BGR->Gray: %v", err) }
	defer gray.Close()

	ch, _ := gray.Channels()
	if ch != 1 { t.Errorf("Gray channels = %d, want 1", ch) }
	m, _ := gray.ColorModel()
	if m != Gray { t.Errorf("Gray model = %v, want Gray", m) }

	// Gray -> BGR
	bgr, err := r.ConvertModel(gray, BGR)
	if err != nil { t.Fatalf("ConvertModel Gray->BGR: %v", err) }
	defer bgr.Close()

	ch2, _ := bgr.Channels()
	if ch2 != 3 { t.Errorf("BGR channels = %d, want 3", ch2) }
	m2, _ := bgr.ColorModel()
	if m2 != BGR { t.Errorf("BGR model = %v, want BGR", m2) }

	// Same model -> clone
	clone, err := r.ConvertModel(bgr, BGR)
	if err != nil { t.Fatalf("ConvertModel same: %v", err) }
	defer clone.Close()
}

// ---------------------------------------------------------------------------
// Batch 7: Photo APIs
// ---------------------------------------------------------------------------

func TestDLLFastNlMeansDenoising(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC1)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.FastNlMeansDenoising(src, 10, 7, 21)
	if err != nil { t.Fatalf("FastNlMeansDenoising: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("dims = %dx%d, want 50x50", rows, cols)
	}
}

func TestDLLFastNlMeansDenoisingColored(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.FastNlMeansDenoisingColored(src, 10, 10, 7, 21)
	if err != nil { t.Fatalf("FastNlMeansDenoisingColored: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("dims = %dx%d, want 50x50", rows, cols)
	}
	m, _ := dst.ColorModel()
	if m != BGR { t.Errorf("model = %v, want BGR", m) }
}

func TestDLLDetailEnhance(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.DetailEnhance(src, 10, 0.15)
	if err != nil { t.Fatalf("DetailEnhance: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("dims = %dx%d, want 50x50", rows, cols)
	}
}

func TestDLLEdgePreservingFilter(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.EdgePreservingFilter(src, RecursFilter, 60, 0.4)
	if err != nil { t.Fatalf("EdgePreservingFilter: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("dims = %dx%d, want 50x50", rows, cols)
	}
}

func TestDLLStylization(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.NewMat(50, 50, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer src.Close()

	dst, err := r.Stylization(src, 60, 0.45)
	if err != nil { t.Fatalf("Stylization: %v", err) }
	defer dst.Close()

	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	if rows != 50 || cols != 50 {
		t.Errorf("dims = %dx%d, want 50x50", rows, cols)
	}
}

// ---------------------------------------------------------------------------
// Batch 8: Features2d APIs
// ---------------------------------------------------------------------------

func TestDLLFAST(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.Zeros(200, 200, CV8UC1)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer src.Close()

	// Draw features for detection
	r.Rectangle(src, Rect{X: 50, Y: 50, Width: 100, Height: 100}, color.Gray{Y: 255}, 2)

	kps, err := r.FAST(src, 50, true)
	if err != nil { t.Fatalf("FAST: %v", err) }

	if len(kps) == 0 {
		t.Log("FAST: no keypoints found (possible with simple image)")
	} else {
		t.Logf("FAST: found %d keypoints", len(kps))
		for i, kp := range kps {
			if i >= 3 { break }
			t.Logf("  kp[%d]: x=%.1f y=%.1f size=%.1f", i, kp.X, kp.Y, kp.Size)
		}
	}
}

func TestDLLORBDetectCompute(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.Zeros(200, 200, CV8UC1)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer src.Close()

	r.Rectangle(src, Rect{X: 50, Y: 50, Width: 100, Height: 100}, color.Gray{Y: 255}, 2)
	r.Circle(src, Point{X: 100, Y: 100}, 30, color.Gray{Y: 200}, -1)

	kps, desc, err := r.ORBDetectCompute(src, 500, 1.2, 8)
	if err != nil { t.Fatalf("ORBDetectCompute: %v", err) }
	defer desc.Close()

	dr, _ := desc.Rows()
	dc, _ := desc.Cols()
	t.Logf("ORB: %d keypoints, descriptor %dx%d", len(kps), dr, dc)
	if len(kps) > 0 {
		kp := kps[0]
		t.Logf("  first kp: x=%.1f y=%.1f size=%.1f angle=%.1f", kp.X, kp.Y, kp.Size, kp.Angle)
	}
}

func TestDLLBFMatch(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	src, err := r.Zeros(200, 200, CV8UC1)
	if err != nil { t.Fatalf("Zeros: %v", err) }
	defer src.Close()

	r.Rectangle(src, Rect{X: 50, Y: 50, Width: 100, Height: 100}, color.Gray{Y: 255}, 2)
	r.Circle(src, Point{X: 100, Y: 100}, 30, color.Gray{Y: 200}, -1)

	kps1, desc1, err := r.ORBDetectCompute(src, 500, 1.2, 8)
	if err != nil { t.Fatalf("ORB1: %v", err) }
	defer desc1.Close()

	if len(kps1) < 2 {
		t.Skip("Not enough keypoints for BFMatch test")
	}

	kps2, desc2, err := r.ORBDetectCompute(src, 500, 1.2, 8)
	if err != nil { t.Fatalf("ORB2: %v", err) }
	defer desc2.Close()

	matches, err := r.BFMatch(desc1, desc2, NormHamming)
	if err != nil { t.Fatalf("BFMatch: %v", err) }

	t.Logf("BFMatch: %d matches out of %d and %d keypoints", len(matches), len(kps1), len(kps2))
	if len(matches) > 0 {
		m := matches[0]
		t.Logf("  first match: query=%d train=%d dist=%.1f", m.QueryIdx, m.TrainIdx, m.Distance)
	}
}

// ---------------------------------------------------------------------------
// Batch 9: Core extras (Diag, AtByte, SetByte)
// ---------------------------------------------------------------------------

func TestDLLMatDiag(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	eye, err := r.Eye(3, 3, CV8UC1)
	if err != nil { t.Fatalf("Eye: %v", err) }
	defer eye.Close()

	diag, err := eye.Diag()
	if err != nil { t.Fatalf("Diag: %v", err) }
	defer diag.Close()

	rows, _ := diag.Rows()
	cols, _ := diag.Cols()
	if rows != 3 || cols != 1 {
		t.Errorf("Diag dims = %dx%d, want 3x1", rows, cols)
	}
}

func TestDLLMatAtSetByte(t *testing.T) {
	r := newTestRuntime(t)
	defer r.Close()

	mat, err := r.NewMat(10, 10, CV8UC3)
	if err != nil { t.Fatalf("NewMat: %v", err) }
	defer mat.Close()

	// Set a pixel
	err = mat.SetByte(5, 5, 0, 42)
	if err != nil { t.Fatalf("SetByte: %v", err) }

	// Get it back
	val, err := mat.AtByte(5, 5, 0)
	if err != nil { t.Fatalf("AtByte: %v", err) }

	if val != 42 {
		t.Errorf("AtByte = %d, want 42", val)
	}
}
