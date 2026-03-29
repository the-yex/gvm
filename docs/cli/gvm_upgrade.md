## gvm upgrade

将 GVM 升级到最新版本

### 使用方法

```bash
gvm upgrade [flags]
```

### 说明

此命令会从 GitHub Releases 下载最新的 GVM 版本，并自动替换当前安装。

### 使用示例

```bash
# 升级到最新版本
gvm upgrade

# 查看当前版本
gvm --version
```

### 升级流程

1. 检查 GitHub Releases 获取最新版本
2. 下载对应平台的压缩包
3. 替换 `~/.gvm/gvm` 二进制文件
4. 清理临时文件

### 注意事项

- 升级过程中不影响已安装的 Go 版本
- 建议定期升级以获取最新功能和修复

### 相关命令

- [gvm config](gvm_config.md) - 配置镜像源影响下载速度