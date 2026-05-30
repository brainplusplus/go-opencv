@echo off
REM Build goopencv.dll from goopencv_abi.cpp + opencv-mobile static libs
REM Works both locally and in CI (GitHub Actions windows-2022 runner)

setlocal enabledelayedexpansion

REM --- Detect repo root (parent of build-tools/) ---
set "REPO_ROOT=%~dp0.."

REM --- Setup MSVC (CI has it at a standard path, local has it at a specific path) ---
if defined VSCMD_VER (
    echo Already in VS Developer Command Prompt.
) else (
    set "VS_PATH=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools"
    set "VS_PATH_ALT=C:\Program Files\Microsoft Visual Studio\2022\Enterprise"
    set "VS_PATH_COMM=C:\Program Files\Microsoft Visual Studio\2022\Community"
    set "VS_PATH_PREVIEW=C:\Program Files\Microsoft Visual Studio\2022\Preview"

    if exist "!VS_PATH!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else if exist "!VS_PATH_ALT!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH_ALT!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else if exist "!VS_PATH_COMM!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH_COMM!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else if exist "!VS_PATH_PREVIEW!\VC\Auxiliary\Build\vcvars64.bat" (
        call "!VS_PATH_PREVIEW!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1
    ) else (
        echo ERROR: Visual Studio 2022 not found. Install BuildTools or run from VS Developer Command Prompt.
        exit /b 1
    )
)

REM --- Resolve paths relative to repo root ---
set "OPENCV_ROOT=%REPO_ROOT%\build-tools\opencv-mobile-4.13.0-windows-vs2022\opencv-mobile-4.13.0-windows-vs2022"
set "INCLUDE_DIR=%OPENCV_ROOT%\x64\include"
set "LIB_DIR=%OPENCV_ROOT%\x64\x64\vc17\staticlib"
set "OUTPUT=%REPO_ROOT%\dist\goopencv.dll"
set "SOURCE=%REPO_ROOT%\backend\goopencv_abi.cpp"

if not exist "%OPENCV_ROOT%" (
    echo ERROR: opencv-mobile SDK not found at %OPENCV_ROOT%
    echo Download from https://github.com/nihui/opencv-mobile/releases/latest/download/opencv-mobile-4.13.0-windows-vs2022.zip
    exit /b 1
)

if not exist "%REPO_ROOT%\dist" mkdir "%REPO_ROOT%\dist"

echo === Compiling goopencv.dll ===
echo Source:  %SOURCE%
echo Include: %INCLUDE_DIR%
echo Lib:     %LIB_DIR%
echo Output:  %OUTPUT%
echo.

cl /LD /nologo /O2 /utf-8 /W3 /EHsc /MD /utf-8 ^
    /Fe:"%OUTPUT%" ^
    /I"%INCLUDE_DIR%" ^
    "%SOURCE%" ^
    "%LIB_DIR%\opencv_core4130.lib" ^
    "%LIB_DIR%\opencv_imgproc4130.lib" ^
    "%LIB_DIR%\opencv_photo4130.lib" ^
    "%LIB_DIR%\opencv_features2d4130.lib" ^
    "%LIB_DIR%\opencv_highgui4130.lib" ^
    ole32.lib ^
    user32.lib ^
    gdi32.lib

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo === COMPILATION FAILED ===
    exit /b %ERRORLEVEL%
)

echo.
echo === SUCCESS ===
dir "%OUTPUT%"
