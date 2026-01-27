# Go 语言安装指南（Windows）

## 方法 1：自动下载安装（推荐）

### 步骤 1：下载 Go 安装包

打开 PowerShell 或命令提示符，执行以下命令下载 Go 1.21.6（推荐版本）：

```powershell
# 创建下载目录
mkdir C:\Go-Installer -Force

# 下载 Go 安装包（约 140MB）
Invoke-WebRequest -Uri "https://golang.google.cn/dl/go1.21.6.windows-amd64.msi" -OutFile "C:\Go-Installer\go1.21.6.windows-amd64.msi"

# 或者使用官方镜像
# Invoke-WebRequest -Uri "https://go.dev/dl/go1.21.6.windows-amd64.msi" -OutFile "C:\Go-Installer\go1.21.6.windows-amd64.msi"
```

### 步骤 2：运行安装程序

```powershell
# 启动安装程序
Start-Process "C:\Go-Installer\go1.21.6.windows-amd64.msi"
```

**安装向导说明：**
- 点击 "Next"
- 接受许可协议
- 安装位置保持默认：`C:\Program Files\Go`
- 点击 "Install"
- 完成后点击 "Finish"

### 步骤 3：配置环境变量（自动）

安装程序会自动配置以下环境变量：
- `GOROOT`: `C:\Program Files\Go`
- `Path`: 添加 `%GOROOT%\bin`

### 步骤 4：配置 GOPATH（可选但推荐）

```powershell
# 设置 GOPATH 到用户目录
[System.Environment]::SetEnvironmentVariable("GOPATH", "$env:USERPROFILE\go", "User")

# 添加 GOPATH\bin 到 PATH
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$env:USERPROFILE\go\bin*") {
    [System.Environment]::SetEnvironmentVariable("Path", "$currentPath;$env:USERPROFILE\go\bin", "User")
}
```

### 步骤 5：验证安装

**重新打开命令提示符或 PowerShell**，然后执行：

```bash
go version
```

应该显示：`go version go1.21.6 windows/amd64`

---

## 方法 2：手动下载安装

### 步骤 1：手动下载

打开浏览器，访问以下地址之一：

- **国内镜像（推荐）**：https://golang.google.cn/dl/
- **官方地址**：https://go.dev/dl/

下载文件：`go1.21.6.windows-amd64.msi`

### 步骤 2：手动安装

双击下载的 `.msi` 文件，按照安装向导完成安装。

### 步骤 3：手动配置环境变量

1. 右键点击"此电脑" -> "属性"
2. 点击"高级系统设置"
3. 点击"环境变量"
4. 在"系统变量"中找到 `Path`，点击"编辑"
5. 添加 `C:\Program Files\Go\bin`
6. 点击"确定"保存

### 步骤 4：配置 GOPATH

在"用户变量"中：
1. 点击"新建"
2. 变量名：`GOPATH`
3. 变量值：`C:\Users\你的用户名\go`
4. 在 `Path` 中添加：`%GOPATH%\bin`

---

## 方法 3：使用包管理器安装（Chocolatey）

如果您安装了 Chocolatey 包管理器，可以使用：

```powershell
# 以管理员身份运行 PowerShell
choco install golang -y
```

---

## 配置 Go 国内镜像（重要！）

由于网络原因，建议配置 Go 模块代理：

### 方法 A：临时配置

```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

### 方法 B：永久配置

在用户环境变量中添加：
- 变量名：`GOPROXY`
- 变量值：`https://goproxy.cn,direct`

或使用命令：

```powershell
[System.Environment]::SetEnvironmentVariable("GOPROXY", "https://goproxy.cn,direct", "User")
```

---

## 验证完整安装

执行以下命令验证：

```bash
# 1. 检查 Go 版本
go version

# 2. 检查 Go 环境
go env

# 3. 检查 GOPATH
go env GOPATH

# 4. 检查 GOPROXY
go env GOPROXY
```

---

## 常见问题

### Q1: 命令提示符找不到 go 命令

**解决方案：**
- 确保安装完成后重新打开了命令提示符
- 检查环境变量 Path 中是否包含 `C:\Program Files\Go\bin`
- 重启电脑

### Q2: 下载速度慢

**解决方案：**
- 使用国内镜像：https://golang.google.cn/dl/
- 或使用代理下载

### Q3: 安装后 go version 显示旧版本

**解决方案：**
- 卸载旧版本
- 重新安装新版本
- 确保环境变量指向正确的目录

---

## 安装完成后的下一步

1. 进入项目目录：
```bash
cd "d:/claude space/CharonOMS"
```

2. 下载项目依赖：
```bash
go mod download
```

3. 运行项目：
```bash
go run cmd/server/main.go
```

---

## 推荐的 IDE

- **VS Code**（推荐）：https://code.visualstudio.com/
  - 安装 Go 插件：在 VS Code 中搜索 "Go" 插件并安装

- **GoLand**：https://www.jetbrains.com/go/
  - JetBrains 专业的 Go IDE（付费，有免费试用）

---

## 快速命令参考

```bash
# 查看版本
go version

# 查看环境配置
go env

# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 运行程序
go run main.go

# 编译程序
go build

# 清理缓存
go clean -cache
```
