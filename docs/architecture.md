# Architecture

## Decision

go-opencv is a real Go library at module path `github.com/brainplusplus/go-opencv`.

It must remain CGO-free. OpenCV compute runs in a standalone wasm guest loaded by a Go host runtime. Public API is Go-native but intentionally shaped after OpenCV.js concepts.

## Non-decision

`references/wasm/opencv.wasm` is not the production backend target. Audit showed it depends on Emscripten/OpenCV.js glue. It remains reference material for contract and parity tests.

`references/go-wasmcv` is a contract seed. Its generated WIT bindings model wasm guests importing host OpenCV functions, while go-opencv runs the inverse direction: Go host calling OpenCV wasm backend. Reuse names, resources, module boundaries, and enum coverage; do not expose generated WIT/CM types in public API.

## Layers

1. Public package `github.com/brainplusplus/go-opencv`
   - `Runtime`
   - `Mat`
   - OpenCV.js-like types: `Point`, `Size`, `Rect`, `Scalar`
   - operations: `CvtColor`, `Resize`, then broader OpenCV modules

2. Internal runtime host
   - pure Go wasm runtime via wazero
   - module lifecycle
   - memory and pointer marshalling
   - error-code translation

3. Internal ABI spec
   - stable flat exports like `goopencv_mat_new`
   - no embind
   - no JS glue
   - no DOM/browser dependencies

4. Backend wasm build
   - future OpenCV build with explicit C-style exports
   - generated Go bindings for broad coverage

## ABI principles

- Handles cross boundary as integers.
- Buffers cross boundary as pointer+length.
- API returns explicit error codes.
- Go owns high-level lifecycle.
- Wasm owns OpenCV object memory.
- No OpenCV C++ object layout exposed to Go.

## Full-feature path

Full feature means broad OpenCV module coverage via generated ABI bindings, not hand-written wrappers for every symbol.

Planned coverage groups:

- core: `Mat`, scalar/vector types, arithmetic, array ops
- imgcodecs: read/write, encode/decode from bytes
- imgproc: color, resize, filters, threshold, morphology, contours
- features2d/calib3d: keypoints, matchers, transforms
- video/objdetect/dnn: separate build profiles if binary size/perf require it

## Verification strategy

- Go unit tests for lifecycle and error handling.
- ABI conformance tests against tiny fixture wasm.
- Golden image tests comparing output bytes/statistics.
- Parity tests derived from OpenCV.js typings and examples.
- Memory leak tests: every handle allocated must be released.
