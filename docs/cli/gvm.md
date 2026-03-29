## gvm

简单、快速、优雅的 Go 版本管理工具

### 简介

GVM 是一款 Go 版本管理工具，类似于 Node.js 的 `nvm` 或 Python 的 `pyenv`。它可以帮助你：

- 安装和管理多个 Go 版本
- 在不同版本之间快速切换
- 通过交互式界面轻松操作
- 使用国内镜像加速下载

非常适合需要在不同项目中使用不同 Go 版本的开发者。

### 安装

```bash
# GitHub 源（推荐）
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash

# Gitee 源（国内加速）
curl -sSL https://raw.githubusercontent.com/the-yex/gvm/main/install.sh | bash -s -- --source gitee
```

安装完成后重启终端，或执行 `source ~/.bashrc`（或 `~/.zshrc`）。

### 快速示例

```bash
# 查看远程可用版本
gvm list -r

# 安装 Go 1.21
gvm install 1.21

# 切换到 Go 1.21
gvm use 1.21

# 验证
go version
```

### 子命令

GVM 提供以下核心命令：

| 命令 | 说明 | 详细文档 |
|------|------|----------|
| [gvm list](gvm_list.md) | 列出 Go 版本 | 查看本地或远程版本 |
| [gvm install](gvm_install.md) | 安装 Go 版本 | 安装指定版本 |
| [gvm use](gvm_use.md) | 切换 Go 版本 | 切换到指定版本 |
| [gvm uninstall](gvm_uninstall.md) | 卸载 Go 版本 | 移除已安装版本 |
| [gvm new](gvm_new.md) | 创建新项目 | 使用指定版本创建项目 |
| [gvm upgrade](gvm_upgrade.md) | 升级 GVM | 更新到最新版本 |
| [gvm config](gvm_config.md) | 管理配置 | 查看/设置/删除配置 |

### 全局选项

```
  -h, --help      帮助信息
  -v, --verbose   详细输出模式
```

### 更多信息

- [README.md](../README.md) - 项目概述和完整指南
- [CONTRIBUTING.md](../CONTRIBUTING.md) - 贡献指南
- [GitHub Issues](https://github.com/the-yex/gvm/issues) - 问题反馈