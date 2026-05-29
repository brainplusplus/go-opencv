// goopencv_abi.cpp — Real OpenCV flat ABI wrapper
// Compile via build-tools/build-goopencv.bat (MSVC x64, static link opencv-mobile)

#include <stdint.h>
#include <opencv2/core.hpp>
#include <opencv2/imgproc.hpp>

// DLL export annotation (Windows only; extern "C" suffices on Unix).
#ifdef _WIN32
#define ABI_EXPORT __declspec(dllexport)
#else
#define ABI_EXPORT
#endif

// Error codes matching internal/abi/spec.go
constexpr int32_t OK = 0;
constexpr int32_t ERR_INVALID_ARGUMENT = 1;
constexpr int32_t ERR_INVALID_HANDLE = 2;
constexpr int32_t ERR_OUT_OF_MEMORY = 3;
constexpr int32_t ERR_OPENCV = 4;
constexpr int32_t ERR_UNSUPPORTED = 5;

// ---------------------------------------------------------------------------
// Mat lifecycle
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_mat_new(int32_t rows, int32_t cols, int32_t type) {
    if (rows <= 0 || cols <= 0) return 0;
    try {
        auto* m = new cv::Mat(rows, cols, type);
        return (uint64_t)(intptr_t)m;
    } catch (std::exception&) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_delete(uint64_t handle) {
    if (handle == 0) return ERR_INVALID_HANDLE;
    try {
        delete (cv::Mat*)(intptr_t)handle;
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_rows(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return ((cv::Mat*)(intptr_t)handle)->rows;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_cols(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return ((cv::Mat*)(intptr_t)handle)->cols;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_type(uint64_t handle) {
    if (handle == 0) return ERR_INVALID_HANDLE;
    try {
        return ((cv::Mat*)(intptr_t)handle)->type();
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_clone(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        auto* src = (cv::Mat*)(intptr_t)handle;
        auto* dst = new cv::Mat(src->clone());
        return (uint64_t)(intptr_t)dst;
    } catch (std::bad_alloc&) {
        return 0;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_new_from_data(
    const uint8_t* data, int32_t rows, int32_t cols, int32_t type
) {
    if (!data || rows <= 0 || cols <= 0) return 0;
    try {
        auto* m = new cv::Mat(rows, cols, type);
        size_t total = (size_t)rows * (size_t)cols * m->elemSize();
        memcpy(m->data, data, total);
        return (uint64_t)(intptr_t)m;
    } catch (std::bad_alloc&) {
        return 0;
    } catch (...) {
        return 0;
    }
}

// ---------------------------------------------------------------------------
// Mat extras
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_mat_empty(uint64_t handle) {
    if (handle == 0) return 1;
    try {
        return ((cv::Mat*)(intptr_t)handle)->empty() ? 1 : 0;
    } catch (...) {
        return 1;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_elem_size(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return (int32_t)((cv::Mat*)(intptr_t)handle)->elemSize();
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_data_ptr(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        auto* m = (cv::Mat*)(intptr_t)handle;
        return (uint64_t)(intptr_t)m->data;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_step(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return (int32_t)((cv::Mat*)(intptr_t)handle)->step;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_channels(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return ((cv::Mat*)(intptr_t)handle)->channels();
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_total(uint64_t handle) {
    if (handle == 0) return 0;
    try {
        return (int32_t)((cv::Mat*)(intptr_t)handle)->total();
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_row(uint64_t handle, int32_t row) {
    if (handle == 0) return 0;
    try {
        auto* src = (cv::Mat*)(intptr_t)handle;
        auto* dst = new cv::Mat(src->row(row));
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_col(uint64_t handle, int32_t col) {
    if (handle == 0) return 0;
    try {
        auto* src = (cv::Mat*)(intptr_t)handle;
        auto* dst = new cv::Mat(src->col(col));
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_region(uint64_t handle, int32_t x, int32_t y, int32_t w, int32_t h) {
    if (handle == 0) return 0;
    try {
        auto* src = (cv::Mat*)(intptr_t)handle;
        auto* dst = new cv::Mat((*src)(cv::Rect(x, y, w, h)).clone());
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_reshape(uint64_t handle, int32_t cn, int32_t rows) {
    if (handle == 0) return 0;
    try {
        auto* src = (cv::Mat*)(intptr_t)handle;
        auto* dst = new cv::Mat(src->reshape(cn, rows));
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_set_to(uint64_t handle, double v0, double v1, double v2, double v3) {
    if (handle == 0) return ERR_INVALID_HANDLE;
    try {
        ((cv::Mat*)(intptr_t)handle)->setTo(cv::Scalar(v0, v1, v2, v3));
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_mat_convert_to(uint64_t src_handle, uint64_t dst_handle, int32_t rtype, double alpha, double beta) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        src->convertTo(*dst, rtype, alpha, beta);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_zeros(int32_t rows, int32_t cols, int32_t type) {
    try {
        auto m = cv::Mat::zeros(rows, cols, type);
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_ones(int32_t rows, int32_t cols, int32_t type) {
    try {
        auto m = cv::Mat::ones(rows, cols, type);
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_mat_eye(int32_t rows, int32_t cols, int32_t type) {
    try {
        auto m = cv::Mat::eye(rows, cols, type);
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: filtering
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_cvt_color(uint64_t src_handle, uint64_t dst_handle, int32_t code) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::cvtColor(*src, *dst, code);
        return OK;
    } catch (cv::Exception& e) {
        (void)e;
        return ERR_OPENCV;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_resize(uint64_t src_handle, uint64_t dst_handle, int32_t width, int32_t height) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::resize(*src, *dst, cv::Size(width, height));
        return OK;
    } catch (cv::Exception& e) {
        (void)e;
        return ERR_OPENCV;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_blur(uint64_t src_handle, uint64_t dst_handle, int32_t k_width, int32_t k_height) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::blur(*src, *dst, cv::Size(k_width, k_height));
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_gaussian_blur(uint64_t src_handle, uint64_t dst_handle, int32_t k_width, int32_t k_height, double sigma_x) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::GaussianBlur(*src, *dst, cv::Size(k_width, k_height), sigma_x);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_median_blur(uint64_t src_handle, uint64_t dst_handle, int32_t ksize) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::medianBlur(*src, *dst, ksize);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_threshold(uint64_t src_handle, uint64_t dst_handle, double thresh, double maxval, int32_t type) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::threshold(*src, *dst, thresh, maxval, type);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_adaptive_threshold(uint64_t src_handle, uint64_t dst_handle, double maxval, int32_t adaptive_method, int32_t threshold_type, int32_t block_size, double c) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::adaptiveThreshold(*src, *dst, maxval, adaptive_method, threshold_type, block_size, c);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_canny(uint64_t src_handle, uint64_t dst_handle, double threshold1, double threshold2) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Canny(*src, *dst, threshold1, threshold2);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: geometry/gradient
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_flip(uint64_t src_handle, uint64_t dst_handle, int32_t flip_code) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::flip(*src, *dst, flip_code);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_sobel(uint64_t src_handle, uint64_t dst_handle, int32_t ddepth, int32_t dx, int32_t dy, int32_t ksize, double scale, double delta) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Sobel(*src, *dst, ddepth, dx, dy, ksize, scale, delta);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_laplacian(uint64_t src_handle, uint64_t dst_handle, int32_t ddepth, int32_t ksize, double scale, double delta) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Laplacian(*src, *dst, ddepth, ksize, scale, delta);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_transpose(uint64_t src_handle, uint64_t dst_handle) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::transpose(*src, *dst);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: morphology
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_erode(uint64_t src_handle, uint64_t dst_handle, uint64_t kernel_handle, int32_t anchor_x, int32_t anchor_y, int32_t iterations) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Mat kernel;
        if (kernel_handle != 0) kernel = *(cv::Mat*)(intptr_t)kernel_handle;
        cv::Point anchor(anchor_x, anchor_y);
        if (anchor_x < 0 && anchor_y < 0) anchor = cv::Point(-1, -1);
        cv::erode(*src, *dst, kernel, anchor, iterations);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_dilate(uint64_t src_handle, uint64_t dst_handle, uint64_t kernel_handle, int32_t anchor_x, int32_t anchor_y, int32_t iterations) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Mat kernel;
        if (kernel_handle != 0) kernel = *(cv::Mat*)(intptr_t)kernel_handle;
        cv::Point anchor(anchor_x, anchor_y);
        if (anchor_x < 0 && anchor_y < 0) anchor = cv::Point(-1, -1);
        cv::dilate(*src, *dst, kernel, anchor, iterations);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_morphology_ex(uint64_t src_handle, uint64_t dst_handle, int32_t op, uint64_t kernel_handle, int32_t anchor_x, int32_t anchor_y, int32_t iterations) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::Mat kernel;
        if (kernel_handle != 0) kernel = *(cv::Mat*)(intptr_t)kernel_handle;
        cv::Point anchor(anchor_x, anchor_y);
        if (anchor_x < 0 && anchor_y < 0) anchor = cv::Point(-1, -1);
        cv::morphologyEx(*src, *dst, op, kernel, anchor, iterations);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_imgproc_get_structuring_element(int32_t shape, int32_t kwidth, int32_t kheight) {
    try {
        auto m = cv::getStructuringElement(shape, cv::Size(kwidth, kheight));
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: histogram
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_equalize_hist(uint64_t src_handle, uint64_t dst_handle) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::equalizeHist(*src, *dst);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_normalize(uint64_t src_handle, uint64_t dst_handle, double alpha, double beta, int32_t norm_type) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        cv::normalize(*src, *dst, alpha, beta, norm_type);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: drawing
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_rectangle(
    uint64_t img_handle,
    int32_t x1, int32_t y1, int32_t x2, int32_t y2,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness
) {
    if (img_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        cv::rectangle(*img, cv::Point(x1, y1), cv::Point(x2, y2),
                      cv::Scalar(b, g, r, a), thickness);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_circle(
    uint64_t img_handle,
    int32_t cx, int32_t cy, int32_t radius,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness
) {
    if (img_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        cv::circle(*img, cv::Point(cx, cy), radius,
                   cv::Scalar(b, g, r, a), thickness);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_line(
    uint64_t img_handle,
    int32_t x1, int32_t y1, int32_t x2, int32_t y2,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness
) {
    if (img_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        cv::line(*img, cv::Point(x1, y1), cv::Point(x2, y2),
                 cv::Scalar(b, g, r, a), thickness);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_put_text(
    uint64_t img_handle,
    const char* text, int32_t text_len,
    int32_t org_x, int32_t org_y,
    int32_t font_face, double font_scale,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness, int32_t line_type, int32_t bottom_left_origin
) {
    if (img_handle == 0 || !text) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        std::string txt(text, text_len);
        cv::putText(*img, txt, cv::Point(org_x, org_y),
                    font_face, font_scale,
                    cv::Scalar(b, g, r, a), thickness, line_type, bottom_left_origin != 0);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_fill_poly(
    uint64_t img_handle,
    const int32_t* pts, int32_t npts, int32_t ncontours,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t line_type, int32_t shift
) {
    if (img_handle == 0 || !pts) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        // pts layout: [ncontours arrays of (x,y) pairs]
        // npts is array of contour lengths
        // For simplicity, single contour: pts = [{x0,y0},{x1,y1},...], npts = count
        std::vector<cv::Point> points;
        for (int i = 0; i < ncontours; i++) {
            points.push_back(cv::Point(pts[i*2], pts[i*2+1]));
        }
        const cv::Point* ppt[1] = { points.data() };
        int npt[1] = { (int)points.size() };
        cv::fillPoly(*img, ppt, npt, 1, cv::Scalar(b, g, r, a), line_type, shift);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_arrowed_line(
    uint64_t img_handle,
    int32_t x1, int32_t y1, int32_t x2, int32_t y2,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness, int32_t line_type, int32_t shift, double tip_length
) {
    if (img_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        cv::arrowedLine(*img, cv::Point(x1, y1), cv::Point(x2, y2),
                        cv::Scalar(b, g, r, a), thickness, line_type, shift, tip_length);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: contours
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_find_contours(
    uint64_t src_handle,
    uint64_t contours_vec_handle,
    int32_t mode, int32_t method,
    int32_t offset_x, int32_t offset_y
) {
    if (src_handle == 0 || contours_vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* contours = (std::vector<std::vector<cv::Point>>*)(intptr_t)contours_vec_handle;
        cv::findContours(*src, *contours, mode, method, cv::Point(offset_x, offset_y));
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_draw_contours(
    uint64_t img_handle,
    uint64_t contours_vec_handle,
    int32_t contour_idx,
    uint8_t r, uint8_t g, uint8_t b, uint8_t a,
    int32_t thickness
) {
    if (img_handle == 0 || contours_vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* img = (cv::Mat*)(intptr_t)img_handle;
        auto* contours = (std::vector<std::vector<cv::Point>>*)(intptr_t)contours_vec_handle;
        cv::drawContours(*img, *contours, contour_idx, cv::Scalar(b, g, r, a), thickness);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT double goopencv_imgproc_contour_area(uint64_t vec_handle) {
    if (vec_handle == 0) return 0;
    try {
        auto* pts = (std::vector<cv::Point>*)(intptr_t)vec_handle;
        return cv::contourArea(*pts);
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT double goopencv_imgproc_arc_length(uint64_t vec_handle, int32_t closed) {
    if (vec_handle == 0) return 0;
    try {
        auto* pts = (std::vector<cv::Point>*)(intptr_t)vec_handle;
        return cv::arcLength(*pts, closed != 0);
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_bounding_rect(uint64_t vec_handle, int32_t* out_x, int32_t* out_y, int32_t* out_w, int32_t* out_h) {
    if (vec_handle == 0 || !out_x || !out_y || !out_w || !out_h) return ERR_INVALID_HANDLE;
    try {
        auto* pts = (std::vector<cv::Point>*)(intptr_t)vec_handle;
        cv::Rect r = cv::boundingRect(*pts);
        *out_x = r.x;
        *out_y = r.y;
        *out_w = r.width;
        *out_h = r.height;
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_min_enclosing_circle(uint64_t vec_handle, double* out_cx, double* out_cy, double* out_radius) {
    if (vec_handle == 0 || !out_cx || !out_cy || !out_radius) return ERR_INVALID_HANDLE;
    try {
        auto* pts = (std::vector<cv::Point>*)(intptr_t)vec_handle;
        cv::Point2f center;
        float radius;
        cv::minEnclosingCircle(*pts, center, radius);
        *out_cx = center.x;
        *out_cy = center.y;
        *out_radius = radius;
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_moments(
    uint64_t vec_handle, int32_t binary,
    double* out_m00, double* out_m10, double* out_m01,
    double* out_m20, double* out_m11, double* out_m02,
    double* out_m30, double* out_m21, double* out_m12, double* out_m03
) {
    if (vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* pts = (std::vector<cv::Point>*)(intptr_t)vec_handle;
        cv::Moments m = cv::moments(*pts, binary != 0);
        if (out_m00) *out_m00 = m.m00;
        if (out_m10) *out_m10 = m.m10;
        if (out_m01) *out_m01 = m.m01;
        if (out_m20) *out_m20 = m.m20;
        if (out_m11) *out_m11 = m.m11;
        if (out_m02) *out_m02 = m.m02;
        if (out_m30) *out_m30 = m.m30;
        if (out_m21) *out_m21 = m.m21;
        if (out_m12) *out_m12 = m.m12;
        if (out_m03) *out_m03 = m.m03;
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: Hough
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_hough_lines(
    uint64_t src_handle, uint64_t vec_handle,
    double rho, double theta, int32_t threshold,
    double srn, double stn
) {
    if (src_handle == 0 || vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* lines = (std::vector<cv::Vec2f>*)(intptr_t)vec_handle;
        cv::HoughLines(*src, *lines, rho, theta, threshold, srn, stn);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_hough_lines_p(
    uint64_t src_handle, uint64_t vec_handle,
    double rho, double theta, int32_t threshold,
    double min_line_length, double max_line_gap
) {
    if (src_handle == 0 || vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* lines = (std::vector<cv::Vec4i>*)(intptr_t)vec_handle;
        cv::HoughLinesP(*src, *lines, rho, theta, threshold, min_line_length, max_line_gap);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_hough_circles(
    uint64_t src_handle, uint64_t vec_handle,
    int32_t method, double dp, double min_dist,
    double param1, double param2, int32_t min_radius, int32_t max_radius
) {
    if (src_handle == 0 || vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* circles = (std::vector<cv::Vec3f>*)(intptr_t)vec_handle;
        cv::HoughCircles(*src, *circles, method, dp, min_dist, param1, param2, min_radius, max_radius);
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

// ---------------------------------------------------------------------------
// Image processing — imgproc: warp
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_imgproc_warp_affine(uint64_t src_handle, uint64_t dst_handle, uint64_t m_handle, int32_t dst_w, int32_t dst_h) {
    if (src_handle == 0 || dst_handle == 0 || m_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        auto* m = (cv::Mat*)(intptr_t)m_handle;
        cv::warpAffine(*src, *dst, *m, cv::Size(dst_w, dst_h));
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT int32_t goopencv_imgproc_warp_perspective(uint64_t src_handle, uint64_t dst_handle, uint64_t m_handle, int32_t dst_w, int32_t dst_h) {
    if (src_handle == 0 || dst_handle == 0 || m_handle == 0) return ERR_INVALID_HANDLE;
    try {
        auto* src = (cv::Mat*)(intptr_t)src_handle;
        auto* dst = (cv::Mat*)(intptr_t)dst_handle;
        auto* m = (cv::Mat*)(intptr_t)m_handle;
        cv::warpPerspective(*src, *dst, *m, cv::Size(dst_w, dst_h));
        return OK;
    } catch (...) {
        return ERR_OPENCV;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_imgproc_get_rotation_matrix2d(double center_x, double center_y, double angle, double scale) {
    try {
        auto m = cv::getRotationMatrix2D(cv::Point2f((float)center_x, (float)center_y), angle, scale);
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

extern "C" ABI_EXPORT uint64_t goopencv_imgproc_get_affine_transform(
    double src0x, double src0y, double src1x, double src1y, double src2x, double src2y,
    double dst0x, double dst0y, double dst1x, double dst1y, double dst2x, double dst2y
) {
    try {
        cv::Point2f srcPts[3] = {
            cv::Point2f((float)src0x, (float)src0y),
            cv::Point2f((float)src1x, (float)src1y),
            cv::Point2f((float)src2x, (float)src2y)
        };
        cv::Point2f dstPts[3] = {
            cv::Point2f((float)dst0x, (float)dst0y),
            cv::Point2f((float)dst1x, (float)dst1y),
            cv::Point2f((float)dst2x, (float)dst2y)
        };
        auto m = cv::getAffineTransform(srcPts, dstPts);
        auto* dst = new cv::Mat(m);
        return (uint64_t)(intptr_t)dst;
    } catch (...) {
        return 0;
    }
}

// ---------------------------------------------------------------------------
// Core operations
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT int32_t goopencv_core_bitwise_and(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::bitwise_and(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_bitwise_or(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::bitwise_or(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_bitwise_xor(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::bitwise_xor(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_bitwise_not(uint64_t src_handle, uint64_t dst_handle) {
    if (src_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::bitwise_not(*(cv::Mat*)(intptr_t)src_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_add(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::add(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_subtract(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::subtract(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_multiply(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle, double scale, int32_t dtype) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::multiply(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle, scale, dtype);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_divide(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle, double scale, int32_t dtype) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::divide(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle, scale, dtype);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_abs_diff(uint64_t src1_handle, uint64_t src2_handle, uint64_t dst_handle) {
    if (src1_handle == 0 || src2_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::absdiff(*(cv::Mat*)(intptr_t)src1_handle, *(cv::Mat*)(intptr_t)src2_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_min_max_loc(
    uint64_t src_handle,
    double* out_min, double* out_max,
    int32_t* out_min_x, int32_t* out_min_y,
    int32_t* out_max_x, int32_t* out_max_y
) {
    if (src_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::Point minLoc, maxLoc;
        double minVal, maxVal;
        cv::minMaxLoc(*(cv::Mat*)(intptr_t)src_handle, &minVal, &maxVal, &minLoc, &maxLoc);
        if (out_min) *out_min = minVal;
        if (out_max) *out_max = maxVal;
        if (out_min_x) *out_min_x = minLoc.x;
        if (out_min_y) *out_min_y = minLoc.y;
        if (out_max_x) *out_max_x = maxLoc.x;
        if (out_max_y) *out_max_y = maxLoc.y;
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_mean_std_dev(uint64_t src_handle, uint64_t mean_handle, uint64_t stddev_handle) {
    if (src_handle == 0 || mean_handle == 0 || stddev_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::meanStdDev(*(cv::Mat*)(intptr_t)src_handle, *(cv::Mat*)(intptr_t)mean_handle, *(cv::Mat*)(intptr_t)stddev_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_count_non_zero(uint64_t src_handle) {
    if (src_handle == 0) return -1;
    try {
        return (int32_t)cv::countNonZero(*(cv::Mat*)(intptr_t)src_handle);
    } catch (...) { return -1; }
}

// ---------------------------------------------------------------------------
// Vector helpers — std::vector<cv::Point> (contour)
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_vec_points_new() {
    try {
        auto* v = new std::vector<cv::Point>();
        return (uint64_t)(intptr_t)v;
    } catch (...) { return 0; }
}

extern "C" ABI_EXPORT void goopencv_vec_points_push(uint64_t handle, int32_t x, int32_t y) {
    if (handle == 0) return;
    ((std::vector<cv::Point>*)(intptr_t)handle)->push_back(cv::Point(x, y));
}

extern "C" ABI_EXPORT int32_t goopencv_vec_points_len(uint64_t handle) {
    if (handle == 0) return 0;
    return (int32_t)((std::vector<cv::Point>*)(intptr_t)handle)->size();
}

extern "C" ABI_EXPORT int32_t goopencv_vec_points_get(uint64_t handle, int32_t idx, int32_t* out_x, int32_t* out_y) {
    if (handle == 0 || !out_x || !out_y) return ERR_INVALID_HANDLE;
    auto* v = (std::vector<cv::Point>*)(intptr_t)handle;
    if (idx < 0 || idx >= (int32_t)v->size()) return ERR_INVALID_ARGUMENT;
    *out_x = (*v)[idx].x;
    *out_y = (*v)[idx].y;
    return OK;
}

extern "C" ABI_EXPORT void goopencv_vec_points_delete(uint64_t handle) {
    if (handle == 0) return;
    delete (std::vector<cv::Point>*)(intptr_t)handle;
}

// ---------------------------------------------------------------------------
// Vector helpers — std::vector<std::vector<cv::Point>> (contours)
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_vec_vec_points_new() {
    try {
        auto* v = new std::vector<std::vector<cv::Point>>();
        return (uint64_t)(intptr_t)v;
    } catch (...) { return 0; }
}

extern "C" ABI_EXPORT void goopencv_vec_vec_points_push(uint64_t handle, uint64_t contour_handle) {
    if (handle == 0 || contour_handle == 0) return;
    auto* v = (std::vector<std::vector<cv::Point>>*)(intptr_t)handle;
    auto* c = (std::vector<cv::Point>*)(intptr_t)contour_handle;
    v->push_back(*c);
}

extern "C" ABI_EXPORT int32_t goopencv_vec_vec_points_len(uint64_t handle) {
    if (handle == 0) return 0;
    return (int32_t)((std::vector<std::vector<cv::Point>>*)(intptr_t)handle)->size();
}

extern "C" ABI_EXPORT uint64_t goopencv_vec_vec_points_get(uint64_t handle, int32_t idx) {
    if (handle == 0) return 0;
    auto* v = (std::vector<std::vector<cv::Point>>*)(intptr_t)handle;
    if (idx < 0 || idx >= (int32_t)v->size()) return 0;
    auto* c = new std::vector<cv::Point>((*v)[idx]);
    return (uint64_t)(intptr_t)c;
}

extern "C" ABI_EXPORT void goopencv_vec_vec_points_delete(uint64_t handle) {
    if (handle == 0) return;
    delete (std::vector<std::vector<cv::Point>>*)(intptr_t)handle;
}

// ---------------------------------------------------------------------------
// Vector helpers — std::vector<double>
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_vec_double_new() {
    try { auto* v = new std::vector<double>(); return (uint64_t)(intptr_t)v; }
    catch (...) { return 0; }
}

extern "C" ABI_EXPORT double goopencv_vec_double_get(uint64_t handle, int32_t idx) {
    if (handle == 0) return 0;
    auto* v = (std::vector<double>*)(intptr_t)handle;
    if (idx < 0 || idx >= (int32_t)v->size()) return 0;
    return (*v)[idx];
}

extern "C" ABI_EXPORT int32_t goopencv_vec_double_len(uint64_t handle) {
    if (handle == 0) return 0;
    return (int32_t)((std::vector<double>*)(intptr_t)handle)->size();
}

extern "C" ABI_EXPORT void goopencv_vec_double_delete(uint64_t handle) {
    if (handle == 0) return;
    delete (std::vector<double>*)(intptr_t)handle;
}

// ---------------------------------------------------------------------------
// Vector helpers — std::vector<int>
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_vec_int_new() {
    try { auto* v = new std::vector<int>(); return (uint64_t)(intptr_t)v; }
    catch (...) { return 0; }
}

extern "C" ABI_EXPORT int32_t goopencv_vec_int_get(uint64_t handle, int32_t idx) {
    if (handle == 0) return 0;
    auto* v = (std::vector<int>*)(intptr_t)handle;
    if (idx < 0 || idx >= (int32_t)v->size()) return 0;
    return (*v)[idx];
}

extern "C" ABI_EXPORT int32_t goopencv_vec_int_len(uint64_t handle) {
    if (handle == 0) return 0;
    return (int32_t)((std::vector<int>*)(intptr_t)handle)->size();
}

extern "C" ABI_EXPORT void goopencv_vec_int_delete(uint64_t handle) {
    if (handle == 0) return;
    delete (std::vector<int>*)(intptr_t)handle;
}

// ---------------------------------------------------------------------------
// Vector helpers — std::vector<cv::Mat> (for split)
// ---------------------------------------------------------------------------

extern "C" ABI_EXPORT uint64_t goopencv_vec_mat_new() {
    try { auto* v = new std::vector<cv::Mat>(); return (uint64_t)(intptr_t)v; }
    catch (...) { return 0; }
}

extern "C" ABI_EXPORT void goopencv_vec_mat_push(uint64_t handle, uint64_t mat_handle) {
    if (handle == 0 || mat_handle == 0) return;
    auto* v = (std::vector<cv::Mat>*)(intptr_t)handle;
    auto* m = (cv::Mat*)(intptr_t)mat_handle;
    v->push_back(m->clone());
}

extern "C" ABI_EXPORT int32_t goopencv_vec_mat_len(uint64_t handle) {
    if (handle == 0) return 0;
    return (int32_t)((std::vector<cv::Mat>*)(intptr_t)handle)->size();
}

extern "C" ABI_EXPORT uint64_t goopencv_vec_mat_get(uint64_t handle, int32_t idx) {
    if (handle == 0) return 0;
    auto* v = (std::vector<cv::Mat>*)(intptr_t)handle;
    if (idx < 0 || idx >= (int32_t)v->size()) return 0;
    auto* m = new cv::Mat((*v)[idx].clone());
    return (uint64_t)(intptr_t)m;
}

extern "C" ABI_EXPORT void goopencv_vec_mat_delete(uint64_t handle) {
    if (handle == 0) return;
    delete (std::vector<cv::Mat>*)(intptr_t)handle;
}

// ---------------------------------------------------------------------------
// Vector helpers — HoughLines returns Vec2f (rho, theta)
// ---------------------------------------------------------------------------

// We reuse vec_double for Hough lines: each line = 2 doubles (rho, theta)
// HoughLinesP returns Vec4i: each line = 4 ints (x1,y1,x2,y2)
// HoughCircles returns Vec3f: each circle = 3 doubles (x,y,r)

// Split/Merge
extern "C" ABI_EXPORT int32_t goopencv_core_split(uint64_t src_handle, uint64_t vec_handle) {
    if (src_handle == 0 || vec_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::split(*(cv::Mat*)(intptr_t)src_handle, *(std::vector<cv::Mat>*)(intptr_t)vec_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}

extern "C" ABI_EXPORT int32_t goopencv_core_merge(uint64_t vec_handle, uint64_t dst_handle) {
    if (vec_handle == 0 || dst_handle == 0) return ERR_INVALID_HANDLE;
    try {
        cv::merge(*(std::vector<cv::Mat>*)(intptr_t)vec_handle, *(cv::Mat*)(intptr_t)dst_handle);
        return OK;
    } catch (...) { return ERR_OPENCV; }
}
