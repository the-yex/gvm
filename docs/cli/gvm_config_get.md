## gvm config get

获取指定配置项的值

### 使用方法

```bash
gvm config get <key> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `key` | 配置项名称 |

### 使用示例

```bash
# 获取镜像源配置
gvm config get mirror
# 输出: https://mirrors.aliyun.com/golang/

# 获取额外 Go 安装目录
gvm config get goroots
```

### 常用配置项

| 配置项 | 说明 |
|--------|------|
| `mirror` | Go 版本下载镜像源 |
| `goroots` | 额外的 Go 安装目录 |

### 相关命令

- [gvm config](gvm_config.md) - 配置管理概述
- [gvm config set](gvm_config_set.md) - 设置配置项