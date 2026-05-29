// go-opencv examples
//
// Run:  go run .
//
// Or run individual demos:
//
//	go run . -demo=basic     # Read image, print info
//	go run . -demo=convert   # Color space conversion
//	go run . -demo=resize    # Image scaling
//	go run . -demo=edges     # Canny edge detection
//	go run . -demo=filter    # GaussianBlur + Sobel + Morphology
//	go run . -demo=contours  # Find and draw contours
//	go run . -demo=warp      # Rotation via WarpAffine
//	go run . -demo=io        # IMRead + IMWrite roundtrip
//	go run . -demo=all       # Full pipeline (default)

package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"os"

	"github.com/brainplusplus/go-opencv"
)

func main() {
	demo := flag.String("demo", "all", "demo to run: basic, convert, resize, edges, filter, contours, warp, io, all")
	flag.Parse()

	r, err := opencv.New(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "opencv.New: %v\n", err)
		fmt.Fprintln(os.Stderr, "Tip: build goopencv.dll first with: build-tools\\build-goopencv.bat")
		os.Exit(1)
	}
	defer r.Close()

	fmt.Println("=== go-opencv Examples ===")
	fmt.Printf("Library version: %s\n\n", opencv.Version)

	switch *demo {
	case "basic":
		demoBasic(r)
	case "convert":
		demoConvert(r)
	case "resize":
		demoResize(r)
	case "edges":
		demoEdges(r)
	case "filter":
		demoFilter(r)
	case "contours":
		demoContours(r)
	case "warp":
		demoWarp(r)
	case "io":
		demoIO(r)
	default:
		demoBasic(r)
		fmt.Println()
		demoConvert(r)
		fmt.Println()
		demoResize(r)
		fmt.Println()
		demoEdges(r)
		fmt.Println()
		demoFilter(r)
		fmt.Println()
		demoContours(r)
		fmt.Println()
		demoWarp(r)
		fmt.Println()
		demoIO(r)
	}
}

// ── 01. Basic: read image and print properties ─────────────────────────────

func demoBasic(r *opencv.Runtime) {
	fmt.Println("--- 01. Basic: Image Info ---")

	img, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer img.Close()

	rows, _ := img.Rows()
	cols, _ := img.Cols()
	typ, _ := img.Type()
	ch, _ := img.Channels()
	model, _ := img.ColorModel()

	fmt.Printf("  File:      sample.png\n")
	fmt.Printf("  Size:      %d x %d pixels\n", cols, rows)
	fmt.Printf("  Type:      %d (CV_8UC3 = BGR)\n", typ)
	fmt.Printf("  Channels:  %d\n", ch)
	fmt.Printf("  ColorModel: %v\n", model)
	fmt.Println("  OK: Image loaded")
}

// ── 02. Color Conversion ────────────────────────────────────────────────────

func demoConvert(r *opencv.Runtime) {
	fmt.Println("--- 02. Color: BGR -> Gray ---")

	src, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	gray, err := r.NewMat(300, 400, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer gray.Close()

	if err := r.CvtColor(src, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	dstType, _ := gray.Type()
	fmt.Printf("  Input type:  BGR (CV_8UC3)\n")
	fmt.Printf("  Output type: %d (CV_8UC1 = Gray)\n", dstType)
	fmt.Println("  OK: Color conversion")
}

// ── 03. Resize ──────────────────────────────────────────────────────────────

func demoResize(r *opencv.Runtime) {
	fmt.Println("--- 03. Resize: 400x300 -> 200x150 ---")

	src, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	dst, err := r.NewMat(150, 200, opencv.CV8UC3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer dst.Close()

	if err := r.Resize(src, dst, opencv.Size{Width: 200, Height: 150}); err != nil {
		fmt.Fprintf(os.Stderr, "  Resize: %v\n", err)
		return
	}

	drows, _ := dst.Rows()
	dcols, _ := dst.Cols()
	fmt.Printf("  Before: 400 x 300\n")
	fmt.Printf("  After:  %d x %d\n", dcols, drows)
	fmt.Println("  OK: Resize")
}

// ── 04. Canny Edge Detection ────────────────────────────────────────────────

func demoEdges(r *opencv.Runtime) {
	fmt.Println("--- 04. Canny Edge Detection ---")

	src, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	gray, err := r.NewMat(300, 400, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat (gray): %v\n", err)
		return
	}
	defer gray.Close()

	if err := r.CvtColor(src, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	edges, err := r.Canny(gray, 50, 150)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Canny: %v\n", err)
		return
	}
	defer edges.Close()

	erows, _ := edges.Rows()
	ecols, _ := edges.Cols()
	fmt.Printf("  Edges: %d x %d\n", ecols, erows)
	fmt.Println("  OK: Edge detection")
}

// ── 05. Filter: Blur + Sobel + Morphology ───────────────────────────────────

func demoFilter(r *opencv.Runtime) {
	fmt.Println("--- 05. Filter Pipeline ---")

	src, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	// GaussianBlur
	blurred, err := r.GaussianBlur(src, opencv.Size{Width: 5, Height: 5}, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GaussianBlur: %v\n", err)
		return
	}
	defer blurred.Close()
	fmt.Println("  GaussianBlur: OK")

	// Convert to gray for Sobel
	gray, err := r.NewMat(300, 400, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat gray: %v\n", err)
		return
	}
	defer gray.Close()
	if err := r.CvtColor(blurred, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	// Sobel
	sobel, err := r.Sobel(gray, opencv.CV8U, 1, 0, 3, 1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Sobel: %v\n", err)
		return
	}
	defer sobel.Close()
	fmt.Println("  Sobel:        OK")

	// Erode
	kernel, err := r.GetStructuringElement(opencv.MorphRect, opencv.Size{Width: 3, Height: 3})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GetStructuringElement: %v\n", err)
		return
	}
	defer kernel.Close()

	eroded, err := r.Erode(sobel, kernel, opencv.Point{X: -1, Y: -1}, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Erode: %v\n", err)
		return
	}
	defer eroded.Close()
	fmt.Println("  Erode:        OK")

	// MorphologyEx (Close)
	closed, err := r.MorphologyEx(eroded, opencv.MorphClose, kernel, opencv.Point{X: -1, Y: -1}, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  MorphologyEx: %v\n", err)
		return
	}
	defer closed.Close()
	fmt.Println("  MorphClose:   OK")
}

// ── 06. Contours ────────────────────────────────────────────────────────────

func demoContours(r *opencv.Runtime) {
	fmt.Println("--- 06. Contours ---")

	// Create a synthetic image with shapes
	img, err := r.NewMat(200, 200, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer img.Close()

	// Draw filled shapes
	r.Rectangle(img, opencv.Rect{X: 20, Y: 20, Width: 60, Height: 60}, color.Gray{Y: 255}, -1)
	r.Circle(img, opencv.Point{X: 140, Y: 140}, 40, color.Gray{Y: 255}, -1)

	contours, err := r.FindContours(img, opencv.RetrievalExternal, opencv.ChainApproxSimple)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  FindContours: %v\n", err)
		return
	}

	fmt.Printf("  Found %d contours\n", len(contours))
	for i, c := range contours {
		area, _ := r.ContourArea(c)
		rect, _ := r.BoundingRect(c)
		fmt.Printf("  Contour %d: %d points, area=%.0f, rect=%v\n", i, len(c), area, rect)
	}

	// Draw contours on a color image
	colorImg, err := r.NewMat(200, 200, opencv.CV8UC3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat color: %v\n", err)
		return
	}
	defer colorImg.Close()

	r.DrawContours(colorImg, contours, -1, color.RGBA{0, 255, 0, 255}, 2)
	fmt.Println("  OK: Contours drawn")
}

// ── 07. Warp: Rotate Image ─────────────────────────────────────────────────

func demoWarp(r *opencv.Runtime) {
	fmt.Println("--- 07. Warp: Rotate 45 degrees ---")

	src, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	rows, _ := src.Rows()
	cols, _ := src.Cols()

	// Get rotation matrix (45 degrees around center)
	M, err := r.GetRotationMatrix2D(opencv.Point{X: int32(cols / 2), Y: int32(rows / 2)}, 45, 1.0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GetRotationMatrix2D: %v\n", err)
		return
	}
	defer M.Close()

	// Apply warp
	dst, err := r.WarpAffine(src, M, opencv.Size{Width: int32(cols), Height: int32(rows)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  WarpAffine: %v\n", err)
		return
	}
	defer dst.Close()

	drows, _ := dst.Rows()
	dcols, _ := dst.Cols()
	fmt.Printf("  Rotated %dx%d -> %dx%d\n", cols, rows, dcols, drows)
	fmt.Println("  OK: Warp affine rotation")
}

// ── 08. IO: Read + Process + Write roundtrip ────────────────────────────────

func demoIO(r *opencv.Runtime) {
	fmt.Println("--- 08. IO: Read -> Gray -> Write ---")

	// Read in BGR (default)
	img, err := r.IMRead("sample.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer img.Close()

	model, _ := img.ColorModel()
	fmt.Printf("  Read: colorModel=%v\n", model)

	// Convert to gray
	gray, err := r.NewMat(300, 400, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer gray.Close()

	if err := r.CvtColor(img, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	// Write gray image
	if err := r.IMWrite("output_gray.png", gray); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}
	fmt.Println("  Wrote: output_gray.png")

	// Read it back and verify
	back, err := r.IMRead("output_gray.png", opencv.Gray)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead back: %v\n", err)
		return
	}
	defer back.Close()

	backModel, _ := back.ColorModel()
	backRows, _ := back.Rows()
	backCols, _ := back.Cols()
	fmt.Printf("  Read back: %dx%d model=%v\n", backCols, backRows, backModel)
	fmt.Println("  OK: IO roundtrip")
}
