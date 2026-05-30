// Command examples demonstrates go-opencv features against a sample image.
//
// Usage (from examples/):
//
//	go run . -demo=basic
//	go run . -demo=io -in sample.png -out output
//	go run . -demo=all
package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	opencv "github.com/brainplusplus/go-opencv"
)

func main() {
	demo := flag.String("demo", "all", "demo to run: basic, convert, resize, edges, filter, contours, warp, io, all")
	input := flag.String("in", "sample.png", "input image path for image-based demos")
	outputDir := flag.String("out", "output", "output directory for write demos")
	flag.Parse()

	r, err := opencv.New(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "opencv.New: %v\n", err)
		fmt.Fprintln(os.Stderr, "This example expects a supported platform with an embedded native runtime.")
		os.Exit(1)
	}
	defer r.Close()

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", *outputDir, err)
		os.Exit(1)
	}

	fmt.Println("=== go-opencv examples ===")
	fmt.Printf("Library version: %s\n", opencv.Version)
	fmt.Printf("Input image:      %s\n", *input)
	fmt.Printf("Output directory: %s\n\n", *outputDir)

	env := demoEnv{runtime: r, input: *input, outputDir: *outputDir}

	switch *demo {
	case "basic":
		env.demoBasic()
	case "convert":
		env.demoConvert()
	case "resize":
		env.demoResize()
	case "edges":
		env.demoEdges()
	case "filter":
		env.demoFilter()
	case "contours":
		env.demoContours()
	case "warp":
		env.demoWarp()
	case "io":
		env.demoIO()
	default:
		env.demoBasic()
		fmt.Println()
		env.demoConvert()
		fmt.Println()
		env.demoResize()
		fmt.Println()
		env.demoEdges()
		fmt.Println()
		env.demoFilter()
		fmt.Println()
		env.demoContours()
		fmt.Println()
		env.demoWarp()
		fmt.Println()
		env.demoIO()
	}
}

type demoEnv struct {
	runtime   *opencv.Runtime
	input     string
	outputDir string
}

func (e demoEnv) mustRead() (*opencv.Mat, int, int, error) {
	img, err := e.runtime.IMRead(e.input)
	if err != nil {
		return nil, 0, 0, err
	}
	rows, _ := img.Rows()
	cols, _ := img.Cols()
	return img, rows, cols, nil
}

func (e demoEnv) out(name string) string {
	return filepath.Join(e.outputDir, name)
}

func (e demoEnv) demoBasic() {
	fmt.Println("--- 01. Basic: image info ---")

	img, err := e.runtime.IMRead(e.input)
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

	fmt.Printf("  File:       %s\n", e.input)
	fmt.Printf("  Size:       %d x %d pixels\n", cols, rows)
	fmt.Printf("  Type:       %d\n", typ)
	fmt.Printf("  Channels:   %d\n", ch)
	fmt.Printf("  ColorModel: %v\n", model)
	fmt.Println("  OK: Image loaded")
}

func (e demoEnv) demoConvert() {
	fmt.Println("--- 02. Color: BGR -> Gray ---")

	src, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	gray, err := e.runtime.NewMat(rows, cols, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer gray.Close()

	if err := e.runtime.CvtColor(src, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	outPath := e.out("convert-gray.png")
	if err := e.runtime.IMWrite(outPath, gray); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Color conversion")
}

func (e demoEnv) demoResize() {
	fmt.Println("--- 03. Resize: 50% scale ---")

	src, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	dstRows := max(1, rows/2)
	dstCols := max(1, cols/2)
	dst, err := e.runtime.NewMat(dstRows, dstCols, opencv.CV8UC3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer dst.Close()

	if err := e.runtime.Resize(src, dst, opencv.Size{Width: int32(dstCols), Height: int32(dstRows)}); err != nil {
		fmt.Fprintf(os.Stderr, "  Resize: %v\n", err)
		return
	}

	outPath := e.out("resize-half.png")
	if err := e.runtime.IMWrite(outPath, dst); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Before: %d x %d\n", cols, rows)
	fmt.Printf("  After:  %d x %d\n", dstCols, dstRows)
	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Resize")
}

func (e demoEnv) demoEdges() {
	fmt.Println("--- 04. Canny edge detection ---")

	src, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	gray, err := e.runtime.NewMat(rows, cols, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat (gray): %v\n", err)
		return
	}
	defer gray.Close()

	if err := e.runtime.CvtColor(src, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	edges, err := e.runtime.Canny(gray, 50, 150)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Canny: %v\n", err)
		return
	}
	defer edges.Close()

	outPath := e.out("edges-canny.png")
	if err := e.runtime.IMWrite(outPath, edges); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Edge detection")
}

func (e demoEnv) demoFilter() {
	fmt.Println("--- 05. Filter pipeline ---")

	src, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	blurred, err := e.runtime.GaussianBlur(src, opencv.Size{Width: 5, Height: 5}, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GaussianBlur: %v\n", err)
		return
	}
	defer blurred.Close()

	gray, err := e.runtime.NewMat(rows, cols, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat gray: %v\n", err)
		return
	}
	defer gray.Close()
	if err := e.runtime.CvtColor(blurred, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	sobel, err := e.runtime.Sobel(gray, opencv.CV8U, 1, 0, 3, 1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Sobel: %v\n", err)
		return
	}
	defer sobel.Close()

	kernel, err := e.runtime.GetStructuringElement(opencv.MorphRect, opencv.Size{Width: 3, Height: 3})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GetStructuringElement: %v\n", err)
		return
	}
	defer kernel.Close()

	closed, err := e.runtime.MorphologyEx(sobel, opencv.MorphClose, kernel, opencv.Point{X: -1, Y: -1}, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  MorphologyEx: %v\n", err)
		return
	}
	defer closed.Close()

	outPath := e.out("filter-pipeline.png")
	if err := e.runtime.IMWrite(outPath, closed); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Blur + Sobel + MorphClose")
}

func (e demoEnv) demoContours() {
	fmt.Println("--- 06. Contours ---")

	img, err := e.runtime.NewMat(200, 200, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer img.Close()

	if err := e.runtime.Rectangle(img, opencv.Rect{X: 20, Y: 20, Width: 60, Height: 60}, color.Gray{Y: 255}, -1); err != nil {
		fmt.Fprintf(os.Stderr, "  Rectangle: %v\n", err)
		return
	}
	if err := e.runtime.Circle(img, opencv.Point{X: 140, Y: 140}, 40, color.Gray{Y: 255}, -1); err != nil {
		fmt.Fprintf(os.Stderr, "  Circle: %v\n", err)
		return
	}

	contours, err := e.runtime.FindContours(img, opencv.RetrievalExternal, opencv.ChainApproxSimple)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  FindContours: %v\n", err)
		return
	}

	fmt.Printf("  Found %d contours\n", len(contours))
	for i, c := range contours {
		area, _ := e.runtime.ContourArea(c)
		rect, _ := e.runtime.BoundingRect(c)
		fmt.Printf("  Contour %d: %d points, area=%.0f, rect=%v\n", i, len(c), area, rect)
	}

	colorImg, err := e.runtime.NewMat(200, 200, opencv.CV8UC3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat color: %v\n", err)
		return
	}
	defer colorImg.Close()

	if err := e.runtime.DrawContours(colorImg, contours, -1, color.RGBA{0, 255, 0, 255}, 2); err != nil {
		fmt.Fprintf(os.Stderr, "  DrawContours: %v\n", err)
		return
	}

	outPath := e.out("contours-drawn.png")
	if err := e.runtime.IMWrite(outPath, colorImg); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Contours drawn")
}

func (e demoEnv) demoWarp() {
	fmt.Println("--- 07. Warp: rotate 45 degrees ---")

	src, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer src.Close()

	M, err := e.runtime.GetRotationMatrix2D(opencv.Point{X: int32(cols / 2), Y: int32(rows / 2)}, 45, 1.0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  GetRotationMatrix2D: %v\n", err)
		return
	}
	defer M.Close()

	dst, err := e.runtime.WarpAffine(src, M, opencv.Size{Width: int32(cols), Height: int32(rows)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  WarpAffine: %v\n", err)
		return
	}
	defer dst.Close()

	outPath := e.out("warp-rotate.png")
	if err := e.runtime.IMWrite(outPath, dst); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	fmt.Printf("  Output: %s\n", outPath)
	fmt.Println("  OK: Warp affine rotation")
}

func (e demoEnv) demoIO() {
	fmt.Println("--- 08. IO: read -> gray -> write ---")

	img, rows, cols, err := e.mustRead()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead: %v\n", err)
		return
	}
	defer img.Close()

	model, _ := img.ColorModel()
	fmt.Printf("  Read: %dx%d model=%v\n", cols, rows, model)

	gray, err := e.runtime.NewMat(rows, cols, opencv.CV8UC1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  NewMat: %v\n", err)
		return
	}
	defer gray.Close()

	if err := e.runtime.CvtColor(img, gray, opencv.ColorBGR2Gray); err != nil {
		fmt.Fprintf(os.Stderr, "  CvtColor: %v\n", err)
		return
	}

	outPath := e.out("output-gray.png")
	if err := e.runtime.IMWrite(outPath, gray); err != nil {
		fmt.Fprintf(os.Stderr, "  IMWrite: %v\n", err)
		return
	}

	back, err := e.runtime.IMRead(outPath, opencv.Gray)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  IMRead back: %v\n", err)
		return
	}
	defer back.Close()

	backModel, _ := back.ColorModel()
	backRows, _ := back.Rows()
	backCols, _ := back.Cols()
	fmt.Printf("  Wrote:     %s\n", outPath)
	fmt.Printf("  Read back: %dx%d model=%v\n", backCols, backRows, backModel)
	fmt.Println("  OK: IO roundtrip")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
