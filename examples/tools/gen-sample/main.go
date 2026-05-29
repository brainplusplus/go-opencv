// Command gen-sample generates a sample image for examples.
//
// Usage (from examples/):
//
//	go run ./tools/gen-sample/
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	const W, H = 400, 300
	img := image.NewRGBA(image.Rect(0, 0, W, H))

	// Background gradient (top-left to bottom-right)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r := uint8(x * 255 / W)
			g := uint8(y * 255 / H)
			b := uint8((x + y) * 255 / (W + H))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// White circle
	drawCircle(img, 320, 80, 50, color.RGBA{255, 255, 255, 255})

	// Red rectangle
	drawRect(img, 30, 40, 130, 140, color.RGBA{255, 0, 0, 255})

	// Green triangle (approximate with lines)
	drawTriangle(img, 200, 220, 260, 280, 140, 280, color.RGBA{0, 255, 0, 255})

	// Blue filled rectangle (semi-transparent feel via solid color)
	fillRect(img, 80, 180, 180, 260, color.RGBA{0, 100, 255, 255})

	f, err := os.Create("sample.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	println("Generated sample.png (400x300)")
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				img.Set(x, y, c)
			}
		}
	}
}

func drawRect(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == x1 || x == x2 || y == y1 || y == y2 {
				img.Set(x, y, c)
			}
		}
	}
}

func drawTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, c color.RGBA) {
	drawLine(img, x1, y1, x2, y2, c)
	drawLine(img, x2, y2, x3, y3, c)
	drawLine(img, x3, y3, x1, y1, c)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func fillRect(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	draw.Draw(img, image.Rect(x1, y1, x2, y2), &image.Uniform{c}, image.Point{}, draw.Src)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
