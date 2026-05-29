# Performance Benchmark Plan

This plan defines mandatory baseline benchmarks for color-model and conversion-sensitive paths.

## Goals

- Detect regressions in image I/O and color conversions.
- Keep video/frame processing predictable.
- Track allocation behavior for conversion-heavy paths.

## Benchmark Matrix

### Read/Decode

- `BenchmarkIMRead_BRG_1080p`
- `BenchmarkIMRead_RGB_1080p`
- `BenchmarkIMRead_RGBA_1080p`
- `BenchmarkIMRead_Gray_1080p`

### ConvertModel / CvtColor

- `BenchmarkCvtColor_BGR2Gray_1080p`
- `BenchmarkCvtColor_BGR2RGBA_1080p`
- `BenchmarkCvtColor_RGBA2BGR_1080p`
- `BenchmarkCvtColor_BGR2RGBA_4K`

### Save/Encode

- `BenchmarkIMWrite_PNG_1080p_BGR`
- `BenchmarkIMWrite_JPEG_1080p_BGR`
- `BenchmarkIMWrite_PNG_1080p_RGBA`

### Frame-loop (pipeline realism)

- `BenchmarkPipeline_NoConvert_1080p`
- `BenchmarkPipeline_BGR2RGBA_PerFrame_1080p`

## Metrics to watch

- ns/op
- B/op
- allocs/op

## Guardrails

- No hidden auto-conversion in read/process/write path.
- Per-frame conversion must always be explicit in code.
- Conversion API should document copy semantics clearly.
