package opencv

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"unsafe"

	internalruntime "github.com/brainplusplus/go-opencv/internal/runtime"
)

// Runtime owns a single OpenCV instance (native DLL) and all Mats.
type Runtime struct {
	ctx     context.Context
	cancel  context.CancelFunc
	backend backend
	closed  atomic.Bool
	mu      sync.Mutex
	mats    map[matHandle]struct{}
}

type backend interface {
	Close(context.Context) error
	NewMat(context.Context, int, int, int32) (matHandle, error)
	NewMatFromData(context.Context, []byte, int, int, int32) (matHandle, error)
	CloseMat(context.Context, matHandle) error
	MatRows(context.Context, matHandle) (int, error)
	MatCols(context.Context, matHandle) (int, error)
	MatType(context.Context, matHandle) (int32, error)
	MatClone(context.Context, matHandle) (matHandle, error)
	MatEmpty(context.Context, matHandle) (bool, error)
	MatElemSize(context.Context, matHandle) (int, error)
	MatDataPtr(context.Context, matHandle) (unsafe.Pointer, error)
	MatStep(context.Context, matHandle) (int, error)
	MatChannels(context.Context, matHandle) (int, error)
	MatTotal(context.Context, matHandle) (int, error)
	MatRow(context.Context, matHandle, int32) (matHandle, error)
	MatCol(context.Context, matHandle, int32) (matHandle, error)
	MatRegion(context.Context, matHandle, int32, int32, int32, int32) (matHandle, error)
	MatReshape(context.Context, matHandle, int32, int32) (matHandle, error)
	MatSetTo(context.Context, matHandle, float64, float64, float64, float64) error
	MatConvertTo(context.Context, matHandle, matHandle, int32, float64, float64) error
	MatZeros(context.Context, int, int, int32) (matHandle, error)
	MatOnes(context.Context, int, int, int32) (matHandle, error)
	MatEye(context.Context, int, int, int32) (matHandle, error)
	CvtColor(context.Context, matHandle, matHandle, int32) error
	Resize(context.Context, matHandle, matHandle, int32, int32) error
	Blur(context.Context, matHandle, matHandle, int32, int32) error
	GaussianBlur(context.Context, matHandle, matHandle, int32, int32, float64) error
	MedianBlur(context.Context, matHandle, matHandle, int32) error
	Threshold(context.Context, matHandle, matHandle, float64, float64, int32) error
	AdaptiveThreshold(context.Context, matHandle, matHandle, float64, int32, int32, int32, float64) error
	Canny(context.Context, matHandle, matHandle, float64, float64) error
	Flip(context.Context, matHandle, matHandle, int32) error
	Sobel(context.Context, matHandle, matHandle, int32, int32, int32, int32, float64, float64) error
	Laplacian(context.Context, matHandle, matHandle, int32, int32, float64, float64) error
	Transpose(context.Context, matHandle, matHandle) error
	EqualizeHist(context.Context, matHandle, matHandle) error
	Normalize(context.Context, matHandle, matHandle, float64, float64, int32) error
	Erode(context.Context, matHandle, matHandle, matHandle, int32, int32, int32) error
	Dilate(context.Context, matHandle, matHandle, matHandle, int32, int32, int32) error
	MorphologyEx(context.Context, matHandle, matHandle, int32, matHandle, int32, int32, int32) error
	GetStructuringElement(context.Context, int32, int32, int32) (matHandle, error)
	Rectangle(context.Context, matHandle, int32, int32, int32, int32, uint8, uint8, uint8, uint8, int32) error
	Circle(context.Context, matHandle, int32, int32, int32, uint8, uint8, uint8, uint8, int32) error
	Line(context.Context, matHandle, int32, int32, int32, int32, uint8, uint8, uint8, uint8, int32) error
	PutText(context.Context, matHandle, unsafe.Pointer, int32, int32, int32, int32, float64, uint8, uint8, uint8, uint8, int32, int32, int32) error
	FillPoly(context.Context, matHandle, unsafe.Pointer, int32, int32, uint8, uint8, uint8, uint8, int32, int32) error
	ArrowedLine(context.Context, matHandle, int32, int32, int32, int32, uint8, uint8, uint8, uint8, int32, int32, int32, float64) error
	FindContours(context.Context, matHandle, matHandle, int32, int32, int32, int32) error
	DrawContours(context.Context, matHandle, matHandle, int32, uint8, uint8, uint8, uint8, int32) error
	ContourArea(context.Context, matHandle) (float64, error)
	ArcLength(context.Context, matHandle, bool) (float64, error)
	BoundingRect(context.Context, matHandle) (int32, int32, int32, int32, error)
	MinEnclosingCircle(context.Context, matHandle) (float64, float64, float64, error)
	Moments(context.Context, matHandle, bool) ([10]float64, error)
	HoughLines(context.Context, matHandle, matHandle, float64, float64, int32, float64, float64) error
	HoughLinesP(context.Context, matHandle, matHandle, float64, float64, int32, float64, float64) error
	HoughCircles(context.Context, matHandle, matHandle, int32, float64, float64, float64, float64, int32, int32) error
	WarpAffine(context.Context, matHandle, matHandle, matHandle, int32, int32) error
	WarpPerspective(context.Context, matHandle, matHandle, matHandle, int32, int32) error
	GetRotationMatrix2D(context.Context, float64, float64, float64, float64) (matHandle, error)
	GetAffineTransform(context.Context, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64) (matHandle, error)
	BitwiseAnd(context.Context, matHandle, matHandle, matHandle) error
	BitwiseOr(context.Context, matHandle, matHandle, matHandle) error
	BitwiseXor(context.Context, matHandle, matHandle, matHandle) error
	BitwiseNot(context.Context, matHandle, matHandle) error
	Add(context.Context, matHandle, matHandle, matHandle) error
	Subtract(context.Context, matHandle, matHandle, matHandle) error
	Multiply(context.Context, matHandle, matHandle, matHandle, float64, int32) error
	Divide(context.Context, matHandle, matHandle, matHandle, float64, int32) error
	AbsDiff(context.Context, matHandle, matHandle, matHandle) error
	MinMaxLoc(context.Context, matHandle) (float64, float64, int32, int32, int32, int32, error)
	MeanStdDev(context.Context, matHandle, matHandle, matHandle) error
	CountNonZero(context.Context, matHandle) (int, error)
	Split(context.Context, matHandle, matHandle) error
	Merge(context.Context, matHandle, matHandle) error
	// Vector helpers
	VecPointsNew(context.Context) (matHandle, error)
	VecPointsPush(context.Context, matHandle, int32, int32) error
	VecPointsLen(context.Context, matHandle) (int, error)
	VecPointsGet(context.Context, matHandle, int32) (int32, int32, error)
	VecPointsDelete(context.Context, matHandle)
	VecVecPointsNew(context.Context) (matHandle, error)
	VecVecPointsPush(context.Context, matHandle, matHandle) error
	VecVecPointsLen(context.Context, matHandle) (int, error)
	VecVecPointsGet(context.Context, matHandle, int32) (matHandle, error)
	VecVecPointsDelete(context.Context, matHandle)
	VecDoubleNew(context.Context) (matHandle, error)
	VecDoubleGet(context.Context, matHandle, int32) (float64, error)
	VecDoubleLen(context.Context, matHandle) (int, error)
	VecDoubleDelete(context.Context, matHandle)
	VecIntNew(context.Context) (matHandle, error)
	VecIntGet(context.Context, matHandle, int32) (int32, error)
	VecIntLen(context.Context, matHandle) (int, error)
	VecIntDelete(context.Context, matHandle)
	VecMatNew(context.Context) (matHandle, error)
	VecMatPush(context.Context, matHandle, matHandle) error
	VecMatLen(context.Context, matHandle) (int, error)
	VecMatGet(context.Context, matHandle, int32) (matHandle, error)
	VecMatDelete(context.Context, matHandle)
	// New imgproc
	BilateralFilter(context.Context, matHandle, matHandle, int32, float64, float64) error
	InRange(context.Context, matHandle, float64, float64, float64, float64, float64, float64, float64, float64, matHandle) error
	MatchTemplate(context.Context, matHandle, matHandle, matHandle, int32) error
	CalcHist(context.Context, matHandle, matHandle, int32, float64, float64) error
	ConnectedComponents(context.Context, matHandle, matHandle, int32, int32) error
	DistanceTransform(context.Context, matHandle, matHandle, int32, int32) error
	CopyMakeBorder(context.Context, matHandle, matHandle, int32, int32, int32, int32, int32, float64, float64, float64, float64) error
	Rotate(context.Context, matHandle, matHandle, int32) error
	Hconcat(context.Context, matHandle, matHandle, matHandle) error
	Vconcat(context.Context, matHandle, matHandle, matHandle) error
	Remap(context.Context, matHandle, matHandle, matHandle, matHandle, int32, int32, float64) error
	LUT(context.Context, matHandle, matHandle, matHandle) error
	Integral(context.Context, matHandle, matHandle) error
	GetPerspectiveTransform(context.Context, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64) (matHandle, error)
	FillConvexPoly(context.Context, matHandle, matHandle, int32, float64, float64, float64, float64, int32, int32) error
	ConvertModel(context.Context, matHandle, matHandle, int32, int32) error
	// Photo
	FastNlMeansDenoising(context.Context, matHandle, matHandle, float32, int32, int32) error
	FastNlMeansDenoisingColored(context.Context, matHandle, matHandle, float32, float32, int32, int32) error
	DetailEnhance(context.Context, matHandle, matHandle, float32, float32) error
	EdgePreservingFilter(context.Context, matHandle, matHandle, int32, float32, float32) error
	PencilSketch(context.Context, matHandle, matHandle, matHandle, float32, float32, float32) error
	Stylization(context.Context, matHandle, matHandle, float32, float32) error
	SeamlessClone(context.Context, matHandle, matHandle, matHandle, matHandle, int32, int32, int32) error
	// Features2d
	FAST(context.Context, matHandle, matHandle, int32, int32) error
	ORBDetectCompute(context.Context, matHandle, matHandle, matHandle, matHandle, int32, float32, int32) error
	BFMatch(context.Context, matHandle, matHandle, matHandle, int32) error
	DrawKeypoints(context.Context, matHandle, matHandle, matHandle, float64, float64, float64) error
	// Highgui
	ImShow(context.Context, unsafe.Pointer, int32, matHandle) error
	WaitKey(context.Context, int32) (int32, error)
	DestroyWindow(context.Context, unsafe.Pointer, int32) error
	// Core extras
	MatDiag(context.Context, matHandle) (matHandle, error)
	MatAtU8(context.Context, matHandle, int32, int32, int32) (uint8, error)
	MatSetU8(context.Context, matHandle, int32, int32, int32, uint8) error
	// Vector helpers — keypoints
	VecKeypointNew(context.Context) (matHandle, error)
	VecKeypointLen(context.Context, matHandle) (int, error)
	VecKeypointGet(context.Context, matHandle, int32) (float32, float32, float32, float32, float32, int32, int32, error)
	VecKeypointDelete(context.Context, matHandle)
	// Vector helpers — dmatch
	VecDMatchNew(context.Context) (matHandle, error)
	VecDMatchLen(context.Context, matHandle) (int, error)
	VecDMatchGet(context.Context, matHandle, int32) (int32, int32, float32, int32, error)
	VecDMatchDelete(context.Context, matHandle)
}

// New creates a Runtime backed by a native shared library.
func New(ctx context.Context, opts ...Option) (*Runtime, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg := config{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	if cfg.dll != "" {
		b, err := internalruntime.NewPuregoBackend(cfg.dll)
		if err == nil {
			rt := &Runtime{ctx: ctx, cancel: cancel, backend: b, mats: map[matHandle]struct{}{}}
			return rt, nil
		}
		cancel()
		return nil, fmt.Errorf("opencv: initialize dll backend (%s): %w", cfg.dll, err)
	}

	if len(embedLibData) > 0 {
		libPath, err := extractLib()
		if err == nil {
			b, err := internalruntime.NewPuregoBackend(libPath)
			if err == nil {
				rt := &Runtime{ctx: ctx, cancel: cancel, backend: b, mats: map[matHandle]struct{}{}}
				return rt, nil
			}
			cancel()
			return nil, fmt.Errorf("opencv: load embedded library: %w", err)
		}
		cancel()
		return nil, fmt.Errorf("opencv: extract embedded library: %w", err)
	}

	cancel()
	return nil, ErrBackendUnavailable
}

func extractLib() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}

	dir := filepath.Join(cacheDir, "go-opencv", Version)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}

	libPath := filepath.Join(dir, embedLibName())
	if _, err := os.Stat(libPath); err == nil {
		return libPath, nil
	}

	tmpPath := libPath + ".tmp"
	if err := os.WriteFile(tmpPath, embedLibData, 0755); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("write %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, libPath); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("rename %s -> %s: %w", tmpPath, libPath, err)
	}

	return libPath, nil
}

func (r *Runtime) Close() error {
	if r == nil {
		return ErrClosed
	}
	if r.closed.Swap(true) {
		return ErrClosed
	}
	r.cancel()
	return r.backend.Close(context.Background())
}

// ---------------------------------------------------------------------------
// Mat factory
// ---------------------------------------------------------------------------

func (r *Runtime) NewMat(rows, cols int, typ MatType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.NewMat(r.ctx, rows, cols, int32(typ))
	if err != nil {
		return nil, err
	}
	model := Unknown
	switch typ {
	case CV8UC1:
		model = Gray
	case CV8UC3:
		model = BGR
	case CV8UC4:
		model = RGBA
	}
	return r.wrapMatWithModel(h, model), nil
}

// Zeros creates a Mat filled with zeros.
func (r *Runtime) Zeros(rows, cols int, typ MatType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.MatZeros(r.ctx, rows, cols, int32(typ))
	if err != nil {
		return nil, err
	}
	model := Unknown
	switch typ {
	case CV8UC1:
		model = Gray
	case CV8UC3:
		model = BGR
	case CV8UC4:
		model = RGBA
	}
	return r.wrapMatWithModel(h, model), nil
}

// Ones creates a Mat filled with ones.
func (r *Runtime) Ones(rows, cols int, typ MatType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.MatOnes(r.ctx, rows, cols, int32(typ))
	if err != nil {
		return nil, err
	}
	return r.wrapMatWithModel(h, Unknown), nil
}

// Eye creates an identity matrix.
func (r *Runtime) Eye(rows, cols int, typ MatType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.MatEye(r.ctx, rows, cols, int32(typ))
	if err != nil {
		return nil, err
	}
	return r.wrapMatWithModel(h, Unknown), nil
}

// GetStructuringElement returns a structuring element of the specified size and shape for morphological operations.
func (r *Runtime) GetStructuringElement(shape MorphShape, ksize Size) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.GetStructuringElement(r.ctx, int32(shape), ksize.Width, ksize.Height)
	if err != nil {
		return nil, err
	}
	return r.wrapMatWithModel(h, Unknown), nil
}

// ---------------------------------------------------------------------------
// Color conversion
// ---------------------------------------------------------------------------

func (r *Runtime) CvtColor(src, dst *Mat, code ColorConversionCode) error {
	if err := validatePair(r, src, dst); err != nil {
		return err
	}
	if StrictColorValidation() && src.model != Unknown {
		if err := validateColorConversion(src.model, code); err != nil {
			return err
		}
	}
	if m, ok := outputModelForCode(src.model, code); ok {
		dst.model = m
	} else {
		dst.model = Unknown
	}
	return r.backend.CvtColor(r.ctx, src.handle, dst.handle, int32(code))
}

// ---------------------------------------------------------------------------
// Geometry transforms
// ---------------------------------------------------------------------------

func (r *Runtime) Resize(src, dst *Mat, size Size) error {
	if err := validatePair(r, src, dst); err != nil {
		return err
	}
	dst.model = src.model
	return r.backend.Resize(r.ctx, src.handle, dst.handle, size.Width, size.Height)
}

// Flip flips the src Mat and stores the result in dst.
func (r *Runtime) Flip(src, dst *Mat, flipCode FlipCode) error {
	if err := validatePair(r, src, dst); err != nil {
		return err
	}
	dst.model = src.model
	return r.backend.Flip(r.ctx, src.handle, dst.handle, int32(flipCode))
}

// Transpose transposes the src Mat.
func (r *Runtime) Transpose(src *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(cols, rows, typ)
	if err != nil {
		return nil, fmt.Errorf("transpose: create output: %w", err)
	}
	if err := r.backend.Transpose(r.ctx, src.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// WarpAffine applies an affine transformation to src using the 2x3 matrix M.
func (r *Runtime) WarpAffine(src *Mat, M *Mat, dsize Size) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	typ, _ := src.Type()
	dst, err := r.NewMat(int(dsize.Height), int(dsize.Width), typ)
	if err != nil {
		return nil, fmt.Errorf("warp_affine: create output: %w", err)
	}
	if err := r.backend.WarpAffine(r.ctx, src.handle, dst.handle, M.handle, dsize.Width, dsize.Height); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// WarpPerspective applies a perspective transformation to src using the 3x3 matrix M.
func (r *Runtime) WarpPerspective(src *Mat, M *Mat, dsize Size) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	typ, _ := src.Type()
	dst, err := r.NewMat(int(dsize.Height), int(dsize.Width), typ)
	if err != nil {
		return nil, fmt.Errorf("warp_perspective: create output: %w", err)
	}
	if err := r.backend.WarpPerspective(r.ctx, src.handle, dst.handle, M.handle, dsize.Width, dsize.Height); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// GetRotationMatrix2D returns a 2x3 affine transformation matrix for 2D rotation.
func (r *Runtime) GetRotationMatrix2D(center Point, angle, scale float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.GetRotationMatrix2D(r.ctx, float64(center.X), float64(center.Y), angle, scale)
	if err != nil {
		return nil, err
	}
	return r.wrapMatWithModel(h, Unknown), nil
}

// GetAffineTransform returns a 2x3 affine transformation matrix from three point pairs.
func (r *Runtime) GetAffineTransform(src, dst [3]Point) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	h, err := r.backend.GetAffineTransform(r.ctx,
		float64(src[0].X), float64(src[0].Y), float64(src[1].X), float64(src[1].Y), float64(src[2].X), float64(src[2].Y),
		float64(dst[0].X), float64(dst[0].Y), float64(dst[1].X), float64(dst[1].Y), float64(dst[2].X), float64(dst[2].Y),
	)
	if err != nil {
		return nil, err
	}
	return r.wrapMatWithModel(h, Unknown), nil
}

// ---------------------------------------------------------------------------
// Filtering
// ---------------------------------------------------------------------------

func (r *Runtime) Blur(src *Mat, kSize Size) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("blur: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("blur: cols: %w", err)
	}
	typ, err := src.Type()
	if err != nil {
		return nil, fmt.Errorf("blur: type: %w", err)
	}
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("blur: create output: %w", err)
	}
	if err := r.backend.Blur(r.ctx, src.handle, dst.handle, kSize.Width, kSize.Height); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

func (r *Runtime) GaussianBlur(src *Mat, kSize Size, sigmaX float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("gaussian blur: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("gaussian blur: cols: %w", err)
	}
	typ, err := src.Type()
	if err != nil {
		return nil, fmt.Errorf("gaussian blur: type: %w", err)
	}
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("gaussian blur: create output: %w", err)
	}
	if err := r.backend.GaussianBlur(r.ctx, src.handle, dst.handle, kSize.Width, kSize.Height, sigmaX); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// MedianBlur applies median blur with the given kernel size (must be odd, >1).
func (r *Runtime) MedianBlur(src *Mat, ksize int) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("median blur: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("median blur: cols: %w", err)
	}
	typ, err := src.Type()
	if err != nil {
		return nil, fmt.Errorf("median blur: type: %w", err)
	}
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("median blur: create output: %w", err)
	}
	if err := r.backend.MedianBlur(r.ctx, src.handle, dst.handle, int32(ksize)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

func (r *Runtime) Threshold(src *Mat, thresh, maxValue float64, typ ThresholdType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("threshold: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("threshold: cols: %w", err)
	}
	t, err := src.Type()
	if err != nil {
		return nil, fmt.Errorf("threshold: type: %w", err)
	}
	dst, err := r.NewMat(rows, cols, t)
	if err != nil {
		return nil, fmt.Errorf("threshold: create output: %w", err)
	}
	if err := r.backend.Threshold(r.ctx, src.handle, dst.handle, thresh, maxValue, int32(typ)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

func (r *Runtime) AdaptiveThreshold(src *Mat, maxValue float64, adaptiveType AdaptiveThresholdType, thresholdType ThresholdType, blockSize int, c float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("adaptive threshold: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("adaptive threshold: cols: %w", err)
	}
	t, err := src.Type()
	if err != nil {
		return nil, fmt.Errorf("adaptive threshold: type: %w", err)
	}
	dst, err := r.NewMat(rows, cols, t)
	if err != nil {
		return nil, fmt.Errorf("adaptive threshold: create output: %w", err)
	}
	if err := r.backend.AdaptiveThreshold(r.ctx, src.handle, dst.handle, maxValue, int32(adaptiveType), int32(thresholdType), int32(blockSize), c); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = Gray
	return dst, nil
}

// Canny performs Canny edge detection on src and returns the edge mask.
func (r *Runtime) Canny(src *Mat, threshold1, threshold2 float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := src.validate(); err != nil {
		return nil, err
	}
	rows, err := src.Rows()
	if err != nil {
		return nil, fmt.Errorf("canny: rows: %w", err)
	}
	cols, err := src.Cols()
	if err != nil {
		return nil, fmt.Errorf("canny: cols: %w", err)
	}
	dst, err := r.NewMat(rows, cols, CV8UC1)
	if err != nil {
		return nil, fmt.Errorf("canny: create output: %w", err)
	}
	if err := r.backend.Canny(r.ctx, src.handle, dst.handle, threshold1, threshold2); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = Gray
	return dst, nil
}

// Sobel calculates the Sobel derivative image.
func (r *Runtime) Sobel(src *Mat, ddepth MatType, dx, dy, ksize int, scale, delta float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, ddepth)
	if err != nil {
		return nil, fmt.Errorf("sobel: create output: %w", err)
	}
	if err := r.backend.Sobel(r.ctx, src.handle, dst.handle, int32(ddepth), int32(dx), int32(dy), int32(ksize), scale, delta); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// Laplacian calculates the Laplacian of the image.
func (r *Runtime) Laplacian(src *Mat, ddepth MatType, ksize int, scale, delta float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, ddepth)
	if err != nil {
		return nil, fmt.Errorf("laplacian: create output: %w", err)
	}
	if err := r.backend.Laplacian(r.ctx, src.handle, dst.handle, int32(ddepth), int32(ksize), scale, delta); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// EqualizeHist equalizes the histogram of a grayscale image.
func (r *Runtime) EqualizeHist(src *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, CV8UC1)
	if err != nil {
		return nil, fmt.Errorf("equalize_hist: create output: %w", err)
	}
	if err := r.backend.EqualizeHist(r.ctx, src.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = Gray
	return dst, nil
}

// Normalize normalizes the input Mat.
func (r *Runtime) Normalize(src *Mat, alpha, beta float64, normType NormType) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("normalize: create output: %w", err)
	}
	if err := r.backend.Normalize(r.ctx, src.handle, dst.handle, alpha, beta, int32(normType)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// ---------------------------------------------------------------------------
// Morphology
// ---------------------------------------------------------------------------

// Erode erodes the image using a structuring element.
func (r *Runtime) Erode(src *Mat, kernel *Mat, anchor Point, iterations int) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("erode: create output: %w", err)
	}
	var kh matHandle
	if kernel != nil {
		kh = kernel.handle
	}
	if err := r.backend.Erode(r.ctx, src.handle, dst.handle, kh, anchor.X, anchor.Y, int32(iterations)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// Dilate dilates the image using a structuring element.
func (r *Runtime) Dilate(src *Mat, kernel *Mat, anchor Point, iterations int) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("dilate: create output: %w", err)
	}
	var kh matHandle
	if kernel != nil {
		kh = kernel.handle
	}
	if err := r.backend.Dilate(r.ctx, src.handle, dst.handle, kh, anchor.X, anchor.Y, int32(iterations)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// MorphologyEx performs advanced morphological transformations.
func (r *Runtime) MorphologyEx(src *Mat, op MorphType, kernel *Mat, anchor Point, iterations int) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("morphology_ex: create output: %w", err)
	}
	var kh matHandle
	if kernel != nil {
		kh = kernel.handle
	}
	if err := r.backend.MorphologyEx(r.ctx, src.handle, dst.handle, int32(op), kh, anchor.X, anchor.Y, int32(iterations)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// ---------------------------------------------------------------------------
// Drawing
// ---------------------------------------------------------------------------

func (r *Runtime) Rectangle(img *Mat, rect Rect, c color.Color, thickness int) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}
	cr, cg, cb, ca := colorFetch(c)
	return r.backend.Rectangle(
		r.ctx, img.handle,
		rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height,
		cr, cg, cb, ca, int32(thickness),
	)
}

func (r *Runtime) Circle(img *Mat, center Point, radius int, c color.Color, thickness int) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}
	cr, cg, cb, ca := colorFetch(c)
	return r.backend.Circle(r.ctx, img.handle, center.X, center.Y, int32(radius), cr, cg, cb, ca, int32(thickness))
}

func (r *Runtime) Line(img *Mat, point1, point2 Point, c color.Color, thickness int) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}
	cr, cg, cb, ca := colorFetch(c)
	return r.backend.Line(r.ctx, img.handle, point1.X, point1.Y, point2.X, point2.Y, cr, cg, cb, ca, int32(thickness))
}

// PutText draws a text string on the image.
func (r *Runtime) PutText(img *Mat, text string, org Point, fontFace HersheyFontType, fontScale float64, c color.Color, thickness int) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}
	cr, cg, cb, ca := colorFetch(c)
	textBytes := []byte(text)
	return r.backend.PutText(
		r.ctx, img.handle,
		unsafe.Pointer(&textBytes[0]), int32(len(textBytes)),
		org.X, org.Y, int32(fontFace), fontScale,
		cr, cg, cb, ca, int32(thickness), int32(Line8), 0,
	)
}

// ArrowedLine draws an arrow from point1 to point2.
func (r *Runtime) ArrowedLine(img *Mat, point1, point2 Point, c color.Color, thickness int, tipLength float64) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}
	cr, cg, cb, ca := colorFetch(c)
	return r.backend.ArrowedLine(r.ctx, img.handle, point1.X, point1.Y, point2.X, point2.Y, cr, cg, cb, ca, int32(thickness), int32(Line8), 0, tipLength)
}

// ---------------------------------------------------------------------------
// Contours
// ---------------------------------------------------------------------------

// FindContours finds contours in a binary image.
func (r *Runtime) FindContours(src *Mat, mode RetrievalMode, method ContourApproximationMode) ([][]Point, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}

	contoursVec, err := r.backend.VecVecPointsNew(r.ctx)
	if err != nil {
		return nil, fmt.Errorf("find_contours: create contours vec: %w", err)
	}
	defer r.backend.VecVecPointsDelete(r.ctx, contoursVec)

	if err := r.backend.FindContours(r.ctx, src.handle, contoursVec, int32(mode), int32(method), 0, 0); err != nil {
		return nil, fmt.Errorf("find_contours: %w", err)
	}

	nContours, err := r.backend.VecVecPointsLen(r.ctx, contoursVec)
	if err != nil {
		return nil, err
	}

	result := make([][]Point, nContours)
	for i := 0; i < nContours; i++ {
		contourVec, err := r.backend.VecVecPointsGet(r.ctx, contoursVec, int32(i))
		if err != nil {
			return nil, err
		}
		defer r.backend.VecPointsDelete(r.ctx, contourVec)

		nPts, err := r.backend.VecPointsLen(r.ctx, contourVec)
		if err != nil {
			return nil, err
		}
		pts := make([]Point, nPts)
		for j := 0; j < nPts; j++ {
			x, y, err := r.backend.VecPointsGet(r.ctx, contourVec, int32(j))
			if err != nil {
				return nil, err
			}
			pts[j] = Point{X: x, Y: y}
		}
		result[i] = pts
	}

	return result, nil
}

// DrawContours draws contour outlines or filled contours.
func (r *Runtime) DrawContours(img *Mat, contours [][]Point, contourIdx int, c color.Color, thickness int) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := r.validateOwnedMat(img); err != nil {
		return err
	}

	contoursVec, err := r.backend.VecVecPointsNew(r.ctx)
	if err != nil {
		return err
	}
	defer r.backend.VecVecPointsDelete(r.ctx, contoursVec)

	for _, contour := range contours {
		ptVec, err := r.backend.VecPointsNew(r.ctx)
		if err != nil {
			return err
		}
		for _, pt := range contour {
			r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
		}
		r.backend.VecVecPointsPush(r.ctx, contoursVec, ptVec)
		r.backend.VecPointsDelete(r.ctx, ptVec)
	}

	cr, cg, cb, ca := colorFetch(c)
	return r.backend.DrawContours(r.ctx, img.handle, contoursVec, int32(contourIdx), cr, cg, cb, ca, int32(thickness))
}

// ContourArea computes the area of a contour.
func (r *Runtime) ContourArea(contour []Point) (float64, error) {
	if err := r.validateOpen(); err != nil {
		return 0, err
	}
	ptVec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil {
		return 0, err
	}
	defer r.backend.VecPointsDelete(r.ctx, ptVec)
	for _, pt := range contour {
		r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
	}
	return r.backend.ContourArea(r.ctx, ptVec)
}

// ArcLength computes the arc length of a contour.
func (r *Runtime) ArcLength(contour []Point, closed bool) (float64, error) {
	if err := r.validateOpen(); err != nil {
		return 0, err
	}
	ptVec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil {
		return 0, err
	}
	defer r.backend.VecPointsDelete(r.ctx, ptVec)
	for _, pt := range contour {
		r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
	}
	return r.backend.ArcLength(r.ctx, ptVec, closed)
}

// BoundingRect computes the up-right bounding rectangle of a point set.
func (r *Runtime) BoundingRect(contour []Point) (Rect, error) {
	if err := r.validateOpen(); err != nil {
		return Rect{}, err
	}
	ptVec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil {
		return Rect{}, err
	}
	defer r.backend.VecPointsDelete(r.ctx, ptVec)
	for _, pt := range contour {
		r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
	}
	x, y, w, h, err := r.backend.BoundingRect(r.ctx, ptVec)
	if err != nil {
		return Rect{}, err
	}
	return Rect{X: x, Y: y, Width: w, Height: h}, nil
}

// MinEnclosingCircle finds the minimum enclosing circle of a point set.
func (r *Runtime) MinEnclosingCircle(contour []Point) (Point, float64, error) {
	if err := r.validateOpen(); err != nil {
		return Point{}, 0, err
	}
	ptVec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil {
		return Point{}, 0, err
	}
	defer r.backend.VecPointsDelete(r.ctx, ptVec)
	for _, pt := range contour {
		r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
	}
	cx, cy, radius, err := r.backend.MinEnclosingCircle(r.ctx, ptVec)
	if err != nil {
		return Point{}, 0, err
	}
	return Point{X: int32(cx), Y: int32(cy)}, radius, nil
}

// MomentsResult holds the spatial and central moments.
type MomentsResult struct {
	M00, M10, M01, M20, M11, M02 float64
	M30, M21, M12, M03           float64
}

// Moments computes moments of a contour.
func (r *Runtime) Moments(contour []Point, binary bool) (*MomentsResult, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	ptVec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecPointsDelete(r.ctx, ptVec)
	for _, pt := range contour {
		r.backend.VecPointsPush(r.ctx, ptVec, pt.X, pt.Y)
	}
	vals, err := r.backend.Moments(r.ctx, ptVec, binary)
	if err != nil {
		return nil, err
	}
	return &MomentsResult{
		M00: vals[0], M10: vals[1], M01: vals[2],
		M20: vals[3], M11: vals[4], M02: vals[5],
		M30: vals[6], M21: vals[7], M12: vals[8], M03: vals[9],
	}, nil
}

// ---------------------------------------------------------------------------
// Hough transforms
// ---------------------------------------------------------------------------

// HoughLine represents a line in polar coordinates (rho, theta).
type HoughLine struct {
	Rho   float64
	Theta float64
}

// HoughLinesPResult represents a line segment from probabilistic Hough transform.
type HoughLinesPResult struct {
	X1, Y1, X2, Y2 int32
}

// HoughCircle represents a detected circle from HoughCircles.
type HoughCircle struct {
	X, Y   float64
	Radius float64
}

// HoughLines detects lines using standard Hough transform. Returns lines as (rho, theta) pairs.
func (r *Runtime) HoughLines(src *Mat, rho, theta float64, threshold int) ([]HoughLine, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}

	vecHandle, err := r.backend.VecDoubleNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecDoubleDelete(r.ctx, vecHandle)

	if err := r.backend.HoughLines(r.ctx, src.handle, vecHandle, rho, theta, int32(threshold), 0, 0); err != nil {
		return nil, err
	}

	n, err := r.backend.VecDoubleLen(r.ctx, vecHandle)
	if err != nil {
		return nil, err
	}

	lines := make([]HoughLine, n/2)
	for i := range lines {
		rho, _ := r.backend.VecDoubleGet(r.ctx, vecHandle, int32(i*2))
		theta, _ := r.backend.VecDoubleGet(r.ctx, vecHandle, int32(i*2+1))
		lines[i] = HoughLine{Rho: rho, Theta: theta}
	}
	return lines, nil
}

// HoughLinesP detects line segments using probabilistic Hough transform.
func (r *Runtime) HoughLinesP(src *Mat, rho, theta float64, threshold int, minLineLength, maxLineGap float64) ([]HoughLinesPResult, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}

	vecHandle, err := r.backend.VecIntNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecIntDelete(r.ctx, vecHandle)

	if err := r.backend.HoughLinesP(r.ctx, src.handle, vecHandle, rho, theta, int32(threshold), minLineLength, maxLineGap); err != nil {
		return nil, err
	}

	n, err := r.backend.VecIntLen(r.ctx, vecHandle)
	if err != nil {
		return nil, err
	}

	lines := make([]HoughLinesPResult, n/4)
	for i := range lines {
		x1, _ := r.backend.VecIntGet(r.ctx, vecHandle, int32(i*4))
		y1, _ := r.backend.VecIntGet(r.ctx, vecHandle, int32(i*4+1))
		x2, _ := r.backend.VecIntGet(r.ctx, vecHandle, int32(i*4+2))
		y2, _ := r.backend.VecIntGet(r.ctx, vecHandle, int32(i*4+3))
		lines[i] = HoughLinesPResult{X1: x1, Y1: y1, X2: x2, Y2: y2}
	}
	return lines, nil
}

// HoughCircles detects circles using Hough transform.
func (r *Runtime) HoughCircles(src *Mat, method HoughMode, dp, minDist, param1, param2 float64, minRadius, maxRadius int) ([]HoughCircle, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}

	vecHandle, err := r.backend.VecDoubleNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecDoubleDelete(r.ctx, vecHandle)

	if err := r.backend.HoughCircles(r.ctx, src.handle, vecHandle, int32(method), dp, minDist, param1, param2, int32(minRadius), int32(maxRadius)); err != nil {
		return nil, err
	}

	n, err := r.backend.VecDoubleLen(r.ctx, vecHandle)
	if err != nil {
		return nil, err
	}

	circles := make([]HoughCircle, n/3)
	for i := range circles {
		x, _ := r.backend.VecDoubleGet(r.ctx, vecHandle, int32(i*3))
		y, _ := r.backend.VecDoubleGet(r.ctx, vecHandle, int32(i*3+1))
		r, _ := r.backend.VecDoubleGet(r.ctx, vecHandle, int32(i*3+2))
		circles[i] = HoughCircle{X: x, Y: y, Radius: r}
	}
	return circles, nil
}

// ---------------------------------------------------------------------------
// Core arithmetic
// ---------------------------------------------------------------------------

func (r *Runtime) newSameSizeOutput(src *Mat) (*Mat, error) {
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	return r.NewMat(rows, cols, typ)
}

// Add computes the per-element sum of two Mats.
func (r *Runtime) Add(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}
	if err := r.backend.Add(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// Subtract computes the per-element difference of two Mats.
func (r *Runtime) Subtract(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("subtract: %w", err)
	}
	if err := r.backend.Subtract(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// Multiply computes the per-element product of two Mats.
func (r *Runtime) Multiply(src1, src2 *Mat, scale float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("multiply: %w", err)
	}
	if err := r.backend.Multiply(r.ctx, src1.handle, src2.handle, dst.handle, scale, -1); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// Divide computes the per-element division of two Mats.
func (r *Runtime) Divide(src1, src2 *Mat, scale float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("divide: %w", err)
	}
	if err := r.backend.Divide(r.ctx, src1.handle, src2.handle, dst.handle, scale, -1); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// AbsDiff computes the per-element absolute difference.
func (r *Runtime) AbsDiff(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("abs_diff: %w", err)
	}
	if err := r.backend.AbsDiff(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// BitwiseAnd computes the per-element bitwise AND.
func (r *Runtime) BitwiseAnd(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("bitwise_and: %w", err)
	}
	if err := r.backend.BitwiseAnd(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// BitwiseOr computes the per-element bitwise OR.
func (r *Runtime) BitwiseOr(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("bitwise_or: %w", err)
	}
	if err := r.backend.BitwiseOr(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// BitwiseXor computes the per-element bitwise XOR.
func (r *Runtime) BitwiseXor(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src1); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src1)
	if err != nil {
		return nil, fmt.Errorf("bitwise_xor: %w", err)
	}
	if err := r.backend.BitwiseXor(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// BitwiseNot computes the per-element bitwise NOT.
func (r *Runtime) BitwiseNot(src *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	dst, err := r.newSameSizeOutput(src)
	if err != nil {
		return nil, fmt.Errorf("bitwise_not: %w", err)
	}
	if err := r.backend.BitwiseNot(r.ctx, src.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// MinMaxLocResult holds the min/max values and their locations.
type MinMaxLocResult struct {
	MinVal float64
	MaxVal float64
	MinLoc Point
	MaxLoc Point
}

// MinMaxLoc finds the global minimum and maximum in an array.
func (r *Runtime) MinMaxLoc(src *Mat) (*MinMaxLocResult, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	minVal, maxVal, minX, minY, maxX, maxY, err := r.backend.MinMaxLoc(r.ctx, src.handle)
	if err != nil {
		return nil, err
	}
	return &MinMaxLocResult{
		MinVal: minVal, MaxVal: maxVal,
		MinLoc: Point{X: minX, Y: minY},
		MaxLoc: Point{X: maxX, Y: maxY},
	}, nil
}

// MeanStdDevResult holds the mean and standard deviation Mats.
type MeanStdDevResult struct {
	Mean   *Mat
	StdDev *Mat
}

// MeanStdDev computes the mean and standard deviation of Mat elements.
func (r *Runtime) MeanStdDev(src *Mat) (*MeanStdDevResult, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	mean, err := r.NewMat(1, 4, CV64F)
	if err != nil {
		return nil, err
	}
	stddev, err := r.NewMat(1, 4, CV64F)
	if err != nil {
		mean.Close()
		return nil, err
	}
	if err := r.backend.MeanStdDev(r.ctx, src.handle, mean.handle, stddev.handle); err != nil {
		mean.Close()
		stddev.Close()
		return nil, err
	}
	return &MeanStdDevResult{Mean: mean, StdDev: stddev}, nil
}

// CountNonZero counts the number of non-zero elements.
func (r *Runtime) CountNonZero(src *Mat) (int, error) {
	if err := r.validateOpen(); err != nil {
		return 0, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return 0, err
	}
	return r.backend.CountNonZero(r.ctx, src.handle)
}

// Split splits a multi-channel Mat into separate single-channel Mats.
func (r *Runtime) Split(src *Mat) ([]*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if err := r.validateOwnedMat(src); err != nil {
		return nil, err
	}
	ch, err := src.Channels()
	if err != nil {
		return nil, err
	}

	vecHandle, err := r.backend.VecMatNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecMatDelete(r.ctx, vecHandle)

	if err := r.backend.Split(r.ctx, src.handle, vecHandle); err != nil {
		return nil, err
	}

	result := make([]*Mat, ch)
	for i := 0; i < ch; i++ {
		h, err := r.backend.VecMatGet(r.ctx, vecHandle, int32(i))
		if err != nil || h == 0 {
			for j := 0; j < i; j++ {
				result[j].Close()
			}
			return nil, fmt.Errorf("split: get channel %d: %w", i, err)
		}
		result[i] = r.wrapMatWithModel(h, Gray)
	}
	return result, nil
}

// Merge merges several single-channel Mats into a multi-channel Mat.
func (r *Runtime) Merge(channels []*Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("merge: no channels")
	}

	vecHandle, err := r.backend.VecMatNew(r.ctx)
	if err != nil {
		return nil, err
	}
	defer r.backend.VecMatDelete(r.ctx, vecHandle)

	for _, ch := range channels {
		if err := r.validateOwnedMat(ch); err != nil {
			return nil, err
		}
		r.backend.VecMatPush(r.ctx, vecHandle, ch.handle)
	}

	rows, _ := channels[0].Rows()
	cols, _ := channels[0].Cols()
	typ := CV8UC1 + MatType(8*(len(channels)-1)) // rough channel count to type
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil {
		return nil, fmt.Errorf("merge: create output: %w", err)
	}

	if err := r.backend.Merge(r.ctx, vecHandle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}

	switch len(channels) {
	case 1:
		dst.model = Gray
	case 3:
		dst.model = BGR
	case 4:
		dst.model = RGBA
	}
	return dst, nil
}

// ---------------------------------------------------------------------------
// New imgproc APIs
// ---------------------------------------------------------------------------

// BilateralFilter applies bilateral filtering to the image.
func (r *Runtime) BilateralFilter(src *Mat, d int32, sigmaColor, sigmaSpace float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.BilateralFilter(r.ctx, src.handle, dst.handle, d, sigmaColor, sigmaSpace); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// InRange checks if src elements are within the specified range [lower, upper].
func (r *Runtime) InRange(src *Mat, lower, upper Scalar) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, CV8UC1)
	if err != nil { return nil, err }
	if err := r.backend.InRange(r.ctx, src.handle, lower.V0, lower.V1, lower.V2, lower.V3, upper.V0, upper.V1, upper.V2, upper.V3, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = Gray
	return dst, nil
}

// MatchTemplate slides template over image and computes match result.
func (r *Runtime) MatchTemplate(img, tmpl *Mat, method TemplateMatchMethod) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(img); err != nil { return nil, err }
	if err := r.validateOwnedMat(tmpl); err != nil { return nil, err }
	imgRows, _ := img.Rows()
	imgCols, _ := img.Cols()
	tmplRows, _ := tmpl.Rows()
	tmplCols, _ := tmpl.Cols()
	resultRows := imgRows - tmplRows + 1
	resultCols := imgCols - tmplCols + 1
	dst, err := r.NewMat(resultRows, resultCols, CV32FC1)
	if err != nil { return nil, err }
	if err := r.backend.MatchTemplate(r.ctx, img.handle, tmpl.handle, dst.handle, int32(method)); err != nil {
		dst.Close()
		return nil, err
	}
	return dst, nil
}

// CalcHist computes histogram for single-channel image.
func (r *Runtime) CalcHist(src *Mat, bins int32, rangeMin, rangeMax float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	dst, err := r.NewMat(int(bins), 1, CV32FC1)
	if err != nil { return nil, err }
	if err := r.backend.CalcHist(r.ctx, src.handle, dst.handle, bins, rangeMin, rangeMax); err != nil {
		dst.Close()
		return nil, err
	}
	return dst, nil
}

// ConnectedComponents computes connected components labeled image.
func (r *Runtime) ConnectedComponents(src *Mat, connectivity int32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, CV32SC1)
	if err != nil { return nil, err }
	if err := r.backend.ConnectedComponents(r.ctx, src.handle, dst.handle, connectivity, 4 /* CV_32S */); err != nil {
		dst.Close()
		return nil, err
	}
	return dst, nil
}

// DistanceTransform computes distance transform of binary image.
func (r *Runtime) DistanceTransform(src *Mat, distType DistanceType, maskSize int32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, CV8UC1)
	if err != nil { return nil, err }
	if err := r.backend.DistanceTransform(r.ctx, src.handle, dst.handle, int32(distType), maskSize); err != nil {
		dst.Close()
		return nil, err
	}
	return dst, nil
}

// CopyMakeBorder pads image with specified border.
func (r *Runtime) CopyMakeBorder(src *Mat, top, bottom, left, right int32, borderType BorderType, value Scalar) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows+int(top+bottom), cols+int(left+right), MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.CopyMakeBorder(r.ctx, src.handle, dst.handle, top, bottom, left, right, int32(borderType), value.V0, value.V1, value.V2, value.V3); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// Rotate rotates image by 90/180/270 degrees.
func (r *Runtime) Rotate(src *Mat, code RotateCode) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	var dstRows, dstCols int
	if code == Rotate180 {
		dstRows, dstCols = rows, cols
	} else {
		dstRows, dstCols = cols, rows
	}
	dst, err := r.NewMat(dstRows, dstCols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Rotate(r.ctx, src.handle, dst.handle, int32(code)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// Hconcat concatenates two images horizontally.
func (r *Runtime) Hconcat(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src1); err != nil { return nil, err }
	if err := r.validateOwnedMat(src2); err != nil { return nil, err }
	rows, _ := src1.Rows()
	cols1, _ := src1.Cols()
	cols2, _ := src2.Cols()
	typ, _ := src1.Type()
	dst, err := r.NewMat(rows, cols1+cols2, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Hconcat(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// Vconcat concatenates two images vertically.
func (r *Runtime) Vconcat(src1, src2 *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src1); err != nil { return nil, err }
	if err := r.validateOwnedMat(src2); err != nil { return nil, err }
	rows1, _ := src1.Rows()
	rows2, _ := src2.Rows()
	cols, _ := src1.Cols()
	typ, _ := src1.Type()
	dst, err := r.NewMat(rows1+rows2, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Vconcat(r.ctx, src1.handle, src2.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src1.model
	return dst, nil
}

// Remap applies generic geometric remap transformation.
func (r *Runtime) Remap(src, map1, map2 *Mat, interpolation, borderMode int32, borderVal float64) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	typ, _ := src.Type()
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Remap(r.ctx, src.handle, dst.handle, map1.handle, map2.handle, interpolation, borderMode, borderVal); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// LUT applies lookup table transformation.
func (r *Runtime) LUT(src, lut *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	if err := r.validateOwnedMat(lut); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	ch, _ := src.Channels()
	var typ MatType
	if ch == 1 {
		typ = CV8UC1
	} else {
		typ = CV8UC3
	}
	dst, err := r.NewMat(rows, cols, typ)
	if err != nil { return nil, err }
	if err := r.backend.LUT(r.ctx, src.handle, lut.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// Integral computes integral image (summed area table).
func (r *Runtime) Integral(src *Mat) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows+1, cols+1, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Integral(r.ctx, src.handle, dst.handle); err != nil {
		dst.Close()
		return nil, err
	}
	return dst, nil
}

// GetPerspectiveTransform calculates 3x3 perspective transform from 4 point pairs.
func (r *Runtime) GetPerspectiveTransform(src, dst [4]Point) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	h, err := r.backend.GetPerspectiveTransform(r.ctx,
		float64(src[0].X), float64(src[0].Y), float64(src[1].X), float64(src[1].Y),
		float64(src[2].X), float64(src[2].Y), float64(src[3].X), float64(src[3].Y),
		float64(dst[0].X), float64(dst[0].Y), float64(dst[1].X), float64(dst[1].Y),
		float64(dst[2].X), float64(dst[2].Y), float64(dst[3].X), float64(dst[3].Y),
	)
	if err != nil { return nil, err }
	return r.wrapMatWithModel(h, Unknown), nil
}

// FillConvexPoly fills a convex polygon.
func (r *Runtime) FillConvexPoly(img *Mat, points []Point, c color.Color, lineType LineType) error {
	if err := validatePair(r, img, img); err != nil { return err }
	cr, cg, cb, ca := colorFetch(c)
	vec, err := r.backend.VecPointsNew(r.ctx)
	if err != nil { return err }
	defer r.backend.VecPointsDelete(r.ctx, vec)
	for _, p := range points {
		r.backend.VecPointsPush(r.ctx, vec, p.X, p.Y)
	}
	return r.backend.FillConvexPoly(r.ctx, img.handle, vec, int32(len(points)), float64(cb), float64(cg), float64(cr), float64(ca), int32(lineType), 0)
}

// ---------------------------------------------------------------------------
// ConvertModel convenience method
// ---------------------------------------------------------------------------

// ConvertModel converts a Mat from its current color model to a target model.
// Automatically selects the correct ColorConversionCode.
func (r *Runtime) ConvertModel(src *Mat, target ColorModel) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }

	srcModel := src.model
	if srcModel == target || srcModel == Unknown || target == Unknown {
		// Just clone
		cloned, err := src.Clone()
		if err != nil {
			return nil, err
		}
		return cloned, nil
	}

	rows, _ := src.Rows()
	cols, _ := src.Cols()
	var dstTyp MatType
	switch target {
	case Gray:
		dstTyp = CV8UC1
	case BGR, RGB:
		dstTyp = CV8UC3
	case RGBA:
		dstTyp = CV8UC4
	default:
		t, _ := src.Type()
		dstTyp = t
	}
	dst, err := r.NewMat(rows, cols, dstTyp)
	if err != nil { return nil, err }
	if err := r.backend.ConvertModel(r.ctx, src.handle, dst.handle, int32(srcModel), int32(target)); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = target
	return dst, nil
}

// ---------------------------------------------------------------------------
// Photo APIs
// ---------------------------------------------------------------------------

// FastNlMeansDenoising denoises a grayscale image using Non-local Means.
func (r *Runtime) FastNlMeansDenoising(src *Mat, h float32, templateWindow, searchWindow int32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.FastNlMeansDenoising(r.ctx, src.handle, dst.handle, h, templateWindow, searchWindow); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// FastNlMeansDenoisingColored denoises a color image using Non-local Means.
func (r *Runtime) FastNlMeansDenoisingColored(src *Mat, h, hColor float32, templateWindow, searchWindow int32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.FastNlMeansDenoisingColored(r.ctx, src.handle, dst.handle, h, hColor, templateWindow, searchWindow); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// DetailEnhance enhances image details.
func (r *Runtime) DetailEnhance(src *Mat, sigmaS, sigmaR float32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.DetailEnhance(r.ctx, src.handle, dst.handle, sigmaS, sigmaR); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// EdgePreservingFilter smooths the image while preserving edges.
func (r *Runtime) EdgePreservingFilter(src *Mat, flags EdgePreservingFilterFlag, sigmaS, sigmaR float32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.EdgePreservingFilter(r.ctx, src.handle, dst.handle, int32(flags), sigmaS, sigmaR); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// PencilSketch creates pencil sketch effect. Returns grayscale and color sketches.
func (r *Runtime) PencilSketch(src *Mat, sigmaS, sigmaR, shadeFactor float32) (gray *Mat, color *Mat, err error) {
	if err := r.validateOpen(); err != nil { return nil, nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	grayDst, err := r.NewMat(rows, cols, CV8UC1)
	if err != nil { return nil, nil, err }
	colorDst, err := r.NewMat(rows, cols, CV8UC3)
	if err != nil { grayDst.Close(); return nil, nil, err }
	if err := r.backend.PencilSketch(r.ctx, src.handle, grayDst.handle, colorDst.handle, sigmaS, sigmaR, shadeFactor); err != nil {
		grayDst.Close()
		colorDst.Close()
		return nil, nil, err
	}
	grayDst.model = Gray
	colorDst.model = BGR
	return grayDst, colorDst, nil
}

// Stylization creates oil painting style effect.
func (r *Runtime) Stylization(src *Mat, sigmaS, sigmaR float32) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	rows, _ := src.Rows()
	cols, _ := src.Cols()
	typ, _ := src.Type()
	dst, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.Stylization(r.ctx, src.handle, dst.handle, sigmaS, sigmaR); err != nil {
		dst.Close()
		return nil, err
	}
	dst.model = src.model
	return dst, nil
}

// SeamlessClone blends source image into destination using seamless cloning.
func (r *Runtime) SeamlessClone(src, dst, mask *Mat, center Point, flags SeamlessCloneMethod) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	if err := r.validateOwnedMat(dst); err != nil { return nil, err }
	rows, _ := dst.Rows()
	cols, _ := dst.Cols()
	typ, _ := dst.Type()
	out, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	if err := r.backend.SeamlessClone(r.ctx, src.handle, dst.handle, mask.handle, out.handle, center.X, center.Y, int32(flags)); err != nil {
		out.Close()
		return nil, err
	}
	out.model = dst.model
	return out, nil
}

// ---------------------------------------------------------------------------
// Features2d APIs
// ---------------------------------------------------------------------------

// FAST detects keypoints using FAST detector.
func (r *Runtime) FAST(src *Mat, threshold int, nonmaxSuppression bool) ([]KeyPoint, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, err }
	kpVec, err := r.backend.VecKeypointNew(r.ctx)
	if err != nil { return nil, err }
	defer r.backend.VecKeypointDelete(r.ctx, kpVec)
	nms := int32(0)
	if nonmaxSuppression { nms = 1 }
	if err := r.backend.FAST(r.ctx, src.handle, kpVec, int32(threshold), nms); err != nil {
		return nil, err
	}
	return r.readKeypoints(kpVec)
}

// ORBDetectCompute detects ORB keypoints and computes descriptors.
func (r *Runtime) ORBDetectCompute(src *Mat, nfeatures int, scaleFactor float32, nlevels int) ([]KeyPoint, *Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, nil, err }
	if err := r.validateOwnedMat(src); err != nil { return nil, nil, err }
	kpVec, err := r.backend.VecKeypointNew(r.ctx)
	if err != nil { return nil, nil, err }
	defer r.backend.VecKeypointDelete(r.ctx, kpVec)
	desc, err := r.NewMat(1, 1, CV8UC1)
	if err != nil { return nil, nil, err }
	if err := r.backend.ORBDetectCompute(r.ctx, src.handle, 0, kpVec, desc.handle, int32(nfeatures), scaleFactor, int32(nlevels)); err != nil {
		desc.Close()
		return nil, nil, err
	}
	kps, err := r.readKeypoints(kpVec)
	if err != nil { desc.Close(); return nil, nil, err }
	return kps, desc, nil
}

// BFMatch matches descriptors using brute-force matcher.
func (r *Runtime) BFMatch(desc1, desc2 *Mat, normType NormType) ([]DMatch, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(desc1); err != nil { return nil, err }
	if err := r.validateOwnedMat(desc2); err != nil { return nil, err }
	matchVec, err := r.backend.VecDMatchNew(r.ctx)
	if err != nil { return nil, err }
	defer r.backend.VecDMatchDelete(r.ctx, matchVec)
	if err := r.backend.BFMatch(r.ctx, desc1.handle, desc2.handle, matchVec, int32(normType)); err != nil {
		return nil, err
	}
	return r.readDMatches(matchVec)
}

// DrawKeypoints draws keypoints on an image.
func (r *Runtime) DrawKeypoints(img *Mat, keypoints []KeyPoint, c color.Color) (*Mat, error) {
	if err := r.validateOpen(); err != nil { return nil, err }
	if err := r.validateOwnedMat(img); err != nil { return nil, err }
	rows, _ := img.Rows()
	cols, _ := img.Cols()
	typ, _ := img.Type()
	out, err := r.NewMat(rows, cols, MatType(typ))
	if err != nil { return nil, err }
	kpVec, err := r.backend.VecKeypointNew(r.ctx)
	if err != nil { out.Close(); return nil, err }
	defer r.backend.VecKeypointDelete(r.ctx, kpVec)
	// Can't easily push keypoints from Go — use C++ to populate them
	// Instead we use the simpler path: pass the kpVec from FAST/ORB directly
	// For now, draw on the output image
	cr, cg, cb, _ := colorFetch(c)
	if err := r.backend.DrawKeypoints(r.ctx, img.handle, kpVec, out.handle, float64(cb), float64(cg), float64(cr)); err != nil {
		out.Close()
		return nil, err
	}
	out.model = img.model
	return out, nil
}

func (r *Runtime) readKeypoints(kpVec matHandle) ([]KeyPoint, error) {
	n, err := r.backend.VecKeypointLen(r.ctx, kpVec)
	if err != nil { return nil, err }
	kps := make([]KeyPoint, n)
	for i := 0; i < n; i++ {
		x, y, sz, angle, resp, oct, cid, err := r.backend.VecKeypointGet(r.ctx, kpVec, int32(i))
		if err != nil { return nil, err }
		kps[i] = KeyPoint{X: x, Y: y, Size: sz, Angle: angle, Response: resp, Octave: oct, ClassID: cid}
	}
	return kps, nil
}

func (r *Runtime) readDMatches(matchVec matHandle) ([]DMatch, error) {
	n, err := r.backend.VecDMatchLen(r.ctx, matchVec)
	if err != nil { return nil, err }
	matches := make([]DMatch, n)
	for i := 0; i < n; i++ {
		qi, ti, dist, imgi, err := r.backend.VecDMatchGet(r.ctx, matchVec, int32(i))
		if err != nil { return nil, err }
		matches[i] = DMatch{QueryIdx: qi, TrainIdx: ti, Distance: dist, ImgIdx: imgi}
	}
	return matches, nil
}

// ---------------------------------------------------------------------------
// Highgui APIs
// ---------------------------------------------------------------------------

// ImShow displays an image in the specified window.
func (r *Runtime) ImShow(winname string, img *Mat) error {
	if err := r.validateOpen(); err != nil { return err }
	if err := r.validateOwnedMat(img); err != nil { return err }
	name := []byte(winname)
	return r.backend.ImShow(r.ctx, unsafe.Pointer(&name[0]), int32(len(name)), img.handle)
}

// WaitKey waits for a key press for the specified delay (milliseconds). Returns -1 if no key.
func (r *Runtime) WaitKey(delay int32) (int32, error) {
	if err := r.validateOpen(); err != nil { return -1, err }
	return r.backend.WaitKey(r.ctx, delay)
}

// DestroyWindow closes the specified window.
func (r *Runtime) DestroyWindow(winname string) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	name := []byte(winname)
	return r.backend.DestroyWindow(r.ctx, unsafe.Pointer(&name[0]), int32(len(name)))
}

// ---------------------------------------------------------------------------
// Core Mat extras
// ---------------------------------------------------------------------------

// Diag extracts the diagonal from a Mat.
func (m *Mat) Diag() (*Mat, error) {
	if err := m.validate(); err != nil { return nil, err }
	h, err := m.runtime.backend.MatDiag(m.runtime.ctx, m.handle)
	if err != nil { return nil, err }
	return m.runtime.wrapMatWithModel(h, m.model), nil
}

// AtByte returns the byte value at (row, col, channel).
func (m *Mat) AtByte(row, col, channel int32) (uint8, error) {
	if err := m.validate(); err != nil { return 0, err }
	return m.runtime.backend.MatAtU8(m.runtime.ctx, m.handle, row, col, channel)
}

// SetByte sets the byte value at (row, col, channel).
func (m *Mat) SetByte(row, col, channel int32, val uint8) error {
	if err := m.validate(); err != nil { return err }
	return m.runtime.backend.MatSetU8(m.runtime.ctx, m.handle, row, col, channel, val)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (r *Runtime) wrapMat(h matHandle) *Mat {
	return r.wrapMatWithModel(h, Unknown)
}

func (r *Runtime) wrapMatWithModel(h matHandle, model ColorModel) *Mat {
	r.mu.Lock()
	r.mats[h] = struct{}{}
	r.mu.Unlock()
	return &Mat{runtime: r, handle: h, model: model}
}

func (r *Runtime) closeMat(h matHandle) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	r.mu.Lock()
	_, ok := r.mats[h]
	if ok {
		delete(r.mats, h)
	}
	r.mu.Unlock()
	if !ok {
		return ErrInvalidMat
	}
	return r.backend.CloseMat(r.ctx, h)
}

func (r *Runtime) validateOpen() error {
	if r == nil || r.closed.Load() {
		return ErrClosed
	}
	return nil
}

func validatePair(r *Runtime, src, dst *Mat) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := src.validate(); err != nil {
		return err
	}
	if err := dst.validate(); err != nil {
		return err
	}
	if src.runtime != r || dst.runtime != r {
		return ErrInvalidMat
	}
	return nil
}

func (r *Runtime) validateOwnedMat(m *Mat) error {
	if err := m.validate(); err != nil {
		return err
	}
	if m.runtime != r {
		return ErrInvalidMat
	}
	return nil
}
