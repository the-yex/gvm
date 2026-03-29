## gvm list

列出本地或远程可用的 Go 版本

### 使用方法

```bash
gvm list [flags]
```

### 选项

```
  -r, --remote             列出远程可用版本
  -m, --mirror string      临时指定镜像源（不保存到配置）
  -t, --type string        版本类型: stable | unstable | archived | all (默认 "all")
  -T, --timeout duration   HTTP 超时时间 (默认 5s)
  -h, --help               帮助信息
```

### 使用示例

#### 查看本地版本

```bash
# 列出已安装的版本（交互式 TUI）
gvm list

# 输出示例
┌─────────────────────────────────────────────┐
│  go1.21.0  ← current                         │
│  go1.22.0                                    │
│  go1.23.0                                    │
└─────────────────────────────────────────────┘
```

在交互界面中，可以使用键盘操作：

| 按键 | 功能 |
|------|------|
| `↑/k` | 向上移动 |
| `↓/j` | 向下移动 |
| `/` | 搜索过滤 |
| `u` | 使用选中版本 |
| `x` | 卸载选中版本 |
| `q` | 退出 |

#### 查看远程版本

```bash
# 列出远程所有可用版本
gvm list -r

# 仅显示稳定版本
gvm list -r -t stable

# 仅显示非稳定版本（beta, rc）
gvm list -r -t unstable

# 仅显示已归档版本
gvm list -r -t archived

# 使用国内镜像加速
gvm list -r -m https://mirrors.aliyun.com/golang/

# 设置超时时间（网络较慢时）
gvm list -r -T 30s

# 组合使用
gvm list -r -t stable -T 30s -m https://mirrors.ustc.edu.cn/golang/
```

在远程版本列表界面中，可以：

| 按键 | 功能 |
|------|------|
| `i` | 安装选中版本 |

### 版本类型说明

| 类型 | 说明 |
|------|------|
| `stable` | 正式发布的稳定版本（推荐生产使用） |
| `unstable` | 预发布版本（beta, rc），用于测试新功能 |
| `archived` | 已归档的旧版本，不再维护 |
| `all` | 显示所有版本（默认） |

### 镜像源推荐

国内用户建议设置镜像源以加速下载：

```bash
# 临时使用
gvm list -r -m https://mirrors.aliyun.com/golang/

# 永久设置
gvm config set mirror https://mirrors.aliyun.com/golang/
```

常用镜像源：

| 镜像 | URL |
|------|-----|
| 阿里云 | `https://mirrors.aliyun.com/golang/` |
| 中科大 | `https://mirrors.ustc.edu.cn/golang/` |
| 中国官方 | `https://golang.google.cn/dl/` |

### 相关命令

- [gvm install](gvm_install.md) - 安装版本
- [gvm use](gvm_use.md) - 切换版本
- [gvm config](gvm_config.md) - 设置默认镜像源