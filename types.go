package opencv

import "image/color"

// MatType mirrors OpenCV/OpenCV.js Mat depth+channel type constants.
type MatType int32

// ColorModel defines logical channel ordering semantics for Mat/image I/O.
// Unknown means model cannot be trusted after ambiguous/custom operations.
type ColorModel int

const (
	Unknown ColorModel = iota
	BGR
	RGB
	RGBA
	Gray
)

const (
	CV8U  MatType = 0
	CV8S  MatType = 1
	CV16U MatType = 2
	CV16S MatType = 3
	CV32S MatType = 4
	CV32F MatType = 5
	CV64F MatType = 6
	CV16F MatType = 7

	CV8UC1 MatType = CV8U
	CV8UC2 MatType = CV8U + 8
	CV8UC3 MatType = CV8U + 16
	CV8UC4 MatType = CV8U + 24

	CV32FC1 MatType = CV32F + 0
	CV32FC2 MatType = CV32F + 8
	CV32FC3 MatType = CV32F + 16
	CV32FC4 MatType = CV32F + 24

	CV32SC1 MatType = CV32S + 0
	CV32SC2 MatType = CV32S + 8
	CV32SC3 MatType = CV32S + 16
	CV32SC4 MatType = CV32S + 24
)

// ColorConversionCode mirrors OpenCV.js color conversion constants.
type ColorConversionCode int32

const (
	ColorBGR2BGRA  ColorConversionCode = 0
	ColorRGB2RGBA  ColorConversionCode = ColorBGR2BGRA
	ColorBGRA2BGR  ColorConversionCode = 1
	ColorRGBA2RGB  ColorConversionCode = ColorBGRA2BGR
	ColorBGR2RGBA  ColorConversionCode = 2
	ColorRGB2BGRA  ColorConversionCode = ColorBGR2RGBA
	ColorRGBA2BGR  ColorConversionCode = 3
	ColorBGRA2RGB  ColorConversionCode = ColorRGBA2BGR
	ColorBGR2RGB   ColorConversionCode = 4
	ColorRGB2BGR   ColorConversionCode = ColorBGR2RGB
	ColorBGRA2RGBA ColorConversionCode = 5
	ColorRGBA2BGRA ColorConversionCode = ColorBGRA2RGBA
	ColorBGR2Gray  ColorConversionCode = 6
	ColorRGB2Gray  ColorConversionCode = 7
	ColorGray2BGR  ColorConversionCode = 8
	ColorGray2RGB  ColorConversionCode = ColorGray2BGR
	ColorGray2BGRA ColorConversionCode = 9
	ColorGray2RGBA ColorConversionCode = ColorGray2BGRA
	ColorBGRA2Gray ColorConversionCode = 10
	ColorRGBA2Gray ColorConversionCode = 11
)

type IMReadFlag int32

const (
	IMReadUnchanged IMReadFlag = -1
	IMReadGrayScale IMReadFlag = 0
	IMReadColor     IMReadFlag = 1
)

type Point struct {
	X int32
	Y int32
}

type Size struct {
	Width  int32
	Height int32
}

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

type Scalar struct {
	V0 float64
	V1 float64
	V2 float64
	V3 float64
}

type RGBAColor struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// RGBA implements color.Color — standard Go image color, R/G/B/A order.
func (c RGBAColor) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// BGR follows OpenCV native channel ordering (Blue, Green, Red, Alpha).
// Useful for direct parity with OpenCV C++/Python Scalar(B, G, R, A) usage.
type BGRColor struct {
	B uint8
	G uint8
	R uint8
	A uint8
}

// RGBA implements color.Color — maps BGR field order to standard R,G,B output.
func (c BGRColor) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// colorFromFetch extracts uint8 R,G,B,A from any color.Color.
// Values are in standard R,G,B order — the ABI layer swaps to BGR for OpenCV.
func colorFetch(c color.Color) (r, g, b, a uint8) {
	r32, g32, b32, a32 := c.RGBA()
	// color.Color RGBA() returns 0-65535 premultiplied values
	// Convert back to 0-255 range
	return uint8(r32 >> 8), uint8(g32 >> 8), uint8(b32 >> 8), uint8(a32 >> 8)
}

type BorderType int32

const (
	BorderConstant BorderType = iota
	BorderReplicate
	BorderReflect
	BorderWrap
	BorderReflect101
	BorderTransparent
	BorderDefault
	BorderIsolated
)

type AdaptiveThresholdType int32

const (
	AdaptiveThresholdMean AdaptiveThresholdType = iota
	AdaptiveThresholdGaussian
)

type ThresholdType int32

const (
	ThresholdBinary ThresholdType = iota
	ThresholdBinaryInv
	ThresholdTrunc
	ThresholdToZero
	ThresholdToZeroInv
	ThresholdMask
	ThresholdOtsu
	ThresholdTriangle
)

type InterpolationType int32

const (
	InterpolationNearest InterpolationType = iota
	InterpolationLinear
	InterpolationCubic
	InterpolationArea
	InterpolationLanczos4
)

type MorphShape int32

const (
	MorphRect MorphShape = iota
	MorphCross
	MorphEllipse
)

type HersheyFontType int32

const (
	FontHersheySimplex HersheyFontType = iota
	FontHersheyPlain
	FontHersheyDuplex
	FontHersheyComplex
	FontHersheyTriplex
	FontHersheyComplexSmall
	FontHersheyScriptSimplex
	FontHersheyScriptComplex
)

type FlipCode int32

const (
	FlipHorizontal FlipCode = 1
	FlipVertical   FlipCode = 0
	FlipBoth       FlipCode = -1
)

type MorphType int32

const (
	MorphErode    MorphType = 0
	MorphDilate   MorphType = 1
	MorphOpen     MorphType = 2
	MorphClose    MorphType = 3
	MorphGradient MorphType = 4
	MorphTophat   MorphType = 5
	MorphBlackhat MorphType = 6
	MorphHitmiss  MorphType = 7
)

type LineType int32

const (
	Line8  LineType = 0
	Line4  LineType = 1
	LineAA LineType = 16
)

type RetrievalMode int32

const (
	RetrievalExternal  RetrievalMode = 0
	RetrievalList      RetrievalMode = 1
	RetrievalCComp     RetrievalMode = 2
	RetrievalTree      RetrievalMode = 3
	RetrievalFloodfill RetrievalMode = 4
)

type ContourApproximationMode int32

const (
	ChainApproxNone     ContourApproximationMode = 1
	ChainApproxSimple   ContourApproximationMode = 2
	ChainApproxTC89L1   ContourApproximationMode = 3
	ChainApproxTC89KCOS ContourApproximationMode = 4
)

type NormType int32

const (
	NormInf      NormType = 1
	NormL1       NormType = 2
	NormL2       NormType = 4
	NormL2Sqr    NormType = 5
	NormHamming  NormType = 6
	NormHamming2 NormType = 7
	NormRelative NormType = 8
	NormMinMax   NormType = 32
)

type HoughMode int32

const (
	HoughStandard      HoughMode = 0
	HoughProbabilistic HoughMode = 1
	HoughMultiScale    HoughMode = 2
	HoughGradient      HoughMode = 3
)

type KeyPoint struct {
	X        float32
	Y        float32
	Size     float32
	Angle    float32
	Response float32
	Octave   int32
	ClassID  int32
}

type DMatch struct {
	QueryIdx int32
	TrainIdx int32
	ImgIdx   int32
	Distance float32
}

type RotateCode int32

const (
	Rotate90Clockwise        RotateCode = 0
	Rotate180                RotateCode = 1
	Rotate90CounterClockwise RotateCode = 2
)

type TemplateMatchMethod int32

const (
	TMSqDiff       TemplateMatchMethod = 0
	TMSqDiffNormed TemplateMatchMethod = 1
	TMCCorr        TemplateMatchMethod = 2
	TMCCorrNormed  TemplateMatchMethod = 3
	TMCCoeff       TemplateMatchMethod = 4
	TMCCoeffNormed TemplateMatchMethod = 5
)

type DistanceType int32

const (
	DistUser    DistanceType = 0
	DistL1      DistanceType = 1
	DistL2      DistanceType = 2
	DistC       DistanceType = 3
	DistL12     DistanceType = 4
	DistFair    DistanceType = 5
	DistWelsch  DistanceType = 6
	DistHuber   DistanceType = 7
	DistMask    DistanceType = 8
)

type EdgePreservingFilterFlag int32

const (
	RecursFilter EdgePreservingFilterFlag = 1
	NormconvFilter EdgePreservingFilterFlag = 2
)

type SeamlessCloneMethod int32

const (
	NormalClone    SeamlessCloneMethod = 1
	MixedClone     SeamlessCloneMethod = 2
	FeatureExchange SeamlessCloneMethod = 3
)

type ORBScoreType int32

const (
	ORBHarrisScore ORBScoreType = 0
	ORBFASTScore   ORBScoreType = 1
)
