@echo off
setlocal enabledelayedexpansion

set VSROOT=D:\Apps\Microsoft Visual Studio 14.0
set VCVARS=%VSROOT%\VC\vcvarsall.bat
set BUILD_DIR=%CD%\build64
set INSTALL_DIR=%CD%\..\libs\windows_amd64
set EMBED_DIR=%CD%\..\native\libs

echo === Setting up VS 2015 x64 environment ===
call "%VCVARS%" x86_amd64
if %ERRORLEVEL% neq 0 (
    echo Failed to set up VS 2015 x64 environment
    exit /b 1
)

echo === Configuring with CMake ===
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
cd /d "%BUILD_DIR%"

cmake .. -G "Visual Studio 14 2015" -A x64
if %ERRORLEVEL% neq 0 (
    echo CMake configuration failed
    exit /b 1
)

echo === Building Release ===
cmake --build . --config Release
if %ERRORLEVEL% neq 0 (
    echo Build failed
    exit /b 1
)

echo === Installing DLL ===
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /Y "%BUILD_DIR%\Release\golibjpeg.dll" "%INSTALL_DIR%"
copy /Y "%BUILD_DIR%\Release\golibjpeg.pdb" "%INSTALL_DIR%" >nul 2>&1

rem Copy to Go embed location
if not exist "%EMBED_DIR%" mkdir "%EMBED_DIR%"
copy /Y "%BUILD_DIR%\Release\golibjpeg.dll" "%EMBED_DIR%\golibjpeg_amd64.dll"

echo === Done ===
dir "%INSTALL_DIR%"
