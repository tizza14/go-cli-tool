@echo off
REM Go CLI Tool 啟動腳本

REM 設置 Go 環境
set PATH=C:\Program Files\Go\bin;%PATH%

REM 檢查是否已經構建
if not exist go-cli-tool.exe (
    echo 正在構建專案...
    go build -o go-cli-tool.exe
    if errorlevel 1 (
        echo 構建失敗！
        pause
        exit /b 1
    )
    echo 構建完成！
)

REM 啟動程式
echo.
echo ========================================
echo   Go CLI Tool
echo ========================================
echo.
go-cli-tool.exe %*
