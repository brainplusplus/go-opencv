# Backend ABI Runtime Contract

Go host loads a wasm module with wazero. Backend module must export stable symbols listed in `docs/abi.md`.

Required exports for current runnable host path:

- `goopencv_mat_new(rows: i32, cols: i32, type: i32) -> i64`
- `goopencv_mat_delete(handle: i64) -> i32`
- `goopencv_mat_rows(handle: i64) -> i32`
- `goopencv_mat_cols(handle: i64) -> i32`
- `goopencv_mat_type(handle: i64) -> i32`
- `goopencv_mat_clone(handle: i64) -> i64`
- `goopencv_imgproc_cvt_color(src: i64, dst: i64, code: i32) -> i32`
- `goopencv_imgproc_resize(src: i64, dst: i64, width: i32, height: i32) -> i32`

Optional startup:

- If `_initialize` is exported, Go host calls it once after instantiation.
- WASI imports are supported through `wasi_snapshot_preview1`.

Error codes:

| Code | Meaning |
|---:|---|
| 0 | OK |
| 1 | invalid argument |
| 2 | invalid handle |
| 3 | out of memory |
| 4 | OpenCV error |
| 5 | unsupported |

Rules:

- Handles are opaque `uint64` values from wasm backend.
- `0` handle is invalid/null.
- Backend owns OpenCV object memory.
- Go owns high-level wrapper lifecycle and calls delete once.
- No C++ exception may cross wasm boundary; backend catches and returns error code.
- Do not depend on Emscripten embind/emval imports.
