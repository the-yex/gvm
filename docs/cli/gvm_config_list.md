## gvm config list

显示所有 GVM 配置项

### 使用方法

```bash
gvm config list [flags]
```

### 使用示例

```bash
# 显示所有配置
gvm config list

# 输出示例
mirror: https://mirrors.aliyun.com/golang/
goroots: []
```

### 配置文件

配置保存在 `~/.gvm/config.yaml`，也可以直接编辑此文件。

### 相关命令

- [gvm config](gvm_config.md) - 配置管理概述
- [gvm config get](gvm_config_get.md) - 获取单个配置项
- [gvm config set](gvm_config_set.md) - 设置配置项