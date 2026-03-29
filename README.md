<h1 align="center">GVM - Go Version Manager</h1>

<p align="center">
  <strong>简单、快速、优雅的 Go 版本管理工具</strong>
</p>

<p align="center">
  <a href="https://github.com/the-yex/gvm/releases">
    <img src="https://img.shields.io/github/v/release/the-yex/gvm?style=flat-square" alt="GitHub Release">
  </a>
  <a href="https://github.com/the-yex/gvm/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/the-yex/gvm?style=flat-square" alt="License">
  </a>
  <a href="https://golang.org">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  </a>
  <a href="https://github.com/the-yex/gvm/issues">
    <img src="https://img.shields.io/github/issues/the-yex/gvm?style=flat-square" alt="Issues">
  </a>
</p>

<p align="center">
  <a href="#安装">安装</a> •
  <a href="#快速上手">快速上手</a> •
  <a href="#功能特性">功能特性</a> •
  <a href="#命令参考">命令参考</a> •
  <a href="#镜像源">镜像源</a> •
  <a href="docs/cli/gvm.md">文档</a> •
  <a href="CONTRIBUTING.md">贡献</a>
</p>

---

## 为什么选择 GVM？

如果你曾在不同项目中使用不同 Go 版本，你一定遇到过版本切换的烦恼。GVM 专为解决这个问题而设计：

| 特性 | GVM | 原生安装 |
|------|-----|---------|
| 多版本共存 | ✅ 轻松管理 | ❌ 手动折腾 |
| 版本切换 | ✅ 一条命令 | ❌ 修改 PATH |
| 交互式 TUI | ✅ 可视化操作 | ❌ 无 |
| 镜像加速 | ✅ 国内友好 | ⚠️ 需手动配置 |
| 项目隔离 | 🚧 开发中 | ❌ 无 |

> 💡 **适合人群**：需要在多个项目中切换 Go 版本的开发者、Go 语言学习者、需要测试不同版本兼容性的库开发者

## 安装

### 一键安装（推荐）

```bash
# GitHub 源（推荐）
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash

# Gitee 源（国内加速）
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --source gitee
```

安装完成后，重启终端或执行：

```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

### 手动安装

从 [Releases](https://github.com/the-yex/gvm/releases) 页面下载对应平台的压缩包：

```bash
# 解压到 ~/.gvm 目录
mkdir -p ~/.gvm && tar -xzf gvm*.tar.gz -C ~/.gvm

# 配置环境变量（添加到 ~/.bashrc 或 ~/.zshrc）
export GVM_HOME="${HOME}/.gvm"
export GOROOT="${GVM_HOME}/go"
export PATH="${GVM_HOME}:${GOROOT}/bin:$PATH"
```

### 支持的平台

| 平台 | 架构 |
|------|------|
| macOS | amd64, arm64 |
| Linux | 386, amd64, arm, arm64, s390x, riscv64 |

## 快速上手

### 1️⃣ 查看可用版本

```bash
# 查看本地已安装版本（交互式 TUI）
gvm list

# 查看远程所有可用版本
gvm list -r

# 按类型筛选
gvm list -r -t stable      # 仅稳定版
gvm list -r -t unstable    # 仅非稳定版
gvm list -r -t archived    # 已归档版本
```

### 2️⃣ 安装 Go 版本

```bash
# 安装指定版本（支持模糊匹配）
gvm install 1.23      # 自动匹配最新 1.23.x
gvm install go1.21.0  # 精确版本
gvm install latest    # 最新稳定版
```

### 3️⃣ 切换版本

```bash
# 切换到指定版本
gvm use 1.21

# 验证
go version
```

### 4️⃣ 交互式操作

在 `gvm list` 或 `gvm list -r` 的交互界面中：

| 按键 | 功能 |
|------|------|
| `↑/k` | 向上移动 |
| `↓/j` | 向下移动 |
| `/` | 搜索过滤 |
| `i` | 安装选中版本 |
| `u` | 使用选中版本 |
| `x` | 卸载选中版本 |
| `q` | 退出 |
| `?` | 查看帮助 |

![TUI Demo](docs/images/ls.png)

## 功能特性

### 🎯 版本管理

- **多版本共存**：安装任意数量的 Go 版本，互不干扰
- **秒级切换**：通过符号链接实现版本切换，无需等待
- **模糊匹配**：输入 `1.23` 自动匹配最新的 `1.23.x`
- **版本筛选**：按 stable/unstable/archived 类型分类查看

### 🖥️ 交互式 TUI

基于 [Bubbletea](https://github.com/charmbracelet/bubbletea) 构建的现代终端界面：

- 实时搜索过滤
- 可视化版本列表
- 进度条显示下载进度
- 一键安装/切换/卸载

![Install Progress](docs/images/install.png)

### 🌐 镜像加速

内置多个镜像源，解决国内下载慢的问题：

```bash
# 设置默认镜像
gvm config set mirror https://mirrors.aliyun.com/golang/

# 临时使用镜像（不保存配置）
gvm list -r -m https://mirrors.ustc.edu.cn/golang/
```

### 🔧 配置管理

```bash
gvm config list           # 查看所有配置
gvm config get mirror     # 获取配置值
gvm config set mirror URL # 设置配置
gvm config unset key      # 删除配置
```

### 📦 项目创建

```bash
gvm new myproject         # 使用当前版本创建项目
gvm new myproject -V 1.21 # 指定 Go 版本
gvm new myproject -m github.com/user/myproject  # 指定 module
```

## 命令参考

| 命令 | 描述 | 示例 |
|------|------|------|
| `gvm list` | 列出版本 | `gvm list -r -t stable` |
| `gvm install` | 安装版本 | `gvm install 1.23` |
| `gvm use` | 切换版本 | `gvm use go1.21.0` |
| `gvm uninstall` | 卸载版本 | `gvm uninstall 1.20` |
| `gvm new` | 创建项目 | `gvm new myapp -V 1.21` |
| `gvm upgrade` | 升级 GVM | `gvm upgrade` |
| `gvm config` | 管理配置 | `gvm config set mirror URL` |

详细命令说明请参阅 [命令文档](docs/cli/gvm.md)。

## 镜像源

| 镜像 | URL | 说明 |
|------|-----|------|
| 官方 | `https://go.dev/dl/` | 官方源，可能较慢 |
| 中国官方 | `https://golang.google.cn/dl/` | 国内官方镜像 |
| 阿里云 | `https://mirrors.aliyun.com/golang/` | 推荐 |
| 中科大 | `https://mirrors.ustc.edu.cn/golang/` | 推荐 |
| 华中科大 | `https://mirrors.hust.edu.cn/golang/` | 高校镜像 |
| 南京大学 | `https://mirrors.nju.edu.cn/golang/` | 高校镜像 |

## 工作原理

GVM 通过符号链接管理 Go 版本：

```
~/.gvm/
├── go          → 符号链接指向当前使用的版本
├── sdk/
│   ├── go1.21.0/   # Go 1.21.0 安装目录
│   ├── go1.22.0/   # Go 1.22.0 安装目录
│   └── go1.23.0/   # Go 1.23.0 安装目录
├── config.yaml     # GVM 配置文件
└── gvm             # GVM 二进制文件
```

环境变量：
- `GOROOT` → `~/.gvm/go`（当前版本的符号链接）
- `GOPATH` → `~/go`（默认，可自定义）
- `PATH` → 包含 `~/.gvm/go/bin`

## 开发路线图

| 版本 | 功能 | 状态 |
|------|------|------|
| v1.0 | 核心命令体系 | ✅ 已完成 |
| v1.1 | 交互式 TUI | ✅ 已完成 |
| v1.2 | `.gvmrc` 项目版本隔离 | 🚧 开发中 |
| v1.3 | `gvm doctor` 环境诊断 | 📋 计划中 |
| v1.4 | Shell 自动补全 | 📋 计划中 |
| v2.0 | 插件系统 | 💡 规划中 |

## 贡献

欢迎参与 GVM 的开发！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献流程。

### 开发者快速开始

```bash
# 克隆仓库
git clone https://github.com/the-yex/gvm.git && cd gvm

# 安装依赖
go mod download

# 本地构建
go build -o gvm .

# 运行测试
go test ./...
```

## 社区

- **问题反馈**：[GitHub Issues](https://github.com/the-yex/gvm/issues)
- **功能建议**：[GitHub Discussions](https://github.com/the-yex/gvm/discussions)
- **邮件联系**：1003941268@qq.com

## 致谢

GVM 的开发离不开以下开源项目：

- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [archiver](https://github.com/mholt/archiver) - 压缩包处理

## 许可证

本项目采用 [MIT 许可证](LICENSE)。

---

<p align="center">
  如果 GVM 对你有帮助，请给项目一个 ⭐ Star 支持一下！
</p>