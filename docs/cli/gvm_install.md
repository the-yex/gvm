## gvm install

安装指定的 Go 版本

### 使用方法

```bash
gvm install <version> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `version` | 要安装的 Go 版本号，支持多种格式 |

### 版本格式

GVM 支持灵活的版本指定方式：

```bash
gvm install 1.23        # 安装最新的 1.23.x 版本
gvm install 1.23.0      # 安装精确版本 1.23.0
gvm install go1.23.0    # 带 go 前缀的版本号
gvm install latest      # 安装最新稳定版本
```

### 使用示例

```bash
# 安装最新的 Go 1.21 版本
gvm install 1.21

# 安装指定的 Go 1.21.0 版本
gvm install go1.21.0

# 安装最新的稳定版本
gvm install latest

# 安装完成后自动切换
gvm install 1.22 && gvm use 1.22
```

### 相关命令

- [gvm list](gvm_list.md) - 查看可用版本
- [gvm use](gvm_use.md) - 切换版本
- [gvm uninstall](gvm_uninstall.md) - 卸载版本