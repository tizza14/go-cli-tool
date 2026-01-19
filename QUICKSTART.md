# 🚀 快速入門指南

## 專案啟動

### 方法 1: 使用啟動腳本（最簡單）

```cmd
REM 雙擊 start.bat 或執行：
start.bat

REM 帶參數執行
start.bat hello --name "您的名字"
start.bat version
```

### 方法 2: 直接執行

```cmd
REM 如果已構建
.\go-cli-tool.exe --help

REM 從源碼運行
go run main.go --help
```

---

## 📋 可用功能

### 1️⃣ Hello 命令 - 問候訊息

```cmd
start.bat hello
start.bat hello --name 張三
start.bat hello --name 張三 --upper
```

### 2️⃣ Version 命令 - 版本資訊

```cmd
start.bat version
```

### 3️⃣ **Task 命令 - 自動化任務執行器** ⭐ 新功能

#### 快速開始

```cmd
REM 1. 創建範例配置
start.bat task init -f my-tasks.yaml --example

REM 2. 查看任務列表
start.bat task list -f my-tasks.yaml

REM 3. 驗證配置
start.bat task validate -f my-tasks.yaml

REM 4. 執行所有任務
start.bat task run -f my-tasks.yaml -v

REM 5. 執行特定任務
start.bat task run -f my-tasks.yaml --id hello
```

#### 使用現有範例

```cmd
REM 簡單範例（3個任務）
start.bat task run -f examples/simple-tasks.yaml -v

REM 演示範例
start.bat task run -f examples/demo-tasks.yaml -v

REM 完整範例（12個任務，包含 Git、Go 構建等）
start.bat task list -f examples/tasks.yaml
```

---

## 💡 自動化任務範例

### 範例 1: 簡單的問候流程

創建 `hello-tasks.yaml`:

```yaml
version: "1.0"

tasks:
  - id: greet
    name: "問候"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Write-Host 'Hello! 任務開始執行' -ForegroundColor Green"

  - id: show-time
    name: "顯示時間"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Get-Date"
    depends_on: [greet]
```

執行：
```cmd
start.bat task run -f hello-tasks.yaml -v
```

### 範例 2: 文件備份

創建 `backup-tasks.yaml`:

```yaml
version: "1.0"

defaults:
  timeout: 60s

tasks:
  - id: check-source
    name: "檢查源目錄"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Test-Path './data'"

  - id: create-backup-dir
    name: "創建備份目錄"
    type: command
    command: powershell
    args:
      - "-Command"
      - "New-Item -ItemType Directory -Path './backup' -Force"
    depends_on: [check-source]

  - id: copy-files
    name: "複製文件"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Copy-Item -Path './data/*' -Destination './backup/' -Recurse -Force"
    depends_on: [create-backup-dir]
    retry_count: 2

  - id: compress
    name: "壓縮備份"
    type: command
    command: powershell
    args:
      - "-Command"
      - "$date = Get-Date -Format 'yyyyMMdd-HHmmss'; Compress-Archive -Path './backup/*' -DestinationPath \"backup-$date.zip\" -Force"
    depends_on: [copy-files]
```

### 範例 3: Git 工作流

創建 `git-tasks.yaml`:

```yaml
version: "1.0"

tasks:
  - id: git-status
    name: "檢查狀態"
    type: command
    command: git
    args: [status, --short]

  - id: git-pull
    name: "拉取更新"
    type: command
    command: git
    args: [pull, origin, main]
    depends_on: [git-status]

  - id: show-log
    name: "顯示提交歷史"
    type: command
    command: git
    args: [log, --oneline, -n, "5"]
    depends_on: [git-pull]
```

### 範例 4: Go 專案構建流程

創建 `build-tasks.yaml`:

```yaml
version: "1.0"

defaults:
  timeout: 300s

tasks:
  - id: fmt
    name: "格式化代碼"
    type: command
    command: go
    args: [fmt, ./...]

  - id: vet
    name: "檢查代碼"
    type: command
    command: go
    args: [vet, ./...]
    depends_on: [fmt]

  - id: test
    name: "運行測試"
    type: command
    command: go
    args: [test, ./..., -v]
    depends_on: [vet]
    retry_count: 1

  - id: build
    name: "構建應用"
    type: command
    command: go
    args: [build, -o, go-cli-tool.exe]
    depends_on: [test]
```

---

## 🎯 常用命令速查

```cmd
REM === 基本命令 ===
start.bat --help                    # 查看幫助
start.bat version                   # 查看版本
start.bat hello                     # 問候訊息

REM === 任務命令 ===
start.bat task --help               # 任務幫助
start.bat task init -f tasks.yaml --example
start.bat task list -f tasks.yaml
start.bat task validate -f tasks.yaml
start.bat task run -f tasks.yaml -v
start.bat task run -f tasks.yaml --id task-name

REM === 開發命令 ===
go build -o go-cli-tool.exe         # 構建
go test ./... -v                    # 測試
go run main.go task run -f tasks.yaml
```

---

## 📚 文檔

- **任務自動化詳細文檔**: [docs/TASK_AUTOMATION.md](docs/TASK_AUTOMATION.md)
- **中文使用指南**: [README_TASK.md](README_TASK.md)
- **完整 README**: [README.md](README.md)
- **架構說明**: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- **API 文檔**: [docs/API.md](docs/API.md)

---

## 🔧 配置 PATH（可選）

如果您想在任何地方使用 `go` 命令，請將 Go 添加到系統 PATH：

1. 打開「系統屬性」→「環境變數」
2. 在「系統變數」中找到 `Path`
3. 添加：`C:\Program Files\Go\bin`
4. 重啟終端

---

## ❓ 常見問題

### Q: 如何創建自己的任務？
A: 執行 `start.bat task init -f my-tasks.yaml --example`，然後編輯該文件。

### Q: 任務執行失敗怎麼辦？
A: 
1. 使用 `-v` 參數查看詳細輸出
2. 檢查命令是否正確
3. 確保依賴的任務已成功執行
4. 使用 `--id` 單獨測試失敗的任務

### Q: 如何並行執行任務？
A: 使用 `-c` 參數：`start.bat task run -f tasks.yaml -c 3`（3個任務並行）

### Q: 支持哪些命令？
A: Windows 上支持所有 PowerShell 命令和已安裝的可執行文件。

---

## 🎓 下一步

1. ✅ 嘗試運行範例任務
2. ✅ 創建您的第一個任務配置
3. ✅ 閱讀詳細文檔了解高級功能
4. ✅ 將常用操作自動化為任務

**開始體驗：**
```cmd
start.bat task run -f examples/simple-tasks.yaml -v
```
