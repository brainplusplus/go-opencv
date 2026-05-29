# Interface Contract

`references/go-wasmcv` is used as the initial interface seed for go-opencv.

The reference is generated from WIT and models wasm guest modules importing OpenCV functions from a host. go-opencv uses the opposite execution direction: Go host calls an OpenCV wasm backend. Still, its resource and module vocabulary is useful.

## Mapping rules

| wasmcv concept | go-opencv public API | go-opencv backend ABI |
|---|---|---|
| `cm.Resource` | opaque Go struct with `Close`/`Delete` | integer handle |
| `cm.Result[..., T, ErrorResult]` | `(T, error)` | error code + out parameter/handle |
| `wasm:cv/mat` | `Mat` methods | `goopencv_mat_*` |
| `wasm:cv/cv` | `Runtime` methods/functions | `goopencv_imgproc_*`, `goopencv_core_*` |
| generated enum | Go named integer type | integer code |

## Seed modules

- `types`: `Size`, `Point`, `Scalar`, `Rect`, `RGBA`, common enums.
- `mat`: `Mat` resource lifecycle, slicing, element access, metadata.
- `cv`: drawing, imgproc, array operations.
- `dnn`: blob params and network operations.
- `features2d`: keypoints, descriptors, matchers.
- `objdetect`: cascade/HOG/face detector resources.

## Public API policy

- Keep root package `github.com/brainplusplus/go-opencv` as user-facing library.
- Do not expose `go.bytecodealliance.org/cm` or WIT generated types publicly.
- Prefer OpenCV.js-ish naming where useful, Go error returns everywhere.
- Every owning object must have `Close() error`; OpenCV.js alias `Delete() error` may exist.
- Backend ABI remains flat and stable, even if public API grows.

## Backend ABI naming

Examples:

- `Mat.Rows()` -> `goopencv_mat_rows`
- `Mat.Clone()` -> `goopencv_mat_clone`
- `Runtime.Canny(...)` -> `goopencv_imgproc_canny`
- `Runtime.Circle(...)` -> `goopencv_imgproc_circle`
- `CascadeClassifier.DetectMultiScale(...)` -> `goopencv_objdetect_cascade_classifier_detect_multi_scale`

## Full-feature strategy

1. Seed contract from wasmcv interfaces.
2. Normalize names and types for Go public API.
3. Generate backend ABI symbol list.
4. Generate Go host call wrappers.
5. Implement C++ OpenCV wasm backend exports.
6. Add parity/golden tests per module.

## Generation

`cmd/goopencv-gen` renders contract-derived artifacts:

```shell
go run ./cmd/goopencv-gen -format abi-md -out docs/abi.md
go run ./cmd/goopencv-gen -format cpp -out backend/goopencv_abi.cpp
```

Generated files must not be hand-edited. Update `internal/contract` then regenerate.
