@echo off

REM ����Ӧ�ó�������
SET APP_NAME=rustapi

REM ���� Linux �汾
SET GOOS=linux
SET GOARCH=amd64
go build -o %APP_NAME%-linux-amd64

REM ���� MacOS �汾
SET GOOS=darwin
SET GOARCH=amd64
go build -o %APP_NAME%-darwin-amd64

REM ���� Windows �汾
SET GOOS=windows
SET GOARCH=amd64
go build -o %APP_NAME%-win-amd64.exe

echo All builds complete!
