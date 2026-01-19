# Go CLI Tool 啟動腳本

# 設置 Go 環境
$env:PATH = "C:\Program Files\Go\bin;" + $env:PATH

# 檢查是否已經構建
if (-not (Test-Path "go-cli-tool.exe")) {
    Write-Host "正在構建專案..." -ForegroundColor Yellow
    go build -o go-cli-tool.exe
    if ($LASTEXITCODE -ne 0) {
        Write-Host "構建失敗！" -ForegroundColor Red
        Read-Host "按 Enter 繼續..."
        exit 1
    }
    Write-Host "構建完成！" -ForegroundColor Green
}

# 啟動程式
Write-Host ""
Write-Host "========================================"
Write-Host "  Go CLI Tool"
Write-Host "========================================"
Write-Host ""

.\go-cli-tool.exe $args
