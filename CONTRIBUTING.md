# Contributing to GVM

感谢你考虑为 GVM 贡献代码！本文档将帮助你了解贡献流程。

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)

## 行为准则

本项目采用贡献者公约作为行为准则。参与本项目即表示你同意遵守其条款。请阅读 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) 了解详情。

## 如何贡献

### 报告 Bug

如果你发现了 bug，请通过 [GitHub Issues](https://github.com/the-yex/gvm/issues) 提交报告。提交时请包含：

- 清晰的标题和描述
- 重现步骤
- 预期行为与实际行为
- 你的环境信息（操作系统、Go 版本、GVM 版本）
- 相关的日志输出（使用 `gvm -v` 获取详细日志）

### 提出新功能

欢迎提出新功能建议！请在 Issue 中详细描述：

- 功能的用途和场景
- 预期的使用方式
- 可能的实现思路

### 提交代码

1. Fork 项目仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 编写代码并添加测试
4. 确保测试通过 (`go test ./...`)
5. 提交更改 (`git commit -m 'feat: add amazing feature'`)
6. 推送到分支 (`git push origin feature/amazing-feature`)
7. 创建 Pull Request

## 开发环境设置

### 前置要求

- Go 1.21+（推荐使用最新稳定版）
- Git

### 本地开发

```bash
# Clone your fork
git clone https://github.com/<your-username>/gvm.git
cd gvm

# Install dependencies
go mod download

# Build
go build -o gvm .

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### 项目结构

```
gvm/
├── cmd/                # CLI 命令实现 (Cobra)
├── pkg/                # 公共包 (version_manager 接口)
├── internal/           # 内部包
│   ├── consts/         # 全局常量定义
│   ├── core/           # 核心功能
│   ├── registry/       # 镜像源注册表
│   ├── version/        # 版本解析与管理
│   ├── tui/            # Bubbletea TUI 组件
│   ├── github/         # GitHub API 交互
│   ├── utils/          # 工具函数
│   └── prettyout/      # 输出美化
├── docs/               # 文档
├── demo/               # TUI 示例代码
└── dist/               # 构建产物 (生成)
```

## 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 函数命名清晰，避免缩写
- 添加必要的注释，特别是公共 API

### 测试要求

- 新功能必须包含单元测试
- 测试覆盖率目标：80%+
- 使用 table-driven 测试模式

```go
func TestVersionParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Version
        wantErr bool
    }{
        {"stable", "go1.21.0", Version{Major: 1, Minor: 21, Patch: 0}, false},
        {"beta", "go1.21beta1", Version{Major: 1, Minor: 21, Beta: 1}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## 提交规范

本项目采用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### 类型 (type)

| 类型 | 描述 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档更新 |
| `style` | 代码格式（不影响逻辑） |
| `refactor` | 重构（不添加功能或修复 bug） |
| `perf` | 性能优化 |
| `test` | 添加或修改测试 |
| `chore` | 构建、工具链等变更 |
| `ci` | CI/CD 配置变更 |

### 示例

```
feat(list): add version filtering with fuzzy search

- Add fuzzy matching algorithm
- Update TUI to support real-time filtering
- Add keyboard shortcut '/' for search mode

Closes #123
```

## Pull Request 流程

### PR 标题

PR 标题应遵循提交规范格式。

### PR 描述

请包含以下内容：

1. **变更说明**：清楚描述本次 PR 的目的
2. **相关 Issue**：关联的 Issue 编号
3. **测试计划**：如何验证变更的正确性
4. **截图**：如果涉及 UI 变化，请附上截图

### Review 要求

- 至少需要一位维护者的 approval
- 所有 CI 检查必须通过
- 代码必须符合规范要求

### 合并策略

- 使用 squash merge 保持历史整洁
- PR 合并后会自动关闭相关 Issue

---

再次感谢你的贡献！如果有任何问题，欢迎在 Issue 中提问或通过邮件联系维护者。