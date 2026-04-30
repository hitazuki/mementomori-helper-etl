# CLAUDE.md

Claude Code 开发指南

## 项目简介

MMTH ETL - mementomori-helper的日志处理工具

## 开发命令

```bash
# 构建
go build -o mmth_etl .

# 运行
./mmth_etl [-output <输出目录>] [-lang <语言>] [-records] <日志文件路径>

# 示例
./mmth_etl ./logs/game_log.json
./mmth_etl -output ./data ./logs/game_log.json
./mmth_etl -lang tw ./logs/game_log.json
./mmth_etl -records=false ./logs/game_log.json

# 测试
go test ./... -v
```

## 项目协作记忆

- 后续新增或修改代码注释时，尽量使用中文；仅在需要满足 Go 导出注释、第三方约定或引用原文时保留英文。

## 命令行参数

| 参数 | 说明 | 默认值 |
| --- | --- | --- |
| `-output` | 输出目录路径 | `./data` |
| `-lang` | 日志语言 (en/tw/ja/ko/auto/dynamic) | `dynamic` |
| `-records` | 保留详细变动记录 | `true` |
| `<日志文件路径>` | 待处理的日志文件路径 | 必填 |

## 规范

### Git Commit

使用约定式提交（Conventional Commits）：

```text
<type>(<scope>): <subject>

<body>

<footer>
```

**type 类型：**

- `feat`: 新功能
- `fix`: 修复
- `docs`: 文档
- `style`: 格式（不影响代码运行的变动）
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

**示例：**

```text
feat(parser): add support for new log format

Add parser for v2 log format with timestamp prefix.

Closes #123
```

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 使用 `go vet` 静态分析
- 使用 `staticcheck` 代码检查
- 函数/结构体需添加文档注释
- 错误处理：优先返回错误而非 panic

**检查命令：**

```bash
gofmt -l .
go vet ./...
staticcheck ./...
```

### Markdown 规范

- 使用 VS Code 内置 Markdown 检查或 npx markdownlint-cli 检查
- 文件末尾保留一个空行
- 标题前后空行
- 代码块指定语言

## 目录结构

```text
├── *.go              # Go 源码
├── go.mod            # 模块定义
├── CLAUDE.md         # 本文件
├── README.md         # 项目说明
├── data/             # 输出数据（gitignore）
├── logs/             # 日志文件（gitignore）
├── scripts/          # 辅助脚本
└── .github/          # GitHub 配置
    ├── workflows/    # CI/CD 工作流
    │   └── release.yml
    └── cliff.toml    # Changelog 配置
```
