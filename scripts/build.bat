@echo off
set GOOS=windows
set GOARCH=amd64

echo Building Windows executable...
go build -o tracker.exe cmd\tracker\main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! 
    echo Executable created as tracker.exe
) else (
    echo Build failed!
)