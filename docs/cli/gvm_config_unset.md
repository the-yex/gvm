## gvm config unset

删除指定的配置项

### 使用方法

```bash
gvm config unset <key> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `key` | 要删除的配置项名称 |

### 使用示例

```bash
# 删除自定义配置项
gvm config unset custom_setting

# 恢复默认镜像源（删除自定义镜像配置后使用默认值）
gvm config unset mirror
```

### 注意

删除配置项后，该配置将使用默认值或不存在。

### 相关命令

- [gvm config](gvm_config.md) - 配置管理概述
- [gvm config set](gvm_config_set.md) - 设置配置项
- [gvm config list](gvm_config_list.md) - 查看所有配置