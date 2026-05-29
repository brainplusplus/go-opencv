package contract

// Module describes one OpenCV surface area for ABI code generation.
type Module struct {
	Name      string
	Resources []Resource
	Functions []Function
}

type Resource struct {
	Name    string
	Methods []Function
}

type Function struct {
	Name      string
	Symbol    string
	Stability Stability
}

type Stability string

const (
	Stable   Stability = "stable"
	Planned  Stability = "planned"
	Research Stability = "research"
)

var Modules = []Module{
	{
		Name: "mat",
		Resources: []Resource{{
			Name: "Mat",
			Methods: []Function{
				{Name: "Close", Symbol: "goopencv_mat_delete", Stability: Stable},
				{Name: "Clone", Symbol: "goopencv_mat_clone", Stability: Stable},
				{Name: "Rows", Symbol: "goopencv_mat_rows", Stability: Stable},
				{Name: "Cols", Symbol: "goopencv_mat_cols", Stability: Stable},
				{Name: "Type", Symbol: "goopencv_mat_type", Stability: Stable},
				{Name: "Empty", Symbol: "goopencv_mat_empty", Stability: Planned},
				{Name: "ElemSize", Symbol: "goopencv_mat_elem_size", Stability: Planned},
				{Name: "Step", Symbol: "goopencv_mat_step", Stability: Planned},
				{Name: "Row", Symbol: "goopencv_mat_row", Stability: Planned},
				{Name: "Col", Symbol: "goopencv_mat_col", Stability: Planned},
				{Name: "Region", Symbol: "goopencv_mat_region", Stability: Planned},
				{Name: "Reshape", Symbol: "goopencv_mat_reshape", Stability: Planned},
			},
		}},
		Functions: []Function{
			{Name: "NewMat", Symbol: "goopencv_mat_new", Stability: Stable},
			{Name: "MatZeros", Symbol: "goopencv_mat_zeros", Stability: Planned},
			{Name: "MatOnes", Symbol: "goopencv_mat_ones", Stability: Planned},
			{Name: "MatMerge", Symbol: "goopencv_mat_merge", Stability: Planned},
		},
	},
	{
		Name: "imgproc",
		Functions: []Function{
			{Name: "CvtColor", Symbol: "goopencv_imgproc_cvt_color", Stability: Stable},
			{Name: "Resize", Symbol: "goopencv_imgproc_resize", Stability: Stable},
			{Name: "Blur", Symbol: "goopencv_imgproc_blur", Stability: Planned},
			{Name: "GaussianBlur", Symbol: "goopencv_imgproc_gaussian_blur", Stability: Planned},
			{Name: "Threshold", Symbol: "goopencv_imgproc_threshold", Stability: Planned},
			{Name: "AdaptiveThreshold", Symbol: "goopencv_imgproc_adaptive_threshold", Stability: Planned},
			{Name: "Canny", Symbol: "goopencv_imgproc_canny", Stability: Planned},
			{Name: "Rectangle", Symbol: "goopencv_imgproc_rectangle", Stability: Planned},
			{Name: "Circle", Symbol: "goopencv_imgproc_circle", Stability: Planned},
			{Name: "Line", Symbol: "goopencv_imgproc_line", Stability: Planned},
			{Name: "PutText", Symbol: "goopencv_imgproc_put_text", Stability: Planned},
		},
	},
	{Name: "features2d", Functions: []Function{{Name: "DetectAndCompute", Symbol: "goopencv_features2d_detect_and_compute", Stability: Research}}},
	{Name: "objdetect", Resources: []Resource{{Name: "CascadeClassifier", Methods: []Function{{Name: "DetectMultiScale", Symbol: "goopencv_objdetect_cascade_classifier_detect_multi_scale", Stability: Research}}}}},
	{Name: "dnn", Functions: []Function{{Name: "BlobFromImage", Symbol: "goopencv_dnn_blob_from_image", Stability: Research}}},
}
