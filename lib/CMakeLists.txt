cmake_minimum_required(VERSION 3.10)
project(golibjpeg LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

if(WIN32)
    add_definitions(-DWIN32)
    if(CMAKE_SIZEOF_VOID_P EQUAL 8)
        add_definitions(-DWIN64)
    endif()
endif()

include(TestBigEndian)
test_big_endian(IS_BIG_ENDIAN)
if(IS_BIG_ENDIAN)
    add_definitions(-DJPG_BIG_ENDIAN=1)
endif()

set(LIBJPEG_DIR ${CMAKE_CURRENT_SOURCE_DIR}/libjpeg)

add_definitions(-DACCUSOFT_CODE=1 -DISO_CODE=1)
add_definitions(-DCHECK_LEVEL=0 -DDEBUG_LEVEL=0 -DBUILD_LIB)

include_directories(
    ${LIBJPEG_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/interface
)

file(GLOB LIBJPEG_SOURCES
    ${LIBJPEG_DIR}/boxes/*.cpp
    ${LIBJPEG_DIR}/codestream/*.cpp
    ${LIBJPEG_DIR}/coding/*.cpp
    ${LIBJPEG_DIR}/colortrafo/*.cpp
    ${LIBJPEG_DIR}/control/*.cpp
    ${LIBJPEG_DIR}/dct/*.cpp
    ${LIBJPEG_DIR}/interface/*.cpp
    ${LIBJPEG_DIR}/io/*.cpp
    ${LIBJPEG_DIR}/marker/*.cpp
    ${LIBJPEG_DIR}/std/*.cpp
    ${LIBJPEG_DIR}/tools/*.cpp
    ${LIBJPEG_DIR}/upsampling/*.cpp
    ${LIBJPEG_DIR}/cmd/*.cpp
)

list(REMOVE_ITEM LIBJPEG_SOURCES ${LIBJPEG_DIR}/cmd/main.cpp)

set(INTERFACE_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/interface/decode.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/interface/streamhook.cpp
)

set(CAPI_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/capi/golibjpeg.cpp
)

add_library(golibjpeg SHARED
    ${LIBJPEG_SOURCES}
    ${INTERFACE_SOURCES}
    ${CAPI_SOURCES}
)

if(MSVC)
    add_definitions(-D_CRT_SECURE_NO_WARNINGS -D_VS10_)

    set_target_properties(golibjpeg PROPERTIES
        WINDOWS_EXPORT_ALL_SYMBOLS ON
    )

    foreach(flag CMAKE_CXX_FLAGS CMAKE_CXX_FLAGS_DEBUG CMAKE_CXX_FLAGS_RELEASE)
        if(${flag} MATCHES "/MD")
            string(REGEX REPLACE "/MD" "/MT" ${flag} "${${flag}}")
        endif()
    endforeach()
endif()

if(CMAKE_COMPILER_IS_GNUCXX OR CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    set_target_properties(golibjpeg PROPERTIES
        CXX_VISIBILITY_PRESET hidden
        VISIBILITY_INLINES_HIDDEN ON
    )
endif()

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
