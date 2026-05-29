package runtime

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/brainplusplus/go-opencv/internal/abi"
	"github.com/ebitengine/purego"
)

// PuregoBackend loads goopencv.dll/.so/.dylib and calls OpenCV functions directly.
type PuregoBackend struct {
	lib uintptr

	matNew         func(rows, cols, typ int32) uint64
	matNewFromData func(data unsafe.Pointer, rows, cols, typ int32) uint64
	matDelete      func(handle uint64) int32
	matRows        func(handle uint64) int32
	matCols        func(handle uint64) int32
	matType        func(handle uint64) int32
	matClone       func(handle uint64) uint64
	matTotal       func(handle uint64) int32
	matRow         func(handle uint64, row int32) uint64
	matCol         func(handle uint64, col int32) uint64
	matRegion      func(handle uint64, x, y, w, h int32) uint64
	matReshape     func(handle uint64, cn, rows int32) uint64
	matSetTo       func(handle uint64, v0, v1, v2, v3 float64) int32
	matConvertTo   func(src, dst uint64, rtype int32, alpha, beta float64) int32
	matZeros       func(rows, cols, typ int32) uint64
	matOnes        func(rows, cols, typ int32) uint64
	matEye         func(rows, cols, typ int32) uint64

	cvtColor func(src, dst uint64, code int32) int32
	resize   func(src, dst uint64, width, height int32) int32

	// Filtering
	blur           func(src, dst uint64, kw, kh int32) int32
	gaussianBlur   func(src, dst uint64, kw, kh int32, sigmaX float64) int32
	medianBlur     func(src, dst uint64, ksize int32) int32
	threshold      func(src, dst uint64, thresh, maxval float64, typ int32) int32
	adaptiveThresh func(src, dst uint64, maxval float64, adaptiveMethod, thresholdType int32, blockSize int32, c float64) int32
	canny          func(src, dst uint64, t1, t2 float64) int32
	flip           func(src, dst uint64, flipCode int32) int32
	sobel          func(src, dst uint64, ddepth, dx, dy, ksize int32, scale, delta float64) int32
	laplacian      func(src, dst uint64, ddepth, ksize int32, scale, delta float64) int32
	transpose      func(src, dst uint64) int32
	equalizeHist   func(src, dst uint64) int32
	normalize      func(src, dst uint64, alpha, beta float64, normType int32) int32

	// Morphology
	erode     func(src, dst, kernel uint64, anchorX, anchorY, iterations int32) int32
	dilate    func(src, dst, kernel uint64, anchorX, anchorY, iterations int32) int32
	morphEx   func(src, dst uint64, op int32, kernel uint64, anchorX, anchorY, iterations int32) int32
	getKernel func(shape, kw, kh int32) uint64

	// Drawing
	rectangle   func(img uint64, x1, y1, x2, y2 int32, r, g, b, a uint8, thickness int32) int32
	circle      func(img uint64, cx, cy, radius int32, r, g, b, a uint8, thickness int32) int32
	line        func(img uint64, x1, y1, x2, y2 int32, r, g, b, a uint8, thickness int32) int32
	putText     func(img uint64, text unsafe.Pointer, textLen int32, orgX, orgY, fontFace int32, fontScale float64, r, g, b, a uint8, thickness, lineType, bottomLeftOrigin int32) int32
	fillPoly    func(img uint64, pts unsafe.Pointer, npts, ncontours int32, r, g, b, a uint8, lineType, shift int32) int32
	arrowedLine func(img uint64, x1, y1, x2, y2 int32, r, g, b, a uint8, thickness, lineType, shift int32, tipLength float64) int32

	// Contours
	findContours       func(src, contoursVec uint64, mode, method, offX, offY int32) int32
	drawContours       func(img, contoursVec uint64, contourIdx int32, r, g, b, a uint8, thickness int32) int32
	contourArea        func(vecHandle uint64) float64
	arcLength          func(vecHandle uint64, closed int32) float64
	boundingRect       func(vecHandle uint64, outX, outY, outW, outH *int32) int32
	minEnclosingCircle func(vecHandle uint64, outCX, outCY, outRadius *float64) int32
	moments            func(vecHandle uint64, binary int32, outM0, outM1, outM2, outM3, outM4, outM5, outM6, outM7, outM8, outM9 *float64) int32

	// Hough
	houghLines   func(src, vecHandle uint64, rho, theta float64, threshold int32, srn, stn float64) int32
	houghLinesP  func(src, vecHandle uint64, rho, theta float64, threshold int32, minLen, maxGap float64) int32
	houghCircles func(src, vecHandle uint64, method int32, dp, minDist, p1, p2 float64, minR, maxR int32) int32

	// Warp
	warpAffine         func(src, dst, m uint64, dstW, dstH int32) int32
	warpPerspective    func(src, dst, m uint64, dstW, dstH int32) int32
	getRotationMat2D   func(cx, cy, angle, scale float64) uint64
	getAffineTransform func(src0x, src0y, src1x, src1y, src2x, src2y, dst0x, dst0y, dst1x, dst1y, dst2x, dst2y float64) uint64

	// Core ops
	bitwiseAnd   func(src1, src2, dst uint64) int32
	bitwiseOr    func(src1, src2, dst uint64) int32
	bitwiseXor   func(src1, src2, dst uint64) int32
	bitwiseNot   func(src, dst uint64) int32
	add          func(src1, src2, dst uint64) int32
	subtract     func(src1, src2, dst uint64) int32
	multiply     func(src1, src2, dst uint64, scale float64, dtype int32) int32
	divide       func(src1, src2, dst uint64, scale float64, dtype int32) int32
	absDiff      func(src1, src2, dst uint64) int32
	minMaxLoc    func(src uint64, outMin, outMax *float64, outMinX, outMinY, outMaxX, outMaxY *int32) int32
	meanStdDev   func(src, mean, stddev uint64) int32
	countNonZero func(src uint64) int32
	split        func(src, vecHandle uint64) int32
	merge        func(vecHandle, dst uint64) int32

	// Mat access
	matEmpty    func(handle uint64) int32
	matElemSz   func(handle uint64) int32
	matDataPtr  func(handle uint64) uint64
	matStep     func(handle uint64) int32
	matChannels func(handle uint64) int32

	// Vector helpers — points
	vecPointsNew    func() uint64
	vecPointsPush   func(handle uint64, x, y int32)
	vecPointsLen    func(handle uint64) int32
	vecPointsGet    func(handle uint64, idx int32, outX, outY *int32) int32
	vecPointsDelete func(handle uint64)

	// Vector helpers — vector of points (contours)
	vecVecPointsNew    func() uint64
	vecVecPointsPush   func(handle, contourHandle uint64)
	vecVecPointsLen    func(handle uint64) int32
	vecVecPointsGet    func(handle uint64, idx int32) uint64
	vecVecPointsDelete func(handle uint64)

	// Vector helpers — double
	vecDoubleNew    func() uint64
	vecDoubleGet    func(handle uint64, idx int32) float64
	vecDoubleLen    func(handle uint64) int32
	vecDoubleDelete func(handle uint64)

	// Vector helpers — int
	vecIntNew    func() uint64
	vecIntGet    func(handle uint64, idx int32) int32
	vecIntLen    func(handle uint64) int32
	vecIntDelete func(handle uint64)

	// Vector helpers — mat
	vecMatNew    func() uint64
	vecMatPush   func(handle, matHandle uint64)
	vecMatLen    func(handle uint64) int32
	vecMatGet    func(handle uint64, idx int32) uint64
	vecMatDelete func(handle uint64)
}

// NewPuregoBackend loads the native shared library and resolves all ABI exports.
func NewPuregoBackend(libPath string) (*PuregoBackend, error) {
	lib, err := openLibrary(libPath)
	if err != nil {
		return nil, fmt.Errorf("purego: open %q: %w", libPath, err)
	}

	b := &PuregoBackend{lib: lib}

	// Resolve required functions — panic on missing symbol
	reqs := []struct {
		name string
		fn   interface{}
	}{
		{abi.MatNew, &b.matNew},
		{abi.MatNewFromData, &b.matNewFromData},
		{abi.MatDelete, &b.matDelete},
		{abi.MatRows, &b.matRows},
		{abi.MatCols, &b.matCols},
		{abi.MatType, &b.matType},
		{abi.MatClone, &b.matClone},
		{abi.CvtColor, &b.cvtColor},
		{abi.Resize, &b.resize},
	}

	var missing []string
	for _, r := range reqs {
		func() {
			defer func() {
				if v := recover(); v != nil {
					missing = append(missing, r.name)
				}
			}()
			purego.RegisterLibFunc(r.fn, lib, r.name)
		}()
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("purego: missing ABI exports: %v", missing)
	}

	// Resolve optional functions (best-effort, no error if missing)
	opts := []struct {
		name string
		fn   interface{}
	}{
		// Mat extras
		{abi.MatEmpty, &b.matEmpty},
		{abi.MatElemSize, &b.matElemSz},
		{abi.MatDataPtr, &b.matDataPtr},
		{abi.MatStep, &b.matStep},
		{abi.MatChannels, &b.matChannels},
		{abi.MatTotal, &b.matTotal},
		{abi.MatRow, &b.matRow},
		{abi.MatCol, &b.matCol},
		{abi.MatRegion, &b.matRegion},
		{abi.MatReshape, &b.matReshape},
		{abi.MatSetTo, &b.matSetTo},
		{abi.MatConvertTo, &b.matConvertTo},
		{abi.MatZeros, &b.matZeros},
		{abi.MatOnes, &b.matOnes},
		{abi.MatEye, &b.matEye},

		// Filtering
		{abi.Blur, &b.blur},
		{abi.GaussianBlur, &b.gaussianBlur},
		{abi.MedianBlur, &b.medianBlur},
		{abi.Threshold, &b.threshold},
		{abi.AdaptiveThreshold, &b.adaptiveThresh},
		{abi.Canny, &b.canny},
		{abi.Flip, &b.flip},
		{abi.Sobel, &b.sobel},
		{abi.Laplacian, &b.laplacian},
		{abi.Transpose, &b.transpose},
		{abi.EqualizeHist, &b.equalizeHist},
		{abi.Normalize, &b.normalize},

		// Morphology
		{abi.Erode, &b.erode},
		{abi.Dilate, &b.dilate},
		{abi.MorphologyEx, &b.morphEx},
		{abi.GetStructuringElement, &b.getKernel},

		// Drawing
		{abi.Rectangle, &b.rectangle},
		{abi.Circle, &b.circle},
		{abi.Line, &b.line},
		{abi.PutText, &b.putText},
		{abi.FillPoly, &b.fillPoly},
		{abi.ArrowedLine, &b.arrowedLine},

		// Contours
		{abi.FindContours, &b.findContours},
		{abi.DrawContours, &b.drawContours},
		{abi.ContourArea, &b.contourArea},
		{abi.ArcLength, &b.arcLength},
		{abi.BoundingRect, &b.boundingRect},
		{abi.MinEnclosingCircle, &b.minEnclosingCircle},
		{abi.Moments, &b.moments},

		// Hough
		{abi.HoughLines, &b.houghLines},
		{abi.HoughLinesP, &b.houghLinesP},
		{abi.HoughCircles, &b.houghCircles},

		// Warp
		{abi.WarpAffine, &b.warpAffine},
		{abi.WarpPerspective, &b.warpPerspective},
		{abi.GetRotationMatrix2D, &b.getRotationMat2D},
		{abi.GetAffineTransform, &b.getAffineTransform},

		// Core ops
		{abi.BitwiseAnd, &b.bitwiseAnd},
		{abi.BitwiseOr, &b.bitwiseOr},
		{abi.BitwiseXor, &b.bitwiseXor},
		{abi.BitwiseNot, &b.bitwiseNot},
		{abi.Add, &b.add},
		{abi.Subtract, &b.subtract},
		{abi.Multiply, &b.multiply},
		{abi.Divide, &b.divide},
		{abi.AbsDiff, &b.absDiff},
		{abi.MinMaxLoc, &b.minMaxLoc},
		{abi.MeanStdDev, &b.meanStdDev},
		{abi.CountNonZero, &b.countNonZero},
		{abi.Split, &b.split},
		{abi.Merge, &b.merge},

		// Vector helpers — points
		{abi.VecNewPoints, &b.vecPointsNew},
		{abi.VecPushPoint, &b.vecPointsPush},
		{abi.VecLenPoints, &b.vecPointsLen},
		{abi.VecGetPoint, &b.vecPointsGet},
		{abi.VecDeletePoints, &b.vecPointsDelete},

		// Vector helpers — vector of points
		{abi.VecNewVecPoints, &b.vecVecPointsNew},
		{abi.VecPushVecPoints, &b.vecVecPointsPush},
		{abi.VecLenVecPoints, &b.vecVecPointsLen},
		{abi.VecGetVecPoints, &b.vecVecPointsGet},
		{abi.VecDeleteVecPoints, &b.vecVecPointsDelete},

		// Vector helpers — double
		{abi.VecNewDouble, &b.vecDoubleNew},
		{abi.VecGetDouble, &b.vecDoubleGet},
		{abi.VecLenDouble, &b.vecDoubleLen},
		{abi.VecDeleteDouble, &b.vecDoubleDelete},

		// Vector helpers — int
		{abi.VecNewInt, &b.vecIntNew},
		{abi.VecGetInt, &b.vecIntGet},
		{abi.VecLenInt, &b.vecIntLen},
		{abi.VecDeleteInt, &b.vecIntDelete},

		// Vector helpers — mat
		{abi.VecNewMat, &b.vecMatNew},
		{abi.VecPushMat, &b.vecMatPush},
		{abi.VecLenMat, &b.vecMatLen},
		{abi.VecGetMat, &b.vecMatGet},
		{abi.VecDeleteMat, &b.vecMatDelete},
	}
	for _, o := range opts {
		func() {
			defer func() { recover() }()
			purego.RegisterLibFunc(o.fn, lib, o.name)
		}()
	}

	return b, nil
}

func (b *PuregoBackend) Close(_ context.Context) error {
	if b.lib != 0 {
		closeLibrary(b.lib)
	}
	b.lib = 0
	return nil
}

// ---------------------------------------------------------------------------
// Mat lifecycle
// ---------------------------------------------------------------------------

func (b *PuregoBackend) NewMat(_ context.Context, rows, cols int, typ int32) (uint64, error) {
	h := b.matNew(int32(rows), int32(cols), typ)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_new returned 0 (rows=%d cols=%d type=%d)", rows, cols, typ)
	}
	return h, nil
}

func (b *PuregoBackend) NewMatFromData(_ context.Context, data []byte, rows, cols int, typ int32) (uint64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("purego: mat_new_from_data: empty data")
	}
	h := b.matNewFromData(unsafe.Pointer(&data[0]), int32(rows), int32(cols), typ)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_new_from_data failed (rows=%d cols=%d type=%d)", rows, cols, typ)
	}
	return h, nil
}

func (b *PuregoBackend) CloseMat(_ context.Context, handle uint64) error {
	code := b.matDelete(handle)
	return errorCode(abi.MatDelete, abi.ErrorCode(code))
}

func (b *PuregoBackend) MatRows(_ context.Context, handle uint64) (int, error) {
	return int(b.matRows(handle)), nil
}

func (b *PuregoBackend) MatCols(_ context.Context, handle uint64) (int, error) {
	return int(b.matCols(handle)), nil
}

func (b *PuregoBackend) MatType(_ context.Context, handle uint64) (int32, error) {
	return b.matType(handle), nil
}

func (b *PuregoBackend) MatClone(_ context.Context, handle uint64) (uint64, error) {
	h := b.matClone(handle)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_clone failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatTotal(_ context.Context, handle uint64) (int, error) {
	if b.matTotal == nil {
		return 0, errNotSupported("mat_total")
	}
	return int(b.matTotal(handle)), nil
}

func (b *PuregoBackend) MatRow(_ context.Context, handle uint64, row int32) (uint64, error) {
	if b.matRow == nil {
		return 0, errNotSupported("mat_row")
	}
	h := b.matRow(handle, row)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_row failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatCol(_ context.Context, handle uint64, col int32) (uint64, error) {
	if b.matCol == nil {
		return 0, errNotSupported("mat_col")
	}
	h := b.matCol(handle, col)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_col failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatRegion(_ context.Context, handle uint64, x, y, w, h int32) (uint64, error) {
	if b.matRegion == nil {
		return 0, errNotSupported("mat_region")
	}
	rh := b.matRegion(handle, x, y, w, h)
	if rh == 0 {
		return 0, fmt.Errorf("purego: mat_region failed")
	}
	return rh, nil
}

func (b *PuregoBackend) MatReshape(_ context.Context, handle uint64, cn, rows int32) (uint64, error) {
	if b.matReshape == nil {
		return 0, errNotSupported("mat_reshape")
	}
	h := b.matReshape(handle, cn, rows)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_reshape failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatSetTo(_ context.Context, handle uint64, v0, v1, v2, v3 float64) error {
	if b.matSetTo == nil {
		return errNotSupported("mat_set_to")
	}
	return errorCode("mat_set_to", abi.ErrorCode(b.matSetTo(handle, v0, v1, v2, v3)))
}

func (b *PuregoBackend) MatConvertTo(_ context.Context, src, dst uint64, rtype int32, alpha, beta float64) error {
	if b.matConvertTo == nil {
		return errNotSupported("mat_convert_to")
	}
	return errorCode("mat_convert_to", abi.ErrorCode(b.matConvertTo(src, dst, rtype, alpha, beta)))
}

func (b *PuregoBackend) MatZeros(_ context.Context, rows, cols int, typ int32) (uint64, error) {
	if b.matZeros == nil {
		return 0, errNotSupported("mat_zeros")
	}
	h := b.matZeros(int32(rows), int32(cols), typ)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_zeros failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatOnes(_ context.Context, rows, cols int, typ int32) (uint64, error) {
	if b.matOnes == nil {
		return 0, errNotSupported("mat_ones")
	}
	h := b.matOnes(int32(rows), int32(cols), typ)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_ones failed")
	}
	return h, nil
}

func (b *PuregoBackend) MatEye(_ context.Context, rows, cols int, typ int32) (uint64, error) {
	if b.matEye == nil {
		return 0, errNotSupported("mat_eye")
	}
	h := b.matEye(int32(rows), int32(cols), typ)
	if h == 0 {
		return 0, fmt.Errorf("purego: mat_eye failed")
	}
	return h, nil
}

// ---------------------------------------------------------------------------
// Mat access
// ---------------------------------------------------------------------------

func (b *PuregoBackend) MatEmpty(_ context.Context, handle uint64) (bool, error) {
	if b.matEmpty == nil {
		return false, errNotSupported("mat_empty")
	}
	return b.matEmpty(handle) != 0, nil
}

func (b *PuregoBackend) MatElemSize(_ context.Context, handle uint64) (int, error) {
	if b.matElemSz == nil {
		return 0, errNotSupported("mat_elem_size")
	}
	return int(b.matElemSz(handle)), nil
}

func (b *PuregoBackend) MatDataPtr(_ context.Context, handle uint64) (unsafe.Pointer, error) {
	if b.matDataPtr == nil {
		return nil, errNotSupported("mat_data_ptr")
	}
	ptr := b.matDataPtr(handle)
	if ptr == 0 {
		return nil, fmt.Errorf("mat_data_ptr: null pointer")
	}
	return unsafe.Pointer(uintptr(ptr)), nil
}

func (b *PuregoBackend) MatStep(_ context.Context, handle uint64) (int, error) {
	if b.matStep == nil {
		return 0, errNotSupported("mat_step")
	}
	return int(b.matStep(handle)), nil
}

func (b *PuregoBackend) MatChannels(_ context.Context, handle uint64) (int, error) {
	if b.matChannels == nil {
		return 0, errNotSupported("mat_channels")
	}
	return int(b.matChannels(handle)), nil
}

// ---------------------------------------------------------------------------
// Image processing — filtering
// ---------------------------------------------------------------------------

func (b *PuregoBackend) CvtColor(_ context.Context, src, dst uint64, code int32) error {
	ec := b.cvtColor(src, dst, code)
	return errorCode(abi.CvtColor, abi.ErrorCode(ec))
}

func (b *PuregoBackend) Resize(_ context.Context, src, dst uint64, width, height int32) error {
	ec := b.resize(src, dst, width, height)
	return errorCode(abi.Resize, abi.ErrorCode(ec))
}

func (b *PuregoBackend) Blur(_ context.Context, src, dst uint64, kw, kh int32) error {
	if b.blur == nil {
		return errNotSupported("blur")
	}
	return errorCode(abi.Blur, abi.ErrorCode(b.blur(src, dst, kw, kh)))
}

func (b *PuregoBackend) GaussianBlur(_ context.Context, src, dst uint64, kw, kh int32, sigmaX float64) error {
	if b.gaussianBlur == nil {
		return errNotSupported("gaussian_blur")
	}
	return errorCode(abi.GaussianBlur, abi.ErrorCode(b.gaussianBlur(src, dst, kw, kh, sigmaX)))
}

func (b *PuregoBackend) MedianBlur(_ context.Context, src, dst uint64, ksize int32) error {
	if b.medianBlur == nil {
		return errNotSupported("median_blur")
	}
	return errorCode(abi.MedianBlur, abi.ErrorCode(b.medianBlur(src, dst, ksize)))
}

func (b *PuregoBackend) Threshold(_ context.Context, src, dst uint64, thresh, maxval float64, typ int32) error {
	if b.threshold == nil {
		return errNotSupported("threshold")
	}
	return errorCode(abi.Threshold, abi.ErrorCode(b.threshold(src, dst, thresh, maxval, typ)))
}

func (b *PuregoBackend) AdaptiveThreshold(_ context.Context, src, dst uint64, maxval float64, adaptiveMethod, thresholdType, blockSize int32, c float64) error {
	if b.adaptiveThresh == nil {
		return errNotSupported("adaptive_threshold")
	}
	return errorCode(abi.AdaptiveThreshold, abi.ErrorCode(b.adaptiveThresh(src, dst, maxval, adaptiveMethod, thresholdType, blockSize, c)))
}

func (b *PuregoBackend) Canny(_ context.Context, src, dst uint64, t1, t2 float64) error {
	if b.canny == nil {
		return errNotSupported("canny")
	}
	return errorCode(abi.Canny, abi.ErrorCode(b.canny(src, dst, t1, t2)))
}

func (b *PuregoBackend) Flip(_ context.Context, src, dst uint64, flipCode int32) error {
	if b.flip == nil {
		return errNotSupported("flip")
	}
	return errorCode(abi.Flip, abi.ErrorCode(b.flip(src, dst, flipCode)))
}

func (b *PuregoBackend) Sobel(_ context.Context, src, dst uint64, ddepth, dx, dy, ksize int32, scale, delta float64) error {
	if b.sobel == nil {
		return errNotSupported("sobel")
	}
	return errorCode(abi.Sobel, abi.ErrorCode(b.sobel(src, dst, ddepth, dx, dy, ksize, scale, delta)))
}

func (b *PuregoBackend) Laplacian(_ context.Context, src, dst uint64, ddepth, ksize int32, scale, delta float64) error {
	if b.laplacian == nil {
		return errNotSupported("laplacian")
	}
	return errorCode(abi.Laplacian, abi.ErrorCode(b.laplacian(src, dst, ddepth, ksize, scale, delta)))
}

func (b *PuregoBackend) Transpose(_ context.Context, src, dst uint64) error {
	if b.transpose == nil {
		return errNotSupported("transpose")
	}
	return errorCode(abi.Transpose, abi.ErrorCode(b.transpose(src, dst)))
}

func (b *PuregoBackend) EqualizeHist(_ context.Context, src, dst uint64) error {
	if b.equalizeHist == nil {
		return errNotSupported("equalize_hist")
	}
	return errorCode(abi.EqualizeHist, abi.ErrorCode(b.equalizeHist(src, dst)))
}

func (b *PuregoBackend) Normalize(_ context.Context, src, dst uint64, alpha, beta float64, normType int32) error {
	if b.normalize == nil {
		return errNotSupported("normalize")
	}
	return errorCode(abi.Normalize, abi.ErrorCode(b.normalize(src, dst, alpha, beta, normType)))
}

// ---------------------------------------------------------------------------
// Morphology
// ---------------------------------------------------------------------------

func (b *PuregoBackend) Erode(_ context.Context, src, dst, kernel uint64, anchorX, anchorY, iterations int32) error {
	if b.erode == nil {
		return errNotSupported("erode")
	}
	return errorCode(abi.Erode, abi.ErrorCode(b.erode(src, dst, kernel, anchorX, anchorY, iterations)))
}

func (b *PuregoBackend) Dilate(_ context.Context, src, dst, kernel uint64, anchorX, anchorY, iterations int32) error {
	if b.dilate == nil {
		return errNotSupported("dilate")
	}
	return errorCode(abi.Dilate, abi.ErrorCode(b.dilate(src, dst, kernel, anchorX, anchorY, iterations)))
}

func (b *PuregoBackend) MorphologyEx(_ context.Context, src, dst uint64, op int32, kernel uint64, anchorX, anchorY, iterations int32) error {
	if b.morphEx == nil {
		return errNotSupported("morphology_ex")
	}
	return errorCode(abi.MorphologyEx, abi.ErrorCode(b.morphEx(src, dst, op, kernel, anchorX, anchorY, iterations)))
}

func (b *PuregoBackend) GetStructuringElement(_ context.Context, shape, kw, kh int32) (uint64, error) {
	if b.getKernel == nil {
		return 0, errNotSupported("get_structuring_element")
	}
	h := b.getKernel(shape, kw, kh)
	if h == 0 {
		return 0, fmt.Errorf("purego: get_structuring_element failed")
	}
	return h, nil
}

// ---------------------------------------------------------------------------
// Drawing
// ---------------------------------------------------------------------------

func (b *PuregoBackend) Rectangle(_ context.Context, img uint64, x1, y1, x2, y2 int32, r, g, bl, a uint8, thickness int32) error {
	if b.rectangle == nil {
		return errNotSupported("rectangle")
	}
	return errorCode(abi.Rectangle, abi.ErrorCode(b.rectangle(img, x1, y1, x2, y2, r, g, bl, a, thickness)))
}

func (b *PuregoBackend) Circle(_ context.Context, img uint64, cx, cy, radius int32, r, g, bl, a uint8, thickness int32) error {
	if b.circle == nil {
		return errNotSupported("circle")
	}
	return errorCode(abi.Circle, abi.ErrorCode(b.circle(img, cx, cy, radius, r, g, bl, a, thickness)))
}

func (b *PuregoBackend) Line(_ context.Context, img uint64, x1, y1, x2, y2 int32, r, g, bl, a uint8, thickness int32) error {
	if b.line == nil {
		return errNotSupported("line")
	}
	return errorCode(abi.Line, abi.ErrorCode(b.line(img, x1, y1, x2, y2, r, g, bl, a, thickness)))
}

func (b *PuregoBackend) PutText(_ context.Context, img uint64, textData unsafe.Pointer, textLen int32, orgX, orgY, fontFace int32, fontScale float64, r, g, bl, a uint8, thickness, lineType, bottomLeftOrigin int32) error {
	if b.putText == nil {
		return errNotSupported("put_text")
	}
	return errorCode(abi.PutText, abi.ErrorCode(b.putText(img, textData, textLen, orgX, orgY, fontFace, fontScale, r, g, bl, a, thickness, lineType, bottomLeftOrigin)))
}

func (b *PuregoBackend) FillPoly(_ context.Context, img uint64, pts unsafe.Pointer, npts, ncontours int32, r, g, bl, a uint8, lineType, shift int32) error {
	if b.fillPoly == nil {
		return errNotSupported("fill_poly")
	}
	return errorCode(abi.FillPoly, abi.ErrorCode(b.fillPoly(img, pts, npts, ncontours, r, g, bl, a, lineType, shift)))
}

func (b *PuregoBackend) ArrowedLine(_ context.Context, img uint64, x1, y1, x2, y2 int32, r, g, bl, a uint8, thickness, lineType, shift int32, tipLength float64) error {
	if b.arrowedLine == nil {
		return errNotSupported("arrowed_line")
	}
	return errorCode(abi.ArrowedLine, abi.ErrorCode(b.arrowedLine(img, x1, y1, x2, y2, r, g, bl, a, thickness, lineType, shift, tipLength)))
}

// ---------------------------------------------------------------------------
// Contours
// ---------------------------------------------------------------------------

func (b *PuregoBackend) FindContours(_ context.Context, src, contoursVec uint64, mode, method, offX, offY int32) error {
	if b.findContours == nil {
		return errNotSupported("find_contours")
	}
	return errorCode(abi.FindContours, abi.ErrorCode(b.findContours(src, contoursVec, mode, method, offX, offY)))
}

func (b *PuregoBackend) DrawContours(_ context.Context, img, contoursVec uint64, contourIdx int32, r, g, bl, a uint8, thickness int32) error {
	if b.drawContours == nil {
		return errNotSupported("draw_contours")
	}
	return errorCode(abi.DrawContours, abi.ErrorCode(b.drawContours(img, contoursVec, contourIdx, r, g, bl, a, thickness)))
}

func (b *PuregoBackend) ContourArea(_ context.Context, vecHandle uint64) (float64, error) {
	if b.contourArea == nil {
		return 0, errNotSupported("contour_area")
	}
	return b.contourArea(vecHandle), nil
}

func (b *PuregoBackend) ArcLength(_ context.Context, vecHandle uint64, closed bool) (float64, error) {
	if b.arcLength == nil {
		return 0, errNotSupported("arc_length")
	}
	c := int32(0)
	if closed {
		c = 1
	}
	return b.arcLength(vecHandle, c), nil
}

func (b *PuregoBackend) BoundingRect(_ context.Context, vecHandle uint64) (int32, int32, int32, int32, error) {
	if b.boundingRect == nil {
		return 0, 0, 0, 0, errNotSupported("bounding_rect")
	}
	var x, y, w, h int32
	err := errorCode(abi.BoundingRect, abi.ErrorCode(b.boundingRect(vecHandle, &x, &y, &w, &h)))
	return x, y, w, h, err
}

func (b *PuregoBackend) MinEnclosingCircle(_ context.Context, vecHandle uint64) (float64, float64, float64, error) {
	if b.minEnclosingCircle == nil {
		return 0, 0, 0, errNotSupported("min_enclosing_circle")
	}
	var cx, cy, r float64
	err := errorCode(abi.MinEnclosingCircle, abi.ErrorCode(b.minEnclosingCircle(vecHandle, &cx, &cy, &r)))
	return cx, cy, r, err
}

func (b *PuregoBackend) Moments(_ context.Context, vecHandle uint64, binary bool) ([10]float64, error) {
	if b.moments == nil {
		return [10]float64{}, errNotSupported("moments")
	}
	var result [10]float64
	bin := int32(0)
	if binary {
		bin = 1
	}
	err := errorCode(abi.Moments, abi.ErrorCode(b.moments(vecHandle, bin,
		&result[0], &result[1], &result[2], &result[3], &result[4],
		&result[5], &result[6], &result[7], &result[8], &result[9],
	)))
	return result, err
}

// ---------------------------------------------------------------------------
// Hough
// ---------------------------------------------------------------------------

func (b *PuregoBackend) HoughLines(_ context.Context, src, vecHandle uint64, rho, theta float64, threshold int32, srn, stn float64) error {
	if b.houghLines == nil {
		return errNotSupported("hough_lines")
	}
	return errorCode(abi.HoughLines, abi.ErrorCode(b.houghLines(src, vecHandle, rho, theta, threshold, srn, stn)))
}

func (b *PuregoBackend) HoughLinesP(_ context.Context, src, vecHandle uint64, rho, theta float64, threshold int32, minLen, maxGap float64) error {
	if b.houghLinesP == nil {
		return errNotSupported("hough_lines_p")
	}
	return errorCode(abi.HoughLinesP, abi.ErrorCode(b.houghLinesP(src, vecHandle, rho, theta, threshold, minLen, maxGap)))
}

func (b *PuregoBackend) HoughCircles(_ context.Context, src, vecHandle uint64, method int32, dp, minDist, p1, p2 float64, minR, maxR int32) error {
	if b.houghCircles == nil {
		return errNotSupported("hough_circles")
	}
	return errorCode(abi.HoughCircles, abi.ErrorCode(b.houghCircles(src, vecHandle, method, dp, minDist, p1, p2, minR, maxR)))
}

// ---------------------------------------------------------------------------
// Warp
// ---------------------------------------------------------------------------

func (b *PuregoBackend) WarpAffine(_ context.Context, src, dst, m uint64, dstW, dstH int32) error {
	if b.warpAffine == nil {
		return errNotSupported("warp_affine")
	}
	return errorCode(abi.WarpAffine, abi.ErrorCode(b.warpAffine(src, dst, m, dstW, dstH)))
}

func (b *PuregoBackend) WarpPerspective(_ context.Context, src, dst, m uint64, dstW, dstH int32) error {
	if b.warpPerspective == nil {
		return errNotSupported("warp_perspective")
	}
	return errorCode(abi.WarpPerspective, abi.ErrorCode(b.warpPerspective(src, dst, m, dstW, dstH)))
}

func (b *PuregoBackend) GetRotationMatrix2D(_ context.Context, cx, cy, angle, scale float64) (uint64, error) {
	if b.getRotationMat2D == nil {
		return 0, errNotSupported("get_rotation_matrix2d")
	}
	h := b.getRotationMat2D(cx, cy, angle, scale)
	if h == 0 {
		return 0, fmt.Errorf("purego: get_rotation_matrix2d failed")
	}
	return h, nil
}

func (b *PuregoBackend) GetAffineTransform(_ context.Context, src0x, src0y, src1x, src1y, src2x, src2y, dst0x, dst0y, dst1x, dst1y, dst2x, dst2y float64) (uint64, error) {
	if b.getAffineTransform == nil {
		return 0, errNotSupported("get_affine_transform")
	}
	h := b.getAffineTransform(src0x, src0y, src1x, src1y, src2x, src2y, dst0x, dst0y, dst1x, dst1y, dst2x, dst2y)
	if h == 0 {
		return 0, fmt.Errorf("purego: get_affine_transform failed")
	}
	return h, nil
}

// ---------------------------------------------------------------------------
// Core ops
// ---------------------------------------------------------------------------

func (b *PuregoBackend) BitwiseAnd(_ context.Context, src1, src2, dst uint64) error {
	if b.bitwiseAnd == nil {
		return errNotSupported("bitwise_and")
	}
	return errorCode(abi.BitwiseAnd, abi.ErrorCode(b.bitwiseAnd(src1, src2, dst)))
}

func (b *PuregoBackend) BitwiseOr(_ context.Context, src1, src2, dst uint64) error {
	if b.bitwiseOr == nil {
		return errNotSupported("bitwise_or")
	}
	return errorCode(abi.BitwiseOr, abi.ErrorCode(b.bitwiseOr(src1, src2, dst)))
}

func (b *PuregoBackend) BitwiseXor(_ context.Context, src1, src2, dst uint64) error {
	if b.bitwiseXor == nil {
		return errNotSupported("bitwise_xor")
	}
	return errorCode(abi.BitwiseXor, abi.ErrorCode(b.bitwiseXor(src1, src2, dst)))
}

func (b *PuregoBackend) BitwiseNot(_ context.Context, src, dst uint64) error {
	if b.bitwiseNot == nil {
		return errNotSupported("bitwise_not")
	}
	return errorCode(abi.BitwiseNot, abi.ErrorCode(b.bitwiseNot(src, dst)))
}

func (b *PuregoBackend) Add(_ context.Context, src1, src2, dst uint64) error {
	if b.add == nil {
		return errNotSupported("add")
	}
	return errorCode(abi.Add, abi.ErrorCode(b.add(src1, src2, dst)))
}

func (b *PuregoBackend) Subtract(_ context.Context, src1, src2, dst uint64) error {
	if b.subtract == nil {
		return errNotSupported("subtract")
	}
	return errorCode(abi.Subtract, abi.ErrorCode(b.subtract(src1, src2, dst)))
}

func (b *PuregoBackend) Multiply(_ context.Context, src1, src2, dst uint64, scale float64, dtype int32) error {
	if b.multiply == nil {
		return errNotSupported("multiply")
	}
	return errorCode(abi.Multiply, abi.ErrorCode(b.multiply(src1, src2, dst, scale, dtype)))
}

func (b *PuregoBackend) Divide(_ context.Context, src1, src2, dst uint64, scale float64, dtype int32) error {
	if b.divide == nil {
		return errNotSupported("divide")
	}
	return errorCode(abi.Divide, abi.ErrorCode(b.divide(src1, src2, dst, scale, dtype)))
}

func (b *PuregoBackend) AbsDiff(_ context.Context, src1, src2, dst uint64) error {
	if b.absDiff == nil {
		return errNotSupported("abs_diff")
	}
	return errorCode(abi.AbsDiff, abi.ErrorCode(b.absDiff(src1, src2, dst)))
}

func (b *PuregoBackend) MinMaxLoc(_ context.Context, src uint64) (float64, float64, int32, int32, int32, int32, error) {
	if b.minMaxLoc == nil {
		return 0, 0, 0, 0, 0, 0, errNotSupported("min_max_loc")
	}
	var minVal, maxVal float64
	var minX, minY, maxX, maxY int32
	err := errorCode(abi.MinMaxLoc, abi.ErrorCode(b.minMaxLoc(src, &minVal, &maxVal, &minX, &minY, &maxX, &maxY)))
	return minVal, maxVal, minX, minY, maxX, maxY, err
}

func (b *PuregoBackend) MeanStdDev(_ context.Context, src, mean, stddev uint64) error {
	if b.meanStdDev == nil {
		return errNotSupported("mean_std_dev")
	}
	return errorCode(abi.MeanStdDev, abi.ErrorCode(b.meanStdDev(src, mean, stddev)))
}

func (b *PuregoBackend) CountNonZero(_ context.Context, src uint64) (int, error) {
	if b.countNonZero == nil {
		return 0, errNotSupported("count_non_zero")
	}
	r := b.countNonZero(src)
	if r < 0 {
		return 0, fmt.Errorf("count_non_zero: OpenCV error")
	}
	return int(r), nil
}

func (b *PuregoBackend) Split(_ context.Context, src, vecHandle uint64) error {
	if b.split == nil {
		return errNotSupported("split")
	}
	return errorCode(abi.Split, abi.ErrorCode(b.split(src, vecHandle)))
}

func (b *PuregoBackend) Merge(_ context.Context, vecHandle, dst uint64) error {
	if b.merge == nil {
		return errNotSupported("merge")
	}
	return errorCode(abi.Merge, abi.ErrorCode(b.merge(vecHandle, dst)))
}

// ---------------------------------------------------------------------------
// Vector helpers — points
// ---------------------------------------------------------------------------

func (b *PuregoBackend) VecPointsNew(_ context.Context) (uint64, error) {
	if b.vecPointsNew == nil {
		return 0, errNotSupported("vec_points_new")
	}
	return b.vecPointsNew(), nil
}

func (b *PuregoBackend) VecPointsPush(_ context.Context, handle uint64, x, y int32) error {
	if b.vecPointsPush == nil {
		return errNotSupported("vec_points_push")
	}
	b.vecPointsPush(handle, x, y)
	return nil
}

func (b *PuregoBackend) VecPointsLen(_ context.Context, handle uint64) (int, error) {
	if b.vecPointsLen == nil {
		return 0, errNotSupported("vec_points_len")
	}
	return int(b.vecPointsLen(handle)), nil
}

func (b *PuregoBackend) VecPointsGet(_ context.Context, handle uint64, idx int32) (int32, int32, error) {
	if b.vecPointsGet == nil {
		return 0, 0, errNotSupported("vec_points_get")
	}
	var x, y int32
	err := errorCode(abi.VecGetPoint, abi.ErrorCode(b.vecPointsGet(handle, idx, &x, &y)))
	return x, y, err
}

func (b *PuregoBackend) VecPointsDelete(_ context.Context, handle uint64) {
	if b.vecPointsDelete != nil {
		b.vecPointsDelete(handle)
	}
}

// ---------------------------------------------------------------------------
// Vector helpers — vector of points
// ---------------------------------------------------------------------------

func (b *PuregoBackend) VecVecPointsNew(_ context.Context) (uint64, error) {
	if b.vecVecPointsNew == nil {
		return 0, errNotSupported("vec_vec_points_new")
	}
	return b.vecVecPointsNew(), nil
}

func (b *PuregoBackend) VecVecPointsPush(_ context.Context, handle, contourHandle uint64) error {
	if b.vecVecPointsPush == nil {
		return errNotSupported("vec_vec_points_push")
	}
	b.vecVecPointsPush(handle, contourHandle)
	return nil
}

func (b *PuregoBackend) VecVecPointsLen(_ context.Context, handle uint64) (int, error) {
	if b.vecVecPointsLen == nil {
		return 0, errNotSupported("vec_vec_points_len")
	}
	return int(b.vecVecPointsLen(handle)), nil
}

func (b *PuregoBackend) VecVecPointsGet(_ context.Context, handle uint64, idx int32) (uint64, error) {
	if b.vecVecPointsGet == nil {
		return 0, errNotSupported("vec_vec_points_get")
	}
	h := b.vecVecPointsGet(handle, idx)
	if h == 0 {
		return 0, fmt.Errorf("vec_vec_points_get: index out of range")
	}
	return h, nil
}

func (b *PuregoBackend) VecVecPointsDelete(_ context.Context, handle uint64) {
	if b.vecVecPointsDelete != nil {
		b.vecVecPointsDelete(handle)
	}
}

// ---------------------------------------------------------------------------
// Vector helpers — double
// ---------------------------------------------------------------------------

func (b *PuregoBackend) VecDoubleNew(_ context.Context) (uint64, error) {
	if b.vecDoubleNew == nil {
		return 0, errNotSupported("vec_double_new")
	}
	return b.vecDoubleNew(), nil
}

func (b *PuregoBackend) VecDoubleGet(_ context.Context, handle uint64, idx int32) (float64, error) {
	if b.vecDoubleGet == nil {
		return 0, errNotSupported("vec_double_get")
	}
	return b.vecDoubleGet(handle, idx), nil
}

func (b *PuregoBackend) VecDoubleLen(_ context.Context, handle uint64) (int, error) {
	if b.vecDoubleLen == nil {
		return 0, errNotSupported("vec_double_len")
	}
	return int(b.vecDoubleLen(handle)), nil
}

func (b *PuregoBackend) VecDoubleDelete(_ context.Context, handle uint64) {
	if b.vecDoubleDelete != nil {
		b.vecDoubleDelete(handle)
	}
}

// ---------------------------------------------------------------------------
// Vector helpers — int
// ---------------------------------------------------------------------------

func (b *PuregoBackend) VecIntNew(_ context.Context) (uint64, error) {
	if b.vecIntNew == nil {
		return 0, errNotSupported("vec_int_new")
	}
	return b.vecIntNew(), nil
}

func (b *PuregoBackend) VecIntGet(_ context.Context, handle uint64, idx int32) (int32, error) {
	if b.vecIntGet == nil {
		return 0, errNotSupported("vec_int_get")
	}
	return b.vecIntGet(handle, idx), nil
}

func (b *PuregoBackend) VecIntLen(_ context.Context, handle uint64) (int, error) {
	if b.vecIntLen == nil {
		return 0, errNotSupported("vec_int_len")
	}
	return int(b.vecIntLen(handle)), nil
}

func (b *PuregoBackend) VecIntDelete(_ context.Context, handle uint64) {
	if b.vecIntDelete != nil {
		b.vecIntDelete(handle)
	}
}

// ---------------------------------------------------------------------------
// Vector helpers — mat
// ---------------------------------------------------------------------------

func (b *PuregoBackend) VecMatNew(_ context.Context) (uint64, error) {
	if b.vecMatNew == nil {
		return 0, errNotSupported("vec_mat_new")
	}
	return b.vecMatNew(), nil
}

func (b *PuregoBackend) VecMatPush(_ context.Context, handle, matHandle uint64) error {
	if b.vecMatPush == nil {
		return errNotSupported("vec_mat_push")
	}
	b.vecMatPush(handle, matHandle)
	return nil
}

func (b *PuregoBackend) VecMatLen(_ context.Context, handle uint64) (int, error) {
	if b.vecMatLen == nil {
		return 0, errNotSupported("vec_mat_len")
	}
	return int(b.vecMatLen(handle)), nil
}

func (b *PuregoBackend) VecMatGet(_ context.Context, handle uint64, idx int32) (uint64, error) {
	if b.vecMatGet == nil {
		return 0, errNotSupported("vec_mat_get")
	}
	h := b.vecMatGet(handle, idx)
	if h == 0 {
		return 0, fmt.Errorf("vec_mat_get: index out of range")
	}
	return h, nil
}

func (b *PuregoBackend) VecMatDelete(_ context.Context, handle uint64) {
	if b.vecMatDelete != nil {
		b.vecMatDelete(handle)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func errNotSupported(op string) error {
	return fmt.Errorf("%s: not supported by this library build", op)
}

// errorCode translates ABI error codes to Go errors.
func errorCode(op string, code abi.ErrorCode) error {
	switch code {
	case abi.OK:
		return nil
	case abi.ErrInvalidArgument:
		return fmt.Errorf("%s: invalid argument", op)
	case abi.ErrInvalidHandle:
		return fmt.Errorf("%s: invalid handle", op)
	case abi.ErrOutOfMemory:
		return fmt.Errorf("%s: out of memory", op)
	case abi.ErrOpenCV:
		return fmt.Errorf("%s: OpenCV error", op)
	case abi.ErrUnsupported:
		return fmt.Errorf("%s: unsupported", op)
	default:
		return fmt.Errorf("%s: unknown error code %d", op, code)
	}
}
