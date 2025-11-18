# GVM - Go Version Manager

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**GVM** 是一款集 **Go 版本管理** 与 **项目管理** 于一体的开发工具，类似 Node.js 的 `nvm` 和 Rust 的 `cargo`。作者结合多种包管理器的经验设计了这款工具，让你在同一台机器上轻松安装、管理和切换多个 Go 版本，非常适合需要在不同项目中使用不同 Go 版本的开发者。

---


- **Go** - 核心语言
- [**Cobra**](https://github.com/spf13/cobra) **v1.10.1** - 强大的现代 CLI 框架
- Go 标准库

---
- **版本管理**
    - `gvm list` – 列出本地或远程 Go 版本（支持交互式操作）
    - `gvm install` – 安装指定版本
    - `gvm use` – 切换 Go 版本
    - `gvm uninstall` – 卸载指定版本
    - `gvm upgrade` – 更新 GVM 本身
- **项目管理**
    - `gvm new` – 创建新项目，可指定 Go 版本与 module(后期期望指定模板初始化项目)
- **配置管理**
    - `gvm config` – 查看、设置和删除 GVM 配置项
---
##  相关截图
```shell
gvm ls   # 列举本地已安装的版本号

# 可以移动上下游标，可以过滤版本号,相关快捷键如下
#↑/k up • ↓/j down • / filter • x uninstall • u use • q quit • ? more
```
![gvm list](/docs/images/ls.png)
```shell
gvm ls -r  # 获取golang官网支持的所有版本号
```
![gvm list -r](/docs/images/ls-r.png)
![gvm list install](/docs/images/ls-install.png)

```shell
gvm install 1.23  # 也可以直接指定版本安装
```
![gvm install](/docs/images/install.png)
## 安装工具

### 安装方式

```bash
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash

# 如果没有科技访问github 可以使用gitee
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --source gitee
```

## 📋 快速上手

### 列出Go版本(当前已支持在列表页交互式安装使用和卸载)

```bash
# 列出本地已安装的Go版本
gvm list

# 列出远程可用的Go版本
gvm list -r

# 列出特定类型的Go版本（稳定版、非稳定版或归档版）
gvm list -r -t stable
gvm list -r -t unstable
gvm list -r -t archived
```

### 安装Go版本

```bash
# 安装特定版本的Go
gvm install go1.21
```

### 切换Go版本

```bash
# 切换到特定版本的Go
gvm use go1.21
```

### 卸载Go版本

```bash
# 卸载特定版本的Go
gvm uninstall go1.21
```
### 配置管理

```bash
# 查看配置
gvm config list

# 获取配置
gvm config get mirror

# 设置配置
gvm config set mirror https://golang.google.cn/dl/

# 删除配置
gvm config unset custom_setting
```

### 创建新项目

```bash
# 使用当前活动的Go版本创建新项目
gvm new myproject

# 使用指定版本号创建新项目
gvm new myproject -V 1.21.0

# 指定module创建项目
gvm new myproject -m github/xxx/myproject
```

### 配置管理

```bash
# 列出所有配置
gvm config list

# 获取特定配置
gvm config get mirror

# 设置配置
gvm config set mirror https://golang.google.cn/dl/

# 删除配置
gvm config unset custom_setting
```

## 命令参考

| 命令              | 描述         |
|-----------------|------------|
| `gvm list`      | 列出Go版本     |
| `gvm install`   | 安装Go版本     |
| `gvm use`       | 切换到特定Go版本  |
| `gvm uninstall` | 卸载Go版本     |
| `gvm new`       | 创建新Go项目    |
| `gvm upgrade`   | 升级最新的gvm版本 |
| `gvm config`    | 管理GVM配置    |

更详细的命令说明请参考[命令文档](docs/cli/gvm.md)。

## 项目结构

```
├── cmd/           # 命令行工具实现
├── docs/          # 文档
│   └── cli/       # 命令行文档
├── internal/      # 内部包
│   ├── consts/    # 常量定义
│   ├── registry/  # 版本注册表
│   ├── version/   # 版本管理
│   └── utils/     # 工具函数
└── pkg/           # 公共包
```
## 🧭 开发路线图

| 阶段      | 功能 | 状态 |
|---------|------|------|
| ✅ v1.0  | 基础命令体系 (list/install/use/uninstall/config) | 已完成 |
| 🚧 v1.2 | `.gvmrc` 项目版本隔离 | 开发中 |
| 🚧 v1.3 | `gvm doctor` 环境诊断工具 | 计划中 |
| 🧩 v1.4 | Shell 自动补全、项目模板系统 | 计划中 |
| 🧠 v2.0 | 插件系统与智能版本推荐 | 规划中 |
## 贡献

欢迎贡献代码、报告问题或提出改进建议！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

本项目采用MIT许可证 - 详情请参阅[LICENSE](LICENSE)文件。

## 联系方式

如有任何问题或建议，请通过以下方式联系我们：

- 项目维护者：[mortal](1003941268@qq.com)
- GitHub Issues：[https://github.com/the-yex/gvm/issues](https://github.com/the-yex/gvm/issues)