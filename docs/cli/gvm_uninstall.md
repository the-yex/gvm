## gvm uninstall

卸载已安装的 Go 版本

### 使用方法

```bash
gvm uninstall <version> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `version` | 要卸载的 Go 版本号 |

### 使用示例

```bash
# 卸载 Go 1.20.0
gvm uninstall go1.20.0

# 卸载最新安装的 1.20.x
gvm uninstall 1.20
```

### 注意事项

- 无法卸载当前正在使用的版本（需先切换到其他版本）
- 卸载后会释放磁盘空间（每个版本约 300-500MB）
- 卸载操作不可撤销，请谨慎操作

### 交互式卸载

在 `gvm list` 交互界面中，按 `x` 可卸载选中的版本：

```bash
gvm list
# 使用 ↑/↓ 选择版本，按 x 卸载
```

### 相关命令

- [gvm list](gvm_list.md) - 查看已安装版本
- [gvm install](gvm_install.md) - 安装版本