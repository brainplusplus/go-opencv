# go-opencv

[![Go Reference](https://pkg.go.dev/badge/github.com/brainplusplus/go-opencv.svg)](https://pkg.go.dev/github.com/brainplusplus/go-opencv)
[![Go Report Card](https://goreportcard.com/badge/github.com/brainplusplus/go-opencv)](https://goreportcard.com/report/github.com/brainplusplus/go-opencv)

CGO-free Go OpenCV library. Prebuilt native binary (`goopencv.dll/.so/.dylib`) wrapping real OpenCV via [opencv-mobile](https://github.com/nicehash/opencv-mobile) — loaded at runtime via [purego](https://github.com/ebitengine/purego).

**No C compiler required at build time.** Users just `go get` and run.

## Features

- **CGO_ENABLED=0** — zero C toolchain at user build time
- **OpenCV-style API** — `IMRead`, `CvtColor`, `GaussianBlur`, `FindContours`, `WarpAffine`, etc.
- **Prebuilt binary auto-load** — embedded `goopencv.dll` extracted to cache on first use
- **Explicit color model** — `BGR`, `RGB`, `RGBA`, `Gray` tracked per-Mat, no hidden conversions
- **96 ABI exports** — core + imgproc + morphology + contours + hough + warp + arithmetic

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    opencv "github.com/brainplusplus/go-opencv"
)

func main() {
    cv, err := opencv.New(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    defer cv.Close()

    // Read image (default: BGR)
    img, err := cv.IMRead("photo.png")
    if err != nil {
        log.Fatal(err)
    }
    defer img.Close()

    rows, _ := img.Rows()
    cols, _ := img.Cols()
    fmt.Printf("Loaded: %dx%d\n", cols, rows)

    // Convert BGR -> Gray
    gray, err := cv.NewMat(rows, cols, opencv.CV8UC1)
    if err != nil {
        log.Fatal(err)
    }
    defer gray.Close()
    cv.CvtColor(img, gray, opencv.ColorBGR2Gray)

    // Gaussian blur
    blurred, err := cv.GaussianBlur(img, opencv.Size{Width: 5, Height: 5}, 0)
    if err != nil {
        log.Fatal(err)
    }
    defer blurred.Close()

    // Write output
    cv.IMWrite("output.png", blurred)
}
```

## API Surface

### Image I/O
| Method | Description |
|--------|-------------|
| `IMRead(path, ...ColorModel)` | Read image file (PNG, JPEG, GIF) |
| `IMReadBytes(data, ...ColorModel)` | Read image from bytes |
| `IMWrite(path, mat, ...ColorModel)` | Write image to file |

### Color
| Method | Description |
|--------|-------------|
| `CvtColor(src, dst, code)` | Convert color space |
| `PutText(img, text, org, font, scale, color, thickness)` | Draw text |

### Filtering
| Method | Description |
|--------|-------------|
| `Blur(src, ksize)` | Box blur |
| `GaussianBlur(src, ksize, sigmaX)` | Gaussian blur |
| `MedianBlur(src, ksize)` | Median blur |
| `EqualizeHist(src)` | Histogram equalization |
| `Sobel(src, ddepth, dx, dy, ksize, scale, delta)` | Sobel derivative |
| `Laplacian(src, ddepth, ksize, scale, delta)` | Laplacian |
| `Threshold(src, thresh, maxval, typ)` | Binary threshold |
| `AdaptiveThreshold(src, maxval, adaptiveType, threshType, blockSize, c)` | Adaptive threshold |

### Morphology
| Method | Description |
|--------|-------------|
| `Erode(src, kernel, anchor, iterations)` | Erosion |
| `Dilate(src, kernel, anchor, iterations)` | Dilation |
| `MorphologyEx(src, op, kernel, anchor, iterations)` | Open/Close/Gradient/TopHat/BlackHat |
| `GetStructuringElement(shape, ksize)` | Create kernel (Rect/Cross/Ellipse) |

### Geometry
| Method | Description |
|--------|-------------|
| `Resize(src, dst, size)` | Resize image |
| `Flip(src, dst, flipCode)` | Flip (Horizontal/Vertical/Both) |
| `Transpose(src)` | Transpose |
| `WarpAffine(src, M, dsize)` | Affine warp |
| `WarpPerspective(src, M, dsize)` | Perspective warp |
| `GetRotationMatrix2D(center, angle, scale)` | 2D rotation matrix |
| `GetAffineTransform(src, dst)` | Affine transform matrix |

### Edge Detection
| Method | Description |
|--------|-------------|
| `Canny(src, threshold1, threshold2)` | Canny edge detector |

### Contours
| Method | Description |
|--------|-------------|
| `FindContours(src, mode, method)` | Find contours in binary image |
| `DrawContours(img, contours, idx, color, thickness)` | Draw contour outlines |
| `ContourArea(contour)` | Compute contour area |
| `ArcLength(contour, closed)` | Compute contour perimeter |
| `BoundingRect(contour)` | Compute bounding rectangle |
| `MinEnclosingCircle(contour)` | Compute minimum enclosing circle |
| `Moments(contour, binary)` | Compute spatial moments |

### Hough Transforms
| Method | Description |
|--------|-------------|
| `HoughLines(src, rho, theta, threshold)` | Standard Hough lines |
| `HoughLinesP(src, rho, theta, threshold, minLen, maxGap)` | Probabilistic Hough lines |
| `HoughCircles(src, method, dp, minDist, p1, p2, minR, maxR)` | Hough circles |

### Drawing
| Method | Description |
|--------|-------------|
| `Rectangle(img, rect, color, thickness)` | Draw rectangle |
| `Circle(img, center, radius, color, thickness)` | Draw circle |
| `Line(img, pt1, pt2, color, thickness)` | Draw line |
| `PutText(img, text, org, font, scale, color, thickness)` | Draw text |
| `ArrowedLine(img, pt1, pt2, color, thickness, tipLen)` | Draw arrow |

### Arithmetic & Logic
| Method | Description |
|--------|-------------|
| `Add(src1, src2)` | Per-element addition |
| `Subtract(src1, src2)` | Per-element subtraction |
| `Multiply(src1, src2, scale)` | Per-element multiplication |
| `Divide(src1, src2, scale)` | Per-element division |
| `AbsDiff(src1, src2)` | Per-element absolute difference |
| `BitwiseAnd/Or/Xor/Not` | Bitwise operations |

### Statistics
| Method | Description |
|--------|-------------|
| `MinMaxLoc(src)` | Find min/max values and locations |
| `MeanStdDev(src)` | Compute mean and standard deviation |
| `CountNonZero(src)` | Count non-zero elements |
| `Normalize(src, alpha, beta, normType)` | Normalize |

### Mat Operations
| Method | Description |
|--------|-------------|
| `NewMat(rows, cols, typ)` | Create empty Mat |
| `Zeros(rows, cols, typ)` | Create zero-initialized Mat |
| `Ones(rows, cols, typ)` | Create ones-initialized Mat |
| `Eye(rows, cols, typ)` | Create identity matrix |
| `Split(src)` | Split multi-channel Mat |
| `Merge(channels)` | Merge single-channel Mats |
| `Row/Col/Region/Reshape` | Sub-matrix extraction |

## Building the Native Backend

Windows (MSVC x64):

```shell
build-tools\build-goopencv.bat
```

This compiles `backend/goopencv_abi.cpp` against opencv-mobile 4.13.0 static libs and outputs `dist/goopencv.dll`.

## Running Tests

```shell
go test ./...
```

## Running Examples

```shell
cd examples
go run . -demo=all
```

## Color Model System

By default, `IMRead` returns BGR (OpenCV native). Specify a model to change:

```go
bgr, _  := cv.IMRead("img.png")           // BGR (default)
rgb, _  := cv.IMRead("img.png", cv.RGB)    // RGB
rgba, _ := cv.IMRead("img.png", cv.RGBA)   // RGBA
gray, _ := cv.IMRead("img.png", cv.Gray)   // Gray
```

Model metadata is tracked per-Mat. Processing ops preserve the source model. `CvtColor` determines the output model deterministically.

Strict validation is available but off by default:

```go
opencv.SetStrictColorValidation(true)
```

## Architecture

```
Go application
    |
    v
Runtime (runtime.go)           -- public API, Mat lifecycle, color metadata
    |
    v
PuregoBackend (purego_backend.go) -- resolve DLL symbols via ebitengine/purego
    |
    v
goopencv.dll (C++ ABI)          -- thin extern "C" wrappers calling cv:: functions
    |
    v
opencv-mobile 4.13.0            -- static-linked OpenCV core + imgproc
```

## License

MIT
