# API Reference — go-opencv

Complete mapping of all supported and unsupported functions, grouped by OpenCV module.

Legend: ✅ Supported | ❌ Stub (returns error) | 🔒 Not exposed (internal only)

---

## Matrix Operations (Core)

| Go API | Description | Status |
|--------|-------------|--------|
| `NewMat(rows, cols, typ)` | Create zero-initialized Mat | ✅ |
| `Mat.Delete()` | Free Mat memory | ✅ |
| `Mat.Clone()` | Deep copy Mat | ✅ |
| `Mat.Rows()` | Number of rows | ✅ |
| `Mat.Cols()` | Number of columns | ✅ |
| `Mat.Type()` | Mat data type (CV_8UC1, etc.) | ✅ |
| `Mat.Channels()` | Number of channels | ✅ |
| `Mat.ElemSize()` | Size of each element in bytes | ✅ |
| `Mat.Step()` | Number of bytes per row | ✅ |
| `Mat.Total()` | Total number of elements | ✅ |
| `Mat.Empty()` | Check if Mat is empty | ✅ |
| `Mat.Row(row)` | Extract single row | ✅ |
| `Mat.Col(col)` | Extract single column | ✅ |
| `Mat.Region(rect)` | Extract ROI sub-matrix | ✅ |
| `Mat.Reshape(ch, rows)` | Reshape without copying data | ✅ |
| `Mat.Diag()` | Extract diagonal | ✅ |
| `Mat.AtByte(row, col, ch)` | Get pixel value (uint8) | ✅ |
| `Mat.SetByte(row, col, ch, val)` | Set pixel value (uint8) | ✅ |
| `Mat.ColorModel()` | Get color model (BGR/Gray/etc.) | ✅ |
| `Mat.IsColorKnown()` | Check if color model is set | ✅ |
| `Mat.CopyTo(dst)` | Copy raw data to byte slice | ✅ |
| `Zeros(rows, cols, typ)` | Create zero Mat | ✅ |
| `Ones(rows, cols, typ)` | Create Mat filled with 1s | ✅ |
| `Eye(rows, cols, typ)` | Create identity Mat | ✅ |
| `Mat.SetTo(value)` | Set all elements to scalar | ✅ |
| `Mat.ConvertTo(typ)` | Convert Mat type | ✅ |

## Arithmetic & Logic (Core)

| Go API | Description | Status |
|--------|-------------|--------|
| `Add(src1, src2)` | Element-wise addition | ✅ |
| `Subtract(src1, src2)` | Element-wise subtraction | ✅ |
| `Multiply(src1, src2, scale)` | Element-wise multiplication | ✅ |
| `Divide(src1, src2, scale)` | Element-wise division | ✅ |
| `AbsDiff(src1, src2)` | Absolute difference | ✅ |
| `BitwiseAnd(src1, src2)` | Bitwise AND | ✅ |
| `BitwiseOr(src1, src2)` | Bitwise OR | ✅ |
| `BitwiseXor(src1, src2)` | Bitwise XOR | ✅ |
| `BitwiseNot(src)` | Bitwise NOT | ✅ |
| `Sqrt(src)` | Element-wise square root | ✅ |
| `Max(src1, src2)` | Element-wise maximum | ✅ |
| `Min(src1, src2)` | Element-wise minimum | ✅ |
| `Compare(src1, src2, cmpop)` | Element-wise comparison | ✅ |
| `CountNonZero(src)` | Count non-zero elements | ✅ |
| `Sum(src)` | Sum of all elements | ✅ |
| `Mean(src)` | Mean of all elements | ✅ |
| `Norm(src)` | Calculate norm | ✅ |
| `MinMaxLoc(src)` | Find min/max values and locations | ✅ |
| `MeanStdDev(src)` | Calculate mean and standard deviation | ✅ |
| `Split(src)` | Split multi-channel Mat into array | ✅ |
| `Merge(channels)` | Merge single-channel Mats into one | ✅ |

## Image I/O (Imgcodecs)

| Go API | Description | Status |
|--------|-------------|--------|
| `IMRead(path, model...)` | Read image file (default BGR) | ✅ |
| `IMReadBytes(data, model...)` | Read image from byte slice | ✅ |
| `IMWrite(path, mat, model...)` | Write image to file (PNG/JPG) | ✅ |

## Color Conversion (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `CvtColor(src, dst, code)` | Convert color space (BGR↔RGB↔Gray↔HSV↔etc.) | ✅ |
| `ConvertModel(src, model)` | Convert using color model enum | ✅ |

## Filtering (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `Blur(src, kSize)` | Simple box blur | ✅ |
| `GaussianBlur(src, kSize, sigmaX)` | Gaussian blur | ✅ |
| `MedianBlur(src, ksize)` | Median blur | ✅ |
| `BilateralFilter(src, d, sigmaColor, sigmaSpace)` | Bilateral filter (edge-preserving) | ✅ |
| `Sobel(src, ddepth, dx, dy, ksize, scale, delta)` | Sobel edge detection | ✅ |
| `Laplacian(src, ddepth, ksize, scale, delta)` | Laplacian edge detection | ✅ |
| `EqualizeHist(src)` | Histogram equalization (grayscale) | ✅ |
| `Erode(src, kernel, anchor, iterations)` | Morphological erosion | ✅ |
| `Dilate(src, kernel, anchor, iterations)` | Morphological dilation | ✅ |
| `MorphologyEx(src, op, kernel, anchor, iterations)` | Advanced morphology (open/close/gradient/tophat/blackhat) | ✅ |
| `GetStructuringElement(shape, ksize)` | Create structuring element for morphology | ✅ |
| `Threshold(src, thresh, maxval, typ)` | Fixed-level thresholding | ✅ |
| `AdaptiveThreshold(src, maxval, adaptiveType, thresholdType, blockSize, c)` | Adaptive thresholding | ✅ |
| `Canny(src, thresh1, thresh2)` | Canny edge detection | ✅ |
| `LUT(src, lut)` | Apply look-up table | ✅ |
| `Normalize(src, alpha, beta, normType)` | Normalize Mat | ✅ |

## Geometric Transforms (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `Resize(src, dst, size)` | Resize image | ✅ |
| `Flip(src, dst, flipCode)` | Flip image (horizontal/vertical/both) | ✅ |
| `Transpose(src)` | Transpose image | ✅ |
| `Rotate(src, code)` | Rotate 90°/180° | ✅ |
| `WarpAffine(src, M, dsize)` | Affine transformation | ✅ |
| `WarpPerspective(src, M, dsize)` | Perspective transformation | ✅ |
| `GetRotationMatrix2D(center, angle, scale)` | 2x3 rotation matrix | ✅ |
| `GetAffineTransform(src, dst)` | 2x3 affine transform matrix from 3 point pairs | ✅ |
| `GetPerspectiveTransform(src, dst)` | 3x3 perspective matrix from 4 point pairs | ✅ |
| `Remap(src, map1, map2, interp, borderMode, borderVal)` | Generic geometric remap | ✅ |
| `CopyMakeBorder(src, top, bottom, left, right, borderType, value)` | Pad image borders | ✅ |
| `Hconcat(src1, src2)` | Horizontal concatenation | ✅ |
| `Vconcat(src1, src2)` | Vertical concatenation | ✅ |
| `Integral(src)` | Calculate integral image | ✅ |

## Drawing (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `PutText(img, text, org, fontFace, fontScale, color, thickness)` | Draw text | ✅ |
| `Rectangle(img, rect, color, thickness)` | Draw rectangle | ✅ |
| `Circle(img, center, radius, color, thickness)` | Draw circle | ✅ |
| `Line(img, pt1, pt2, color, thickness)` | Draw line | ✅ |
| `ArrowedLine(img, pt1, pt2, color, thickness, tipLength)` | Draw arrow | ✅ |
| `FillPoly(img, contours, color, lineType)` | Draw filled polygon | ✅ |
| `FillConvexPoly(img, points, color, lineType)` | Draw filled convex polygon | ✅ |

## Contours (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `FindContours(src, mode, method)` | Find contours in binary image | ✅ |
| `DrawContours(img, contours, idx, color, thickness)` | Draw contours | ✅ |
| `ContourArea(contour)` | Calculate contour area | ✅ |
| `ArcLength(contour, closed)` | Calculate contour perimeter | ✅ |
| `BoundingRect(contour)` | Compute bounding rectangle | ✅ |
| `MinEnclosingCircle(contour)` | Compute minimum enclosing circle | ✅ |
| `Moments(contour, binary)` | Calculate image moments | ✅ |
| `ConnectedComponents(src, connectivity)` | Label connected components | ✅ |

## Hough Transforms (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `HoughLines(src, rho, theta, threshold)` | Standard Hough line detection | ✅ |
| `HoughLinesP(src, rho, theta, threshold, minLen, maxGap)` | Probabilistic Hough lines | ✅ |
| `HoughCircles(src, method, dp, minDist, param1, param2, minR, maxR)` | Hough circle detection | ✅ |

## Histogram (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `CalcHist(src, bins, rangeMin, rangeMax)` | Calculate histogram | ✅ |
| `InRange(src, lower, upper)` | Threshold within range | ✅ |
| `MatchTemplate(img, tmpl, method)` | Template matching | ✅ |

## Distance Transform (Imgproc)

| Go API | Description | Status |
|--------|-------------|--------|
| `DistanceTransform(src, distType, maskSize)` | Calculate distance to zero pixels | ✅ |

## Photo (Computational Photography)

| Go API | Description | Status |
|--------|-------------|--------|
| `FastNlMeansDenoising(src, h, templateWin, searchWin)` | Denoise grayscale image | ✅ |
| `FastNlMeansDenoisingColored(src, h, hColor, templateWin, searchWin)` | Denoise color image | ✅ |
| `DetailEnhance(src, sigmaS, sigmaR)` | Edge-preserving detail enhance | ✅ |
| `EdgePreservingFilter(src, flags, sigmaS, sigmaR)` | Edge-preserving smoothing | ✅ |
| `PencilSketch(src, sigmaS, sigmaR, shadeFactor)` | Pencil sketch effect | ✅ |
| `Stylization(src, sigmaS, sigmaR)` | Stylization (oil painting) effect | ✅ |
| `SeamlessClone(src, dst, mask, center, flags)` | Seamless image cloning | ✅ |

## Feature Detection (Features2d)

| Go API | Description | Status |
|--------|-------------|--------|
| `FAST(src, threshold, nonmaxSuppression)` | FAST corner detection | ✅ |
| `ORBDetectCompute(src, nfeatures, scaleFactor, nlevels)` | ORB feature detection + descriptor | ✅ |
| `BFMatch(desc1, desc2, normType)` | Brute-force descriptor matching | ✅ |
| `DrawKeypoints(img, keypoints, color)` | Draw feature keypoints | ✅ |

## HighGUI (Window Management)

| Go API | Description | Status |
|--------|-------------|--------|
| `ImShow(winname, img)` | Display image in window | ❌ Stub |
| `WaitKey(delay)` | Wait for key press | ❌ Stub |
| `DestroyWindow(winname)` | Close window | ❌ Stub |

> HighGUI functions return `ERR_UNSUPPORTED` — opencv-mobile does not include windowing subsystem.

## Internal Vector Helpers

These are used internally for transferring arrays across the FFI boundary. Not exposed in public API.

| ABI Function | Purpose |
|-------------|---------|
| `vec_points_*` | Point arrays (contours) |
| `vec_vec_points_*` | Array of point arrays (multiple contours) |
| `vec_double_*` | Double arrays (histogram bins) |
| `vec_int_*` | Int arrays (labels) |
| `vec_mat_*` | Mat arrays (channels) |
| `vec_keypoint_*` | KeyPoint arrays (features) |
| `vec_dmatch_*` | DMatch arrays (matching results) |

---

## Summary

| Category | Total | Supported | Stub |
|----------|-------|-----------|------|
| Matrix Ops (Core) | 28 | 28 | 0 |
| Arithmetic & Logic | 20 | 20 | 0 |
| Image I/O | 3 | 3 | 0 |
| Color | 2 | 2 | 0 |
| Filtering | 15 | 15 | 0 |
| Geometric Transforms | 14 | 14 | 0 |
| Drawing | 7 | 7 | 0 |
| Contours | 8 | 8 | 0 |
| Hough | 3 | 3 | 0 |
| Histogram | 3 | 3 | 0 |
| Distance Transform | 1 | 1 | 0 |
| Photo | 7 | 7 | 0 |
| Feature Detection | 4 | 4 | 0 |
| HighGUI | 3 | 0 | **3** |
| **Total** | **118** | **115** | **3** |

All 115 functional APIs are fully wired and tested. Only HighGUI windowing (ImShow, WaitKey, DestroyWindow) is stubbed — opencv-mobile doesn't include the windowing subsystem.

## Platform Support

| Platform | GOARCH | Supported |
|----------|--------|-----------|
| Windows | arm64 | ✅ |
| Windows | amd64 | ✅ |
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| macOS | amd64 (Intel) | ✅ |
| macOS | arm64 (Apple Silicon) | ✅ |
