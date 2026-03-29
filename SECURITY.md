# Security Policy

## Supported Versions

我们积极维护以下版本的 GVM：

| 版本 | 支持状态 |
|------|----------|
| >= 1.2.x | ✅ 活跃维护 |
| < 1.2.0 | ❌ 不再支持 |

## Reporting a Vulnerability

如果你发现安全漏洞，请**不要**通过公开的 GitHub Issue 报告。请通过以下方式私下报告：

1. **发送邮件**：1003941268@qq.com
2. **使用 GitHub Security Advisory**：[提交安全报告](https://github.com/the-yex/gvm/security/advisories/new)

### 报告内容

请包含以下信息：

- 漏洞类型（如 XSS、命令注入等）
- 漏洞影响的版本
- 重现步骤
- 可能的修复建议

### 响应流程

1. 确认收到报告（24小时内）
2. 评估漏洞严重程度
3. 开发修复补丁
4. 发布安全版本
5. 公开披露（经报告者同意后）

### 安全最佳实践

使用 GVM 时，请遵循以下安全建议：

- 仅从官方渠道（GitHub Releases）下载 GVM
- 定期使用 `gvm upgrade` 更新到最新版本
- 不要在配置文件中存储敏感信息
- 检查下载的 Go 版本是否来自可信镜像源