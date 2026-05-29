# API Freeze Notes (v1-prep)

## Core Principle

Single OpenCV-style API surface only.

- Keep OpenCV naming and semantics.
- Add explicit extensions through existing methods.
- No separate Go helper namespace/layer.

## Canonical APIs

- `IMRead(path, ...ColorModel)`
- `IMReadBytes(data, ...ColorModel)`
- `IMWrite(path, mat, ...ColorModel)`
- `CvtColor(src, dst, code)`
- `Resize(src, dst, size)`
- `GaussianBlur(src, kSize, sigmaX)`
- `Threshold`, `AdaptiveThreshold`, `Canny`
- Drawing: `Rectangle`, `Circle`, `Line`

## Color model enum

- `Unknown`
- `BGR`
- `RGB`
- `RGBA`
- `Gray`

## Contracts

- `IMRead(path)` default is `BGR` (OpenCV parity).
- `IMRead(path, model)` performs explicit model conversion.
- No hidden conversions in processing APIs.
- `CvtColor` remains explicit conversion-code driven.
- `Mat` tracks color model metadata best-effort.
- Unknown model must propagate conservatively.

## Strict validation

- Off by default (parity-safe).
- Optional runtime gate:
  - `SetStrictColorValidation(true)`
  - blocks obvious model/code mismatches.

## Naming conflict resolution

- Model enum uses short names (`BGR`, `RGB`, `RGBA`, `Gray`).
- Drawing color structs use explicit names (`BGRColor`, `RGBAColor`).

## Removed legacy aliases

- `ReadImage`, `ReadImageBytes`, `SaveImage` removed.
- Use `IMRead`, `IMReadBytes`, `IMWrite` only.

## Performance posture

- No implicit per-frame conversion.
- IMWrite conversion path uses pooled RGBA images and pooled byte buffers.
- Benchmarks tracked for 1080p/4K conversion and write overrides.
