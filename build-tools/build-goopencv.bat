@echo off
REM Build goopencv.dll from goopencv_abi.cpp + opencv-mobile static libs
REM Works both locally and in CI (GitHub Actions windows-2022 runner)

setlocal enabledelayedexpansion

REM --- Detect repo root ---
REM CI sets GITHUB_WORKSPACE. Locally, derive from script location.
if defined GITHUB_WORKSPACE (
    set "REPO_ROOT=%GITHUB_WORKSPACE%"
) else (
    REM %~dp0 ends with trailing backslash, so parent = %~dp0..
    set "REPO_ROOT=%~dp0.."
)

REM --- Setup MSVC ---
if defined VSCMD_VER (
    echo Already in VS Developer Command Prompt.
) else (
    set "VS_PATH=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools"
    set "VS_PATH_ALT=C:\Program Files\Microsoft Visual Studio\2022\Enterprise"
    set "VS_PATH_COMM=C:\Program Files\Microsoft Visual Studio\2022\Community"

    if exist "!VS_PATH!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else if exist "!VS_PATH_ALT!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH_ALT!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else if exist "!VS_PATH_COMM!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH_COMM!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else (
        echo ERROR: Visual Studio 2022 not found.
        exit /b 1
    )
)

REM --- Resolve SDK root ---
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

REM --- Detect arch directory (x64 first for amd64 target, then arm64) ---
set "ARCH_DIR="
if exist "!OPENCV_ROOT!\x64" (
    set "ARCH_DIR=x64"
)
if "!ARCH_DIR!"=="" (
    if exist "!OPENCV_ROOT!\arm64" (
        set "ARCH_DIR=arm64"
    )
)
if "!ARCH_DIR!"=="" (
    echo ERROR: No x64 or arm64 directory found in !OPENCV_ROOT!
    exit /b 1
)

set "INCLUDE_DIR=!OPENCV_ROOT!\!ARCH_DIR!\include"

REM --- Detect staticlib path ---
set "LIB_DIR="
for %%V in (vc17 vc16 vc15 vc14) do (
    if exist "!OPENCV_ROOT!\!ARCH_DIR!\x64\%%V\staticlib\opencv_core4130.lib" (
        set "LIB_DIR=!OPENCV_ROOT!\!ARCH_DIR!\x64\%%V\staticlib"
    )
    if exist "!OPENCV_ROOT!\!ARCH_DIR!\ARM64\%%V\staticlib\opencv_core4130.lib" (
        set "LIB_DIR=!OPENCV_ROOT!\!ARCH_DIR!\ARM64\%%V\staticlib"
    )
    if exist "!OPENCV_ROOT!\!ARCH_DIR!\x86\%%V\staticlib\opencv_core4130.lib" (
        set "LIB_DIR=!OPENCV_ROOT!\!ARCH_DIR!\x86\%%V\staticlib"
    )
    if exist "!OPENCV_ROOT!\!ARCH_DIR!\%%V\staticlib\opencv_core4130.lib" (
        set "LIB_DIR=!OPENCV_ROOT!\!ARCH_DIR!\%%V\staticlib"
    )
)
if "!LIB_DIR!"=="" (
    echo ERROR: staticlib not found
    dir /s /b "!OPENCV_ROOT!\!ARCH_DIR!\*staticlib" 2>nul
    exit /b 1
)

set "OUTPUT=!REPO_ROOT!\dist\goopencv.dll"
set "SOURCE=!REPO_ROOT!\backend\goopencv_abi.cpp"

if not exist "!REPO_ROOT!\dist" mkdir "!REPO_ROOT!\dist"

echo === Compiling goopencv.dll ===
echo REPO_ROOT:   !REPO_ROOT!
echo OPENCV_ROOT: !OPENCV_ROOT!
echo ARCH_DIR:    !ARCH_DIR!
echo INCLUDE:     !INCLUDE_DIR!
echo LIB:         !LIB_DIR!
echo Output:      !OUTPUT!
echo Checking x64: 
if exist "!OPENCV_ROOT!\x64" (echo FOUND x64) else (echo NOT FOUND x64)
echo Checking arm64:
if exist "!OPENCV_ROOT!\arm64" (echo FOUND arm64) else (echo NOT FOUND arm64)
echo.

cl /LD /nologo /O2 /utf-8 /W3 /EHsc /MD /utf-8 ^
    /Fe:"!OUTPUT!" ^
    /I"!INCLUDE_DIR!" ^
    "!SOURCE!" ^
    "!LIB_DIR!\opencv_core4130.lib" ^
    "!LIB_DIR!\opencv_imgproc4130.lib" ^
    "!LIB_DIR!\opencv_photo4130.lib" ^
    "!LIB_DIR!\opencv_features2d4130.lib" ^
    "!LIB_DIR!\opencv_highgui4130.lib" ^
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
