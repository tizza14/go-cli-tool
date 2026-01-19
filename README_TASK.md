# 🤖 自動化任務執行器使用指南

## 🚀 快速開始

### 1. 創建任務配置文件

```cmd
REM 創建空配置文件
start.bat task init -f my-tasks.yaml

REM 創建帶範例的配置文件
start.bat task init -f my-tasks.yaml --example
```

### 2. 驗證配置文件

```cmd
start.bat task validate -f my-tasks.yaml
```

### 3. 查看所有任務

```cmd
start.bat task list -f my-tasks.yaml
```

### 4. 執行任務

```cmd
REM 執行所有任務
start.bat task run -f my-tasks.yaml

REM 執行特定任務
start.bat task run -f my-tasks.yaml --id hello

REM 啟用詳細輸出
start.bat task run -f my-tasks.yaml -v

REM 並行執行（3個任務同時運行）
start.bat task run -f my-tasks.yaml -c 3
```

## 📋 任務配置範例

### 簡單範例

```yaml
version: "1.0"

defaults:
  timeout: 30s

tasks:
  - id: hello
    name: "問候訊息"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Write-Host '任務開始執行！'"

  - id: backup
    name: "備份文件"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Copy-Item -Path './data' -Destination './backup' -Recurse"
    depends_on:
      - hello
```

### 實用範例

#### 1️⃣ Git 操作自動化

```yaml
tasks:
  - id: git-pull
    name: "拉取最新代碼"
    type: command
    command: git
    args: [pull, origin, main]

  - id: git-status
    name: "檢查狀態"
    type: command
    command: git
    args: [status]
    depends_on: [git-pull]
```

#### 2️⃣ 構建和測試流程

```yaml
tasks:
  - id: clean
    name: "清理舊文件"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Remove-Item -Path './dist' -Recurse -Force -ErrorAction SilentlyContinue"

  - id: build
    name: "構建專案"
    type: command
    command: go
    args: [build, -o, app.exe]
    depends_on: [clean]
    timeout: 120s

  - id: test
    name: "運行測試"
    type: command
    command: go
    args: [test, ./..., -v]
    depends_on: [build]
    retry_count: 2
```

#### 3️⃣ 文件處理

```yaml
tasks:
  - id: create-dir
    name: "創建目錄"
    type: command
    command: powershell
    args:
      - "-Command"
      - "New-Item -ItemType Directory -Path './output' -Force"

  - id: process-files
    name: "處理文件"
    type: command
    command: powershell
    args:
      - "-File"
      - "./scripts/process.ps1"
    depends_on: [create-dir]
    env:
      INPUT_DIR: "./input"
      OUTPUT_DIR: "./output"

  - id: compress
    name: "壓縮結果"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Compress-Archive -Path './output/*' -DestinationPath 'result.zip'"
    depends_on: [process-files]
```

#### 4️⃣ 系統維護

```yaml
tasks:
  - id: disk-check
    name: "檢查磁碟空間"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Get-PSDrive C | Select-Object Used,Free"

  - id: clean-temp
    name: "清理臨時文件"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Remove-Item -Path $env:TEMP\\* -Recurse -Force -ErrorAction SilentlyContinue"

  - id: report
    name: "生成報告"
    type: command
    command: powershell
    args:
      - "-Command"
      - "Get-Date | Out-File -FilePath './maintenance.log' -Append"
    depends_on: [disk-check, clean-temp]
```

## 🎯 功能特性

### ✅ 任務依賴

任務可以依賴其他任務，確保執行順序：

```yaml
tasks:
  - id: task-a
    name: "任務 A"
    type: command
    command: powershell
    args: ["-Command", "Write-Host 'A'"]

  - id: task-b
    name: "任務 B"
    type: command
    command: powershell
    args: ["-Command", "Write-Host 'B'"]
    depends_on: [task-a]  # B 會在 A 之後執行
```

### ✅ 重試機制

任務失敗時自動重試：

```yaml
- id: unstable-task
  name: "不穩定的任務"
  type: command
  command: flaky-script.bat
  retry_count: 3  # 失敗時重試 3 次
```

### ✅ 超時控制

防止任務運行過長：

```yaml
- id: long-task
  name: "長時間任務"
  type: command
  command: long-process.exe
  timeout: 5m  # 5 分鐘後超時
```

### ✅ 環境變數

為任務設置自定義環境變數：

```yaml
- id: deploy
  name: "部署應用"
  type: command
  command: deploy.bat
  env:
    ENVIRONMENT: "production"
    API_KEY: "your-key"
```

### ✅ 工作目錄

設置任務的工作目錄：

```yaml
- id: build-frontend
  name: "構建前端"
  type: command
  command: npm
  args: [run, build]
  workdir: "./frontend"
```

## 💡 使用技巧

1. **從簡單開始**：先創建帶範例的配置文件，然後逐步修改
2. **先驗證再執行**：執行前總是先驗證配置文件
3. **使用詳細模式**：調試時加上 `-v` 參數查看詳細信息
4. **測試單個任務**：用 `--id` 參數測試特定任務
5. **合理設置超時**：避免任務無限期運行
6. **善用依賴關係**：確保任務按正確順序執行
7. **添加描述**：為每個任務添加清晰的描述

## 📦 實際應用場景

- ✅ **CI/CD 流程**：自動化構建、測試、部署
- ✅ **備份任務**：定期備份文件和數據庫
- ✅ **數據處理**：批量處理文件和數據轉換
- ✅ **系統維護**：清理日志、檢查狀態
- ✅ **代碼管理**：Git 操作自動化
- ✅ **環境設置**：自動配置開發環境

## 🔍 故障排除

### 任務立即失敗
- 檢查命令路徑和語法
- 驗證工作目錄是否存在
- 確保依賴任務已完成

### 任務超時
- 增加超時時間
- 檢查進程是否卡住
- 優化命令或腳本

### 找不到命令
- 確保命令在 PATH 中
- 使用完整路徑
- 在 Windows 上使用 powershell 執行內建命令

## 📚 更多信息

- 查看範例文件：`examples/simple-tasks.yaml` 和 `examples/tasks.yaml`
- 詳細文檔：`docs/TASK_AUTOMATION.md`
- 使用幫助：`start.bat task --help`

---

**開始使用：**
```cmd
start.bat task init -f my-tasks.yaml --example
start.bat task run -f my-tasks.yaml -v
```
