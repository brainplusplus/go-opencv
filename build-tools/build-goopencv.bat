@echo off
REM Build goopencv.dll from goopencv_abi.cpp + opencv-mobile static libs

setlocal enabledelayedexpansion

set VS_PATH=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools
set ACTIVATE=!VS_PATH!\Common7\Tools\VsDevCmd.bat

if not exist "!ACTIVATE!" (
    echo ERROR: Visual Studio BuildTools not found at !ACTIVATE!
    exit /b 1
)

call "!VS_PATH!\VC\Auxiliary\Build\vcvars64.bat" >nul 2>&1

set OPENCV_ROOT=D:\golang\go-opencv\build-tools\opencv-mobile-4.13.0-windows-vs2022\opencv-mobile-4.13.0-windows-vs2022
set INCLUDE_DIR=!OPENCV_ROOT!\x64\include
set LIB_DIR=!OPENCV_ROOT!\x64\x64\vc17\staticlib
set OUTPUT=D:\golang\go-opencv\dist\goopencv.dll
set SOURCE=D:\golang\go-opencv\backend\goopencv_abi.cpp

echo === Compiling goopencv.dll ===
echo Source:  !SOURCE!
echo Include: !INCLUDE_DIR!
echo Lib:     !LIB_DIR!
echo Output:  !OUTPUT!
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
    ole32.lib

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo === COMPILATION FAILED ===
    exit /b %ERRORLEVEL%
)

echo.
echo === SUCCESS ===
dir "!OUTPUT!"
