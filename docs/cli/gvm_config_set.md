## gvm config set

设置配置项的值

### 使用方法

```bash
gvm config set <key> <value> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `key` | 配置项名称 |
| `value` | 配置项的值 |

### 使用示例

```bash
# 设置镜像源（国内推荐）
gvm config set mirror https://mirrors.aliyun.com/golang/

# 设置中科大镜像
gvm config set mirror https://mirrors.ustc.edu.cn/golang/

# 设置中国官方镜像
gvm config set mirror https://golang.google.cn/dl/
```

### 镜像源列表

| 镜像源 | URL | 说明 |
|--------|-----|------|
| 官方 | `https://go.dev/dl/` | 国际用户 |
| 中国官方 | `https://golang.google.cn/dl/` | 国内官方 |
| 阿里云 | `https://mirrors.aliyun.com/golang/` | 推荐 |
| 中科大 | `https://mirrors.ustc.edu.cn/golang/` | 推荐 |
| 华中科大 | `https://mirrors.hust.edu.cn/golang/` | 高校 |
| 南京大学 | `https://mirrors.nju.edu.cn/golang/` | 高校 |

### 相关命令

- [gvm config](gvm_config.md) - 配置管理概述
- [gvm config get](gvm_config_get.md) - 获取配置项
- [gvm config unset](gvm_config_unset.md) - 删除配置项