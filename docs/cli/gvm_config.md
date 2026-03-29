## gvm config

管理 GVM 配置项

### 使用方法

```bash
gvm config <command> [flags]
```

### 子命令

| 命令 | 说明 |
|------|------|
| `gvm config list` | 显示所有配置项 |
| `gvm config get <key>` | 获取指定配置项的值 |
| `gvm config set <key> <value>` | 设置配置项 |
| `gvm config unset <key>` | 删除配置项 |

### 主要配置项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `mirror` | Go 版本下载镜像源 | `https://go.dev/dl/` |
| `goroots` | 额外的 Go 安装目录列表 | 空 |

### 使用示例

```bash
# 查看所有配置
gvm config list

# 获取镜像源配置
gvm config get mirror

# 设置国内镜像源（推荐）
gvm config set mirror https://mirrors.aliyun.com/golang/

# 设置中科大镜像源
gvm config set mirror https://mirrors.ustc.edu.cn/golang/

# 删除自定义配置
gvm config unset custom_key
```

### 镜像源推荐

国内用户推荐使用以下镜像源以获得更快的下载速度：

| 镜像源 | URL |
|--------|-----|
| 阿里云 | `https://mirrors.aliyun.com/golang/` |
| 中科大 | `https://mirrors.ustc.edu.cn/golang/` |
| 华中科大 | `https://mirrors.hust.edu.cn/golang/` |
| 南京大学 | `https://mirrors.nju.edu.cn/golang/` |
| 中国官方 | `https://golang.google.cn/dl/` |

### 配置文件位置

配置文件存储在 `~/.gvm/config.yaml`。

### 相关子命令

- [gvm config list](gvm_config_list.md)
- [gvm config get](gvm_config_get.md)
- [gvm config set](gvm_config_set.md)
- [gvm config unset](gvm_config_unset.md)