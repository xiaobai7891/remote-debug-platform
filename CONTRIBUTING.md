# 贡献指南

感谢你对本项目的关注！我们欢迎任何形式的贡献。

## 如何贡献

### 报告 Bug

1. 检查 [Issues](../../issues) 中是否已存在相同问题
2. 如果没有，创建新 Issue，请包含：
   - 问题描述
   - 复现步骤
   - 期望行为
   - 实际行为
   - 环境信息（OS、Go 版本、Node 版本等）

### 功能建议

1. 在 Issues 中提出你的想法
2. 描述功能的用途和价值
3. 如果可能，提供实现思路

### 提交代码

1. Fork 本仓库
2. 创建功能分支: `git checkout -b feature/your-feature`
3. 提交更改: `git commit -m 'Add some feature'`
4. 推送分支: `git push origin feature/your-feature`
5. 创建 Pull Request

## 开发环境

### 后端

```bash
cd backend
go mod tidy
go run main.go
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

## 代码规范

### Go

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 函数和方法需要注释

### Vue/JavaScript

- 使用 2 空格缩进
- 组件使用 PascalCase 命名
- Props 使用 camelCase

## 提交信息规范

```
<type>: <description>

[optional body]
```

Type:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

示例:
```
feat: 添加设备分组功能

- 支持按标签分组
- 支持批量操作
```

## 许可证

贡献的代码将采用 MIT 许可证。
