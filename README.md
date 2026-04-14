<h1 align="center">GVM - Go Version Manager</h1>

<p align="center">
  <strong>让 Go 多版本管理变得顺手、可视、可升级</strong>
</p>

<p align="center">
  在不同项目之间切 Go 版本，不用再手改 PATH，不用再猜当前环境，也不用再自己折腾升级脚本。
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
  <a href="#为什么值得用">为什么值得用</a> •
  <a href="#三十秒上手">三十秒上手</a> •
  <a href="#核心能力">核心能力</a> •
  <a href="#自升级">自升级</a> •
  <a href="#安装">安装</a> •
  <a href="#命令参考">命令参考</a> •
  <a href="docs/cli/gvm.md">文档</a> •
  <a href="CONTRIBUTING.md">贡献</a>
</p>

---

## 为什么值得用

GVM 不是只把 Go 解压到本地这么简单，它把开发者最常见的几个痛点都一起处理了：

- 多版本共存，随时切换，不污染系统 Go
- `gvm list` / `gvm list -r` 提供交互式 TUI，安装、使用、卸载都能直接操作
- `gvm install 1.23` 支持模糊匹配，自动找到最新的 `1.23.x`
- 内置镜像配置与远程版本缓存，国内环境更稳，重复查询更快
- `gvm upgrade` 可自升级，升级过程能看到当前版本、目标版本、平台、升级包和下载进度
- `install.sh` 默认自动安装最新 release，不需要每次手动改脚本里的版本号

如果你经常在这些场景里切版本，GVM 会很顺手：

- 同时维护多个 Go 项目
- 测试不同 Go 版本兼容性
- 学习新版本特性，但又不想动系统环境
- 想给团队一套更统一的 Go 安装和升级方式

## 三十秒上手

### 1. 安装 GVM

```bash
# GitHub 源
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash

# Gitee 源
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --source gitee
```

### 2. 查看远程版本

```bash
gvm list -r
```

### 3. 安装并切换

```bash
gvm install 1.23
gvm use 1.23
go version
```

### 4. 升级 GVM 自己

```bash
gvm upgrade
gvm --version
```

## 界面预览

远程版本列表：

![Remote Versions](docs/images/ls-r.png)

安装过程：

![Install Progress](docs/images/install.png)

## 核心能力

### 1. 多版本管理

- 多个 Go 版本可以同时存在于 `~/.gvm/sdk`
- 当前使用版本通过 `~/.gvm/go` 符号链接切换
- 切换几乎是瞬时完成，不需要重新配置环境

```bash
gvm list
gvm use 1.22
gvm uninstall 1.21.0
```

### 2. 模糊安装与版本筛选

不用每次记完整 patch 版本，也不用自己去翻 release 页面。

```bash
gvm install 1.23      # 自动匹配最新 1.23.x
gvm install go1.21.0  # 精确安装
gvm install latest    # 最新稳定版

gvm list -r -t stable
gvm list -r -t unstable
gvm list -r -t archived
```

### 3. 交互式 TUI

`gvm list` 和 `gvm list -r` 内置交互界面，适合不想反复敲命令的时候。

支持的体验包括：

- 上下移动、搜索过滤
- 一键安装 / 使用 / 卸载
- 下载进度、速度、状态反馈
- 键位提示与失败重试
- 已安装、当前使用、进行中状态可视化

常用按键：

| 按键 | 功能 |
|------|------|
| `↑/k` | 向上移动 |
| `↓/j` | 向下移动 |
| `/` | 搜索过滤 |
| `i` | 安装选中版本 |
| `u` | 使用选中版本 |
| `x` | 卸载选中版本 |
| `r` | 重试失败操作 |
| `q` | 退出 |

### 4. 镜像与缓存

GVM 不只是能换镜像，也会缓存远程版本列表，减少重复请求。

```bash
gvm config set mirror https://mirrors.aliyun.com/golang/
gvm list -r -m https://mirrors.ustc.edu.cn/golang/
gvm list -r --refresh
```

这对国内环境尤其有用：

- 首次拉取走网络
- 短时间重复查看版本会优先用缓存
- 需要强制更新时直接加 `--refresh`

### 5. 自升级

GVM 自己也能升级，而且不是“黑盒替换”。

```bash
gvm upgrade
```

升级界面会展示：

- 当前版本
- 目标版本
- 当前平台
- 匹配到的升级包名称
- 下载进度和已下载大小
- 当前步骤、完成步骤、失败步骤

升级只会替换 `~/.gvm/gvm` 二进制文件，不会动你已经安装好的 Go 版本目录。

### 6. 项目创建

```bash
gvm new myproject
gvm new myproject -V 1.21
gvm new myproject -m github.com/user/myproject
```

## 安装

### 一键安装

```bash
# 默认安装最新 release
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash

# 使用 Gitee 作为优先来源
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --source gitee

# 安装指定版本
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --version 1.2.2
```

安装脚本说明：

- 默认自动解析最新 release
- `--source gitee` 优先使用 Gitee，失败后回退 GitHub
- `--version x.y.z` 可以安装指定版本

安装完成后，重启终端或执行：

```bash
source ~/.bashrc
# 或
source ~/.zshrc
```

### 手动安装

从 [Releases](https://github.com/the-yex/gvm/releases) 页面下载对应平台的压缩包：

```bash
mkdir -p ~/.gvm && tar -xzf gvm*.tar.gz -C ~/.gvm

export GVM_HOME="${HOME}/.gvm"
export GOROOT="${GVM_HOME}/go"
export PATH="${GVM_HOME}:${GOROOT}/bin:$PATH"
```

### 支持平台

| 平台 | 架构 |
|------|------|
| macOS | amd64, arm64 |
| Linux | 386, amd64, arm, arm64, s390x, riscv64 |

## 命令参考

| 命令 | 描述 | 示例 |
|------|------|------|
| `gvm list` | 列出版本 | `gvm list -r -t stable` |
| `gvm install` | 安装版本 | `gvm install 1.23` |
| `gvm use` | 切换版本 | `gvm use go1.21.0` |
| `gvm uninstall` | 卸载版本 | `gvm uninstall 1.20` |
| `gvm new` | 创建项目 | `gvm new myapp -V 1.21` |
| `gvm upgrade` | 升级 GVM 自身 | `gvm upgrade` |
| `gvm config` | 管理配置 | `gvm config set mirror URL` |

详细命令说明请参阅 [docs/cli/gvm.md](docs/cli/gvm.md)。

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

```text
~/.gvm/
├── go          -> 当前使用版本的符号链接
├── sdk/
│   ├── go1.21.0/
│   ├── go1.22.0/
│   └── go1.23.0/
├── config.yaml
└── gvm
```

环境变量：

- `GOROOT` -> `~/.gvm/go`
- `GOPATH` -> `~/go`
- `PATH` -> 包含 `~/.gvm` 与 `~/.gvm/go/bin`

## 维护者发版

```bash
# 一键执行测试、打包、打 tag、push、创建 GitHub Release
./scripts/release.sh 1.2.2
```

现在的发布流程有两个特点：

- `scripts/release.sh` 会串起测试、构建、tag、push、Release 创建
- `build.sh` 会在构建时自动注入版本号，不需要手动改源码常量
- `install.sh` 会自动解析最新 release，发布成功后安装脚本立即生效

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

欢迎参与 GVM 的开发，详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

```bash
git clone https://github.com/the-yex/gvm.git
cd gvm
go mod download
go test ./...
go build -o gvm .
```

## 社区

- 问题反馈：[GitHub Issues](https://github.com/the-yex/gvm/issues)
- 功能建议：[GitHub Discussions](https://github.com/the-yex/gvm/discussions)
- 邮件联系：1003941268@qq.com

## 致谢

GVM 的开发离不开这些优秀项目：

- [Cobra](https://github.com/spf13/cobra)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Viper](https://github.com/spf13/viper)
- [archiver](https://github.com/mholt/archiver)

## 许可证

本项目采用 [MIT](LICENSE) 许可证。

---

<p align="center">
  如果 GVM 对你有帮助，欢迎点一个 Star 支持这个项目。
</p>
