cmake_minimum_required(VERSION 3.20)
project(goopencv_wasm_backend LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

option(GOOPENCV_ENABLE_OPENCV "Link against OpenCV static libs" OFF)

add_library(goopencv_abi_obj OBJECT goopencv_abi.cpp)

target_compile_definitions(goopencv_abi_obj PRIVATE
  $<$<BOOL:${GOOPENCV_ENABLE_OPENCV}>:GOOPENCV_ENABLE_OPENCV=1>
)

add_executable(goopencv_wasm $<TARGET_OBJECTS:goopencv_abi_obj>)

if (GOOPENCV_ENABLE_OPENCV)
  if (NOT DEFINED OpenCV_INCLUDE_DIRS)
    message(FATAL_ERROR "GOOPENCV_ENABLE_OPENCV=ON requires OpenCV_INCLUDE_DIRS")
  endif()
  if (NOT DEFINED OpenCV_LIBRARIES)
    message(FATAL_ERROR "GOOPENCV_ENABLE_OPENCV=ON requires OpenCV_LIBRARIES")
  endif()
  target_include_directories(goopencv_abi_obj PRIVATE ${OpenCV_INCLUDE_DIRS})
  target_link_libraries(goopencv_wasm PRIVATE ${OpenCV_LIBRARIES})
endif()

if (CMAKE_SYSTEM_NAME STREQUAL "WASI")
  target_link_options(goopencv_wasm PRIVATE
    "-Wl,--no-entry"
    "-Wl,--no-gc-sections"
    "-Wl,--export=goopencv_mat_new"
    "-Wl,--export=goopencv_mat_delete"
    "-Wl,--export=goopencv_mat_rows"
    "-Wl,--export=goopencv_mat_cols"
    "-Wl,--export=goopencv_mat_type"
    "-Wl,--export=goopencv_mat_clone"
    "-Wl,--export=goopencv_imgproc_cvt_color"
    "-Wl,--export=goopencv_imgproc_resize"
    "-Wl,--export-memory"
  )
  # Build as bare wasm library without wasi-libc CRT
  target_link_options(goopencv_wasm PRIVATE "-nostartfiles")
endif()
