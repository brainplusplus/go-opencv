package abi

// Function names define the stable CGO-free ABI exported by goopencv native library.
const (
	Malloc = "goopencv_malloc"
	Free   = "goopencv_free"

	MatNew         = "goopencv_mat_new"
	MatNewFromData = "goopencv_mat_new_from_data"
	MatDelete      = "goopencv_mat_delete"
	MatRows        = "goopencv_mat_rows"
	MatCols        = "goopencv_mat_cols"
	MatType        = "goopencv_mat_type"
	MatClone       = "goopencv_mat_clone"
	MatDataPtr     = "goopencv_mat_data_ptr"
	MatEmpty       = "goopencv_mat_empty"
	MatElemSize    = "goopencv_mat_elem_size"
	MatStep        = "goopencv_mat_step"
	MatChannels    = "goopencv_mat_channels"
	MatRow         = "goopencv_mat_row"
	MatCol         = "goopencv_mat_col"
	MatRegion      = "goopencv_mat_region"
	MatReshape     = "goopencv_mat_reshape"
	MatTotal       = "goopencv_mat_total"
	MatSetTo       = "goopencv_mat_set_to"
	MatConvertTo   = "goopencv_mat_convert_to"
	MatZeros       = "goopencv_mat_zeros"
	MatOnes        = "goopencv_mat_ones"
	MatEye         = "goopencv_mat_eye"

	CvtColor              = "goopencv_imgproc_cvt_color"
	Resize                = "goopencv_imgproc_resize"
	Blur                  = "goopencv_imgproc_blur"
	GaussianBlur          = "goopencv_imgproc_gaussian_blur"
	MedianBlur            = "goopencv_imgproc_median_blur"
	Threshold             = "goopencv_imgproc_threshold"
	AdaptiveThreshold     = "goopencv_imgproc_adaptive_threshold"
	Canny                 = "goopencv_imgproc_canny"
	Flip                  = "goopencv_imgproc_flip"
	Sobel                 = "goopencv_imgproc_sobel"
	Laplacian             = "goopencv_imgproc_laplacian"
	Erode                 = "goopencv_imgproc_erode"
	Dilate                = "goopencv_imgproc_dilate"
	MorphologyEx          = "goopencv_imgproc_morphology_ex"
	GetStructuringElement = "goopencv_imgproc_get_structuring_element"
	EqualizeHist          = "goopencv_imgproc_equalize_hist"
	PutText               = "goopencv_imgproc_put_text"
	Rectangle             = "goopencv_imgproc_rectangle"
	Circle                = "goopencv_imgproc_circle"
	Line                  = "goopencv_imgproc_line"
	FillPoly              = "goopencv_imgproc_fill_poly"
	ArrowedLine           = "goopencv_imgproc_arrowed_line"
	FindContours          = "goopencv_imgproc_find_contours"
	DrawContours          = "goopencv_imgproc_draw_contours"
	ContourArea           = "goopencv_imgproc_contour_area"
	ArcLength             = "goopencv_imgproc_arc_length"
	BoundingRect          = "goopencv_imgproc_bounding_rect"
	MinEnclosingCircle    = "goopencv_imgproc_min_enclosing_circle"
	Moments               = "goopencv_imgproc_moments"
	HoughLines            = "goopencv_imgproc_hough_lines"
	HoughLinesP           = "goopencv_imgproc_hough_lines_p"
	HoughCircles          = "goopencv_imgproc_hough_circles"
	WarpAffine            = "goopencv_imgproc_warp_affine"
	WarpPerspective       = "goopencv_imgproc_warp_perspective"
	GetRotationMatrix2D   = "goopencv_imgproc_get_rotation_matrix2d"
	GetAffineTransform    = "goopencv_imgproc_get_affine_transform"
	Transpose             = "goopencv_imgproc_transpose"

	CalcHist     = "goopencv_imgproc_calc_hist"
	Normalize    = "goopencv_imgproc_normalize"
	BitwiseAnd   = "goopencv_core_bitwise_and"
	BitwiseOr    = "goopencv_core_bitwise_or"
	BitwiseXor   = "goopencv_core_bitwise_xor"
	BitwiseNot   = "goopencv_core_bitwise_not"
	MinMaxLoc    = "goopencv_core_min_max_loc"
	MeanStdDev   = "goopencv_core_mean_std_dev"
	Split        = "goopencv_core_split"
	Merge        = "goopencv_core_merge"
	Add          = "goopencv_core_add"
	Subtract     = "goopencv_core_subtract"
	Multiply     = "goopencv_core_multiply"
	Divide       = "goopencv_core_divide"
	AbsDiff      = "goopencv_core_abs_diff"
	Sqrt         = "goopencv_core_sqrt"
	Max          = "goopencv_core_max"
	Min          = "goopencv_core_min"
	Compare      = "goopencv_core_compare"
	CountNonZero = "goopencv_core_count_non_zero"
	Sum          = "goopencv_core_sum"
	Mean         = "goopencv_core_mean"
	Norm         = "goopencv_core_norm"

	// Vector helpers for contour/keypoint transfer
	VecNewPoints       = "goopencv_vec_points_new"
	VecPushPoint       = "goopencv_vec_points_push"
	VecLenPoints       = "goopencv_vec_points_len"
	VecGetPoint        = "goopencv_vec_points_get"
	VecDeletePoints    = "goopencv_vec_points_delete"
	VecNewVecPoints    = "goopencv_vec_vec_points_new"
	VecPushVecPoints   = "goopencv_vec_vec_points_push"
	VecLenVecPoints    = "goopencv_vec_vec_points_len"
	VecGetVecPoints    = "goopencv_vec_vec_points_get"
	VecDeleteVecPoints = "goopencv_vec_vec_points_delete"

	VecNewDouble    = "goopencv_vec_double_new"
	VecGetDouble    = "goopencv_vec_double_get"
	VecLenDouble    = "goopencv_vec_double_len"
	VecDeleteDouble = "goopencv_vec_double_delete"

	VecNewInt    = "goopencv_vec_int_new"
	VecGetInt    = "goopencv_vec_int_get"
	VecLenInt    = "goopencv_vec_int_len"
	VecDeleteInt = "goopencv_vec_int_delete"

	VecNewMat    = "goopencv_vec_mat_new"
	VecPushMat   = "goopencv_vec_mat_push"
	VecLenMat    = "goopencv_vec_mat_len"
	VecGetMat    = "goopencv_vec_mat_get"
	VecDeleteMat = "goopencv_vec_mat_delete"

	// New imgproc
	BilateralFilter      = "goopencv_imgproc_bilateral_filter"
	InRange              = "goopencv_imgproc_in_range"
	MatchTemplate        = "goopencv_imgproc_match_template"
	CalcHistABI          = "goopencv_imgproc_calc_hist"
	ConnectedComponents  = "goopencv_imgproc_connected_components"
	DistanceTransform    = "goopencv_imgproc_distance_transform"
	CopyMakeBorder       = "goopencv_imgproc_copy_make_border"
	Rotate               = "goopencv_imgproc_rotate"
	Hconcat              = "goopencv_imgproc_hconcat"
	Vconcat              = "goopencv_imgproc_vconcat"
	Remap                = "goopencv_imgproc_remap"
	LUT                  = "goopencv_imgproc_lut"
	Integral             = "goopencv_imgproc_integral"
	GetPerspectiveTransform = "goopencv_imgproc_get_perspective_transform"
	FillConvexPoly       = "goopencv_imgproc_fill_convex_poly"
	ConvertModel         = "goopencv_imgproc_convert_model"

	// Photo
	FastNlMeansDenoising        = "goopencv_photo_fast_nl_means_denoising"
	FastNlMeansDenoisingColored = "goopencv_photo_fast_nl_means_denoising_colored"
	DetailEnhance               = "goopencv_photo_detail_enhance"
	EdgePreservingFilter        = "goopencv_photo_edge_preserving_filter"
	PencilSketch                = "goopencv_photo_pencil_sketch"
	Stylization                 = "goopencv_photo_stylization"
	SeamlessClone               = "goopencv_photo_seamless_clone"

	// Features2d
	FAST             = "goopencv_features2d_fast"
	ORBDetectCompute = "goopencv_features2d_orb_detect_and_compute"
	BFMatch          = "goopencv_features2d_bf_match"
	DrawKeypoints    = "goopencv_features2d_draw_keypoints"

	// Highgui
	ImShow         = "goopencv_highgui_imshow"
	WaitKey        = "goopencv_highgui_wait_key"
	DestroyWindow  = "goopencv_highgui_destroy_window"

	// Core extras
	MatDiag  = "goopencv_mat_diag"
	MatAtU8  = "goopencv_mat_at_u8"
	MatSetU8 = "goopencv_mat_set_u8"

	// Vector helpers — keypoints
	VecNewKeypoint    = "goopencv_vec_keypoint_new"
	VecLenKeypoint    = "goopencv_vec_keypoint_len"
	VecGetKeypoint    = "goopencv_vec_keypoint_get"
	VecDeleteKeypoint = "goopencv_vec_keypoint_delete"

	// Vector helpers — dmatch
	VecNewDMatch    = "goopencv_vec_dmatch_new"
	VecLenDMatch    = "goopencv_vec_dmatch_len"
	VecGetDMatch    = "goopencv_vec_dmatch_get"
	VecDeleteDMatch = "goopencv_vec_dmatch_delete"
)

type Handle uint64

// ErrorCode is returned by ABI calls that cannot use Go errors across FFI boundary.
type ErrorCode int32

const (
	OK ErrorCode = iota
	ErrInvalidArgument
	ErrInvalidHandle
	ErrOutOfMemory
	ErrOpenCV
	ErrUnsupported
)
