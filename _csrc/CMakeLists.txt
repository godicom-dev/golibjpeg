cmake_minimum_required(VERSION 3.10)
project(golibjpeg LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# Platform detection
if(WIN32)
    add_definitions(-DWIN32)
    if(CMAKE_SIZEOF_VOID_P EQUAL 8)
        add_definitions(-DWIN64)
    endif()
endif()

# thorfdbg/libjpeg common settings
add_definitions(-DACCUSOFT_CODE=1 -DISO_CODE=1)
add_definitions(-DCHECK_LEVEL=0 -DDEBUG_LEVEL=0 -DBUILD_LIB)

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/libjpeg)

# ---- libjpeg sources ----
file(GLOB LIBJPEG_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/boxes/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/codestream/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/coding/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/colortrafo/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/control/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/dct/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/interface/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/io/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/marker/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/std/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/tools/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/upsampling/*.cpp
)

# Exclude cmd/ and vs*/ files
file(GLOB CMD_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/cmd/*.cpp)
list(REMOVE_ITEM LIBJPEG_SOURCES ${CMD_SOURCES})

file(GLOB VS_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/vs10.0/**/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/vs12.0/**/*.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg/vs15.0/**/*.cpp
)
list(REMOVE_ITEM LIBJPEG_SOURCES ${VS_SOURCES})

# ---- golibjpeg wrapper ----
set(GOLIBJPEG_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/golibjpeg/golibjpeg.cpp
)

add_library(golibjpeg SHARED
    ${LIBJPEG_SOURCES}
    ${GOLIBJPEG_SOURCES}
)

# MSVC-specific
if(MSVC)
    add_definitions(-D_CRT_SECURE_NO_WARNINGS -D_VS10_)

    set_target_properties(golibjpeg PROPERTIES
        WINDOWS_EXPORT_ALL_SYMBOLS ON
    )

    # Static CRT to avoid VC runtime dependency
    foreach(flag CMAKE_CXX_FLAGS CMAKE_CXX_FLAGS_DEBUG CMAKE_CXX_FLAGS_RELEASE)
        if(${flag} MATCHES "/MD")
            string(REGEX REPLACE "/MD" "/MT" ${flag} "${${flag}}")
        endif()
    endforeach()
endif()

# GCC/Clang: hide internal symbols, only export our API
if(CMAKE_COMPILER_IS_GNUCXX OR CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    set_target_properties(golibjpeg PROPERTIES
        CXX_VISIBILITY_PRESET hidden
        VISIBILITY_INLINES_HIDDEN ON
    )
endif()

# Handle _CRT_SECURE_NO_WARNINGS-like things on MSVC
if(MSVC)
    # Already handled above
elseif(UNIX)
    # No CRT deprecation warnings on POSIX
endif()

# Output naming
if(WIN32)
    set_target_properties(golibjpeg PROPERTIES
        PREFIX ""
        SUFFIX ".dll"
    )
elseif(APPLE)
    set_target_properties(golibjpeg PROPERTIES
        PREFIX "lib"
        SUFFIX ".dylib"
    )
else()
    set_target_properties(golibjpeg PROPERTIES
        PREFIX "lib"
        SUFFIX ".so"
    )
endif()
