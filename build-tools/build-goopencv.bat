@echo off
REM Build goopencv.dll from goopencv_abi.cpp + opencv-mobile static libs
REM Supports Windows amd64 and arm64 builds.

setlocal enabledelayedexpansion

if defined GITHUB_WORKSPACE (
    set "REPO_ROOT=%GITHUB_WORKSPACE%"
) else (
    set "REPO_ROOT=%~dp0.."
)

if "%TARGET_ARCH%"=="" set "TARGET_ARCH=amd64"
if /I "%TARGET_ARCH%"=="x64" set "TARGET_ARCH=amd64"

if /I "%TARGET_ARCH%"=="amd64" (
    set "SDK_ARCH_DIR=x64"
    set "MSVC_SUBDIR=x64"
) else if /I "%TARGET_ARCH%"=="arm64" (
    set "SDK_ARCH_DIR=arm64"
    set "MSVC_SUBDIR=ARM64"
) else (
    echo ERROR: Unsupported TARGET_ARCH "%TARGET_ARCH%". Use amd64 or arm64.
    exit /b 1
)

if defined VSCMD_VER (
    echo Already in VS Developer Command Prompt.
) else (
    set "VS_PATH=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools"
    set "VS_PATH_ALT=C:\Program Files\Microsoft Visual Studio\2022\Enterprise"
    set "VS_PATH_COMM=C:\Program Files\Microsoft Visual Studio\2022\Community"
    set "VCVARS_BAT="

    for %%P in ("!VS_PATH!" "!VS_PATH_ALT!" "!VS_PATH_COMM!") do (
        if "!VCVARS_BAT!"=="" (
            if /I "%TARGET_ARCH%"=="amd64" if exist "%%~P\VC\Auxiliary\Build\vcvars64.bat" set "VCVARS_BAT=%%~P\VC\Auxiliary\Build\vcvars64.bat"
            if /I "%TARGET_ARCH%"=="arm64" if exist "%%~P\VC\Auxiliary\Build\vcvarsarm64.bat" set "VCVARS_BAT=%%~P\VC\Auxiliary\Build\vcvarsarm64.bat"
            if /I "%TARGET_ARCH%"=="arm64" if "!VCVARS_BAT!"=="" if exist "%%~P\VC\Auxiliary\Build\vcvarsamd64_arm64.bat" set "VCVARS_BAT=%%~P\VC\Auxiliary\Build\vcvarsamd64_arm64.bat"
        )
    )

    if "!VCVARS_BAT!"=="" (
        echo ERROR: Visual Studio 2022 with %TARGET_ARCH% toolchain not found.
        exit /b 1
    )

    call "!VCVARS_BAT!" >nul 2>&1
)

set "OPENCV_ROOT=%REPO_ROOT%\build-tools\opencv-mobile-4.13.0-windows-vs2022\opencv-mobile-4.13.0-windows-vs2022"
if exist "!OPENCV_ROOT!" goto :found_root
set "OPENCV_ROOT=%REPO_ROOT%\build-tools\opencv-mobile-4.13.0-windows-vs2022"
if exist "!OPENCV_ROOT!" goto :found_root
set "OPENCV_ROOT=%REPO_ROOT%\build-tools\opencv-mobile-4.13.0-windows-vs2019\opencv-mobile-4.13.0-windows-vs2019"
if exist "!OPENCV_ROOT!" goto :found_root
set "OPENCV_ROOT=%REPO_ROOT%\build-tools\opencv-mobile-4.13.0-windows-vs2019"
if exist "!OPENCV_ROOT!" goto :found_root
echo ERROR: opencv-mobile SDK not found
exit /b 1
:found_root

if not exist "!OPENCV_ROOT!\%SDK_ARCH_DIR%" (
    echo ERROR: %SDK_ARCH_DIR% directory not found in !OPENCV_ROOT!
    exit /b 1
)

set "INCLUDE_DIR=!OPENCV_ROOT!\%SDK_ARCH_DIR%\include"
set "LIB_DIR="
for %%V in (vc17 vc16 vc15 vc14) do (
    if exist "!OPENCV_ROOT!\%SDK_ARCH_DIR%\%MSVC_SUBDIR%\%%V\staticlib\opencv_core4130.lib" set "LIB_DIR=!OPENCV_ROOT!\%SDK_ARCH_DIR%\%MSVC_SUBDIR%\%%V\staticlib"
    if exist "!OPENCV_ROOT!\%SDK_ARCH_DIR%\%%V\staticlib\opencv_core4130.lib" set "LIB_DIR=!OPENCV_ROOT!\%SDK_ARCH_DIR%\%%V\staticlib"
)
if "!LIB_DIR!"=="" (
    echo ERROR: staticlib not found for %TARGET_ARCH%
    dir /s /b "!OPENCV_ROOT!\%SDK_ARCH_DIR%\*staticlib" 2>nul
    exit /b 1
)

set "OUTPUT=!REPO_ROOT!\dist\goopencv.dll"
set "SOURCE=!REPO_ROOT!\backend\goopencv_abi.cpp"

if not exist "!REPO_ROOT!\dist" mkdir "!REPO_ROOT!\dist"

echo === Compiling goopencv.dll ===
echo TARGET_ARCH: %TARGET_ARCH%
echo REPO_ROOT:   !REPO_ROOT!
echo OPENCV_ROOT: !OPENCV_ROOT!
echo SDK_ARCH:    %SDK_ARCH_DIR%
echo INCLUDE:     !INCLUDE_DIR!
echo LIB:         !LIB_DIR!
echo Output:      !OUTPUT!
echo.

cl /LD /nologo /O2 /utf-8 /W3 /EHsc /MD ^
    /Fe:"!OUTPUT!" ^
    /I"!INCLUDE_DIR!" ^
    "!SOURCE!" ^
    "!LIB_DIR!\opencv_core4130.lib" ^
    "!LIB_DIR!\opencv_imgproc4130.lib" ^
    "!LIB_DIR!\opencv_photo4130.lib" ^
    "!LIB_DIR!\opencv_features2d4130.lib" ^
    "!LIB_DIR!\opencv_highgui4130.lib" ^
    "!LIB_DIR!\opencv_video4130.lib" ^
    ole32.lib ^
    user32.lib ^
    gdi32.lib

if !ERRORLEVEL! NEQ 0 (
    echo.
    echo === COMPILATION FAILED ===
    exit /b !ERRORLEVEL!
)

echo.
echo === SUCCESS ===
dir "!OUTPUT!"
