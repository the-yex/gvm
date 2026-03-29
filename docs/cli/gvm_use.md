## gvm use

切换到指定的 Go 版本

### 使用方法

```bash
gvm use <version> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `version` | 要使用的 Go 版本号 |

### 版本格式

支持与 `install` 命令相同的版本格式：

```bash
gvm use 1.23        # 使用已安装的 1.23.x 版本
gvm use go1.23.0    # 使用精确版本
```

### 使用示例

```bash
# 切换到 Go 1.21
gvm use 1.21

# 切换并验证
gvm use go1.22.0 && go version

# 切换到最新的已安装版本
gvm use latest
```

### 工作原理

`gvm use` 通过修改符号链接实现版本切换：

```
~/.gvm/go → ~/.gvm/sdk/go1.21.0
```

切换后，`GOROOT` 环境变量指向的目录内容会立即更新。

### 常见问题

**Q: 切换后 `go version` 没变化？**

请确保环境变量正确设置：
```bash
echo $GOROOT  # 应显示 ~/.gvm/go
echo $PATH    # 应包含 ~/.gvm/go/bin
```

**Q: 切换到一个未安装的版本？**

先使用 `gvm install` 安装，或使用 `gvm list -r` 在交互界面中直接安装。

### 相关命令

- [gvm install](gvm_install.md) - 安装版本
- [gvm list](gvm_list.md) - 查看已安装版本