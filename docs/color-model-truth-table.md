# Color Model Truth Table

This document defines how `Mat` color model metadata must be handled per operation.

## Models

- `Unknown`
- `BGR`
- `RGB`
- `RGBA`
- `Gray`

## Rules

- **Preserve**: output model must equal input model.
- **Deterministic**: output model is fully determined by operation parameters.
- **Unknown**: output model cannot be trusted; set `Unknown` conservatively.
- **When in doubt -> Unknown.**

## Table

| Operation | Behavior | Output model rule |
|---|---|---|
| `IMRead(path)` | Deterministic | `BGR` (default) |
| `IMRead(path, model)` | Deterministic | `model` |
| `IMReadBytes(data)` | Deterministic | `BGR` (default) |
| `IMReadBytes(data, model)` | Deterministic | `model` |
| `CvtColor(src, dst, code)` | Deterministic | Based on conversion code (`*2Gray -> Gray`, `*2RGBA/*2BGRA -> RGBA`, `*2RGB/*2BGR -> RGB/BGR`) |
| `Resize(src, dst, size)` | Preserve | `dst = src` |
| `Blur(src)` | Preserve | `dst = src` |
| `GaussianBlur(src)` | Preserve | `dst = src` |
| `Threshold(src)` | Preserve
| `AdaptiveThreshold(src)` | Deterministic | `Gray` (requires single-channel source) |
| `Canny(src)` | Deterministic | `Gray` |
| `Rectangle/Circle/Line` | Preserve | In-place; model unchanged |
| `CopyTo(dstBytes)` | Preserve | Metadata unchanged |
| `IMWrite(path, mat)` | Preserve input metadata | No mutation of `mat` model |
| `Raw pointer wrap / external buffer` | Unknown | `Unknown` unless caller declares model |
| `Custom kernel / unknown transform` | Unknown | `Unknown` unless operation declares output model |

## Validation Modes

- Default mode: parity-first (no strict rejects, OpenCV-like behavior).
- Strict mode (optional): validate model/code compatibility and return error on mismatch.
- Unknown source model bypasses strict mismatch checks (parity-safe), then output model is still updated deterministically when conversion code implies it.
- IMWrite with Unknown model requires explicit override to prevent silent channel-order bugs.

## Examples

```go
img, _ := cv.IMRead("x.png")          // BGR
gray := cv.NewMat()
cv.CvtColor(img, gray, cv.ColorBGR2Gray) // Gray
```

```go
img, _ := cv.IMRead("x.png", cv.RGBA) // RGBA
gray := cv.NewMat()
cv.CvtColor(img, gray, cv.ColorRGBA2Gray) // Gray
```
