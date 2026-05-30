# Supported Platforms

go-opencv ships pre-built native libraries for 4 platform/architecture combinations:

| # | GOOS | GOARCH | Native Library | Target Hardware |
|---|------|--------|----------------|-----------------|
| 1 | `windows` | `amd64` | `goopencv.dll` | Windows x64 |
| 2 | `linux` | `amd64` | `goopencv.so` | Linux x64 (glibc 2.35+) |
| 3 | `linux` | `arm64` | `goopencv-linux-arm64.so` | Linux ARM64 (RPi 5, AWS Graviton, Jetson) |
| 4 | `darwin` | `arm64` | `goopencv.dylib` | Apple Silicon (M1/M2/M3/M4) |

## Why No Intel Mac?

opencv-mobile provides macOS SDK for ARM64 only. Intel Mac (x86_64) is not supported.
Apple discontinued Intel Macs in 2020. If you need Intel Mac support, you must build OpenCV from source.

## Unsupported Platforms

If you run on an unsupported platform, `opencv.New()` returns:

```
opencv: backend unavailable (supported: windows/amd64, linux/amd64, linux/arm64, darwin/arm64)
```

## Autoload

The native library is embedded at compile time via `//go:embed`. No external DLL/SO download needed.

```go
rt, err := opencv.New(ctx)  // auto-detects platform, extracts to cache
```

Manual override:

```go
rt, err := opencv.New(ctx, opencv.WithDLL("/path/to/goopencv.dll"))
```

## Linux glibc Requirement

| Distribution | Minimum Version |
|-------------|----------------|
| Ubuntu 22.04+ | glibc 2.35 |
| Debian Bookworm+ | glibc 2.36 |

Older distributions (Ubuntu 20.04, Debian Buster) are **not supported** due to glibc ABI.
