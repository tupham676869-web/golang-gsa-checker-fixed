@echo off
chcp 65001 >nul
title TurboEngine Client - Golang
cd /d "%~dp0"

echo ==================================================
echo   Dang dong bo thu vien (go mod tidy)...
echo ==================================================
go mod tidy

echo.
echo ==================================================
echo   Dang khoi chay TurboEngine Client...
echo ==================================================
go run .

echo.
echo ==================================================
echo   Chuong trinh da ket thuc.
pause
