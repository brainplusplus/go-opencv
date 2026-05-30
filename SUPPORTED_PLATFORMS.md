# Supported Platforms

go-opencv ships prebuilt native libraries for five runtime targets, with a universal macOS dylib covering both Intel and Apple Silicon.

| # | GOOS | GOARCH | Native Library | Target Hardware |
|---|------|--------|----------------|-----------------|
| 1 | `windows` | `amd64` | `goopencv.dll` | Windows x64 |
| 2 | `linux` | `amd64` | `goopencv.so` | Linux x64 (glibc 2.35+) |
| 3 | `linux` | `arm64` | `goopencv-linux-arm64.so` | Linux ARM64 (RPi 5, AWS Graviton, Jetson) |
| 4 | `darwin` | `amd64` | `goopencv.dylib` | Intel Mac |
| 5 | `darwin` | `arm64` | `goopencv.dylib` | Apple Silicon |

## macOS Universal Build

The macOS artifact is a single universal `goopencv.dylib` with both `x86_64` and `arm64` slices. At runtime, Go loads the same embedded dylib on either Intel Macs or Apple Silicon Macs.

## Unsupported Platforms

If you run on an unsupported platform, `opencv.New()` returns:

```
opencv: backend unavailable (supported: windows/amd64, linux/amd64, linux/arm64, darwin/amd64, darwin/arm64)
```

## Autoload

The native library is embedded at compile time via `//go:embed`. No extra DLL/SO/DYLIB download is required.

```go
rt, err := opencv.New(ctx) // auto-detects platform and extracts to cache
```

Manual override:

```go
rt, err := opencv.New(ctx, opencv.WithDLL("/path/to/goopencv.dylib"))
```

## Linux glibc Requirement

| Distribution | Minimum Version |
|-------------|----------------|
| Ubuntu 22.04+ | glibc 2.35 |
| Debian Bookworm+ | glibc 2.36 |

Older distributions such as Ubuntu 20.04 and Debian Buster are not supported due to glibc ABI requirements in the prebuilt Linux binaries.
