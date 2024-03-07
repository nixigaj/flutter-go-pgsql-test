@echo off

if "%1" == "copy-static" goto :copy-static
if "%1" == "build" goto :build
if "%1" == "build-release" goto :build-release
if "%1" == "copy-static" goto :copy-static
if "%1" == "run" goto :run
if "%1" == "clean" goto :clean

REM Default target
if "%1" == "" goto :build

echo Invalid target: %1
echo Usage: .\make.bat [copy-static^|build^|build-release^|copy-static^|run^|clean]
goto :eof

:copy-static
	if not exist "build" mkdir build
	if not exist "build\config.toml" xcopy /I /Q static\config.toml build\
	goto :eof

:build
	call :copy-static
	go build -o .\build\backend.exe .\cmd\
	goto :eof

:build-release
	call :copy-static
	go build -o .\build\backend.exe .\cmd\
	goto :eof

:run
	cd build
	.\backend.exe
	cd ..
	goto :eof

:clean
	rd /s /q build
	goto :eof
