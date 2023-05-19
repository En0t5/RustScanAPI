@echo off

REM 设置应用程序名称
SET APP_NAME=rustapi

REM 编译 Linux 版本
SET GOOS=linux
SET GOARCH=amd64
go build -o %APP_NAME%-linux-amd64

REM 编译 MacOS 版本
SET GOOS=darwin
SET GOARCH=amd64
go build -o %APP_NAME%-darwin-amd64

REM 编译 Windows 版本
SET GOOS=windows
SET GOARCH=amd64
go build -o %APP_NAME%-win-amd64.exe

echo All builds complete!
