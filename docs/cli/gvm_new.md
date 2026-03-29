## gvm new

使用指定的 Go 版本创建新项目

### 使用方法

```bash
gvm new <project-name> [flags]
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `project-name` | 项目名称（将作为目录名和 module 名） |

### 选项

```
  -V, --version string   使用指定的 Go 版本创建项目（默认当前版本）
  -m, --module string    自定义 module 名称
  -h, --help             帮助信息
```

### 使用示例

```bash
# 使用当前 Go 版本创建项目
gvm new myproject

# 使用 Go 1.21 创建项目
gvm new myproject -V 1.21

# 自定义 module 名称
gvm new myproject -m github.com/username/myproject

# 组合使用
gvm new api-server -V 1.22 -m github.com/myorg/api-server
```

### 生成的项目结构

```bash
gvm new myproject
```

会创建以下结构：

```
myproject/
├── go.mod        # Module 定义文件
├── main.go       # 入口文件
└── .gitignore    # Git 忽略配置
```

### 相关命令

- [gvm use](gvm_use.md) - 切换 Go 版本
- [gvm install](gvm_install.md) - 安装 Go 版本