# MMTH ETL

## 项目概述

MMTH ETL 是 [mementomori-helper](https://github.com/moonheart/mementomori-helper) 的日志转换工具集，用于解析和处理游戏日志数据。

目前支持的功能：

- **物品变动统计**：跟踪钻石/饼干/红水的获取和消耗情况，统一的数据模型和来源映射
- **时空洞窟追踪**：识别洞窟任务执行状态（已执行/已完成/异常），按角色和日期统计
- **战斗日志统计**：识别主线关卡和各种塔的挑战记录，统计尝试次数、通关状态和最后挑战时间
- **多语言支持**：支持英文/繁中/日文/韩文日志解析，可动态检测语言切换

## 功能特点

- **模块化架构**：按功能划分包结构，代码职责清晰
- **多语言支持**：支持 4 种语言日志解析，动态检测语言变化
- **统一数据模型**：钻石/饼干/红水使用相同的 `ChangeRecord` 和统计结构
- **流式日志解析**：从JSON日志文件中流式提取记录，支持GB级大文件
- **智能增量处理**：基于时间戳的二分查找定位，只处理新增日志
- **统一来源映射**：所有物品变动日志应用相同的来源映射规则
- **角色隔离**：每个角色的来源独立追踪
- **内存优化**：可选保留详细记录（默认开启），流式处理不缓存完整数据集

## 架构

```text
mmth-etl/
├── main.go              # 入口：命令行解析、流程协调
├── processor.go         # 日志处理：二分查找、流式读取、动态语言检测
├── i18n/                # 国际化模块
│   ├── i18n.go          # 核心类型：Language, PatternSet, Manager
│   ├── patterns.go      # 多语言正则定义（EN/TW/JA/KO）
│   └── detector.go      # 语言检测器 + ScoreAccumulator
├── parser/              # 解析模块
│   ├── parser.go        # 通用日志解析
│   ├── identify.go      # 日志类型识别
│   ├── source.go        # 来源映射
│   └── extract.go       # 记录提取（含塔名规范化）
├── aggregator/          # 聚合模块
│   ├── change.go        # 变动聚合器（钻石/饼干/红水通用）
│   ├── cave.go          # 洞窟聚合器
│   └── challenge.go     # 挑战聚合器
├── types/               # 类型定义
│   ├── base.go          # 基础类型（ChangeRecord, Stats）
│   ├── cave.go          # 洞窟类型
│   ├── challenge.go     # 挑战类型
│   └── regex.go         # 正则表达式（通过 i18n.Manager 获取）
├── storage/             # 存储模块
│   ├── checkpoint.go    # 检查点管理
│   └── file.go          # 文件读写
├── utils/               # 工具模块
│   └── time.go          # 时间工具
└── .github/             # GitHub 配置
    ├── workflows/       # CI/CD 工作流
    │   └── release.yml  # 发布构建
    └── cliff.toml       # Changelog 配置
```

## 使用方法

### 命令行参数

```bash
./mmth_etl [-output <输出目录>] [-lang <语言>] [-records] [-window <大小>] [-threshold <值>] <日志文件路径>
```

| 参数 | 说明 | 默认值 |
| ------ | ------ | -------- |
| `-output` | 输出目录路径 | `./data` |
| `-lang` | 日志语言 | `dynamic` |
| `-records` | 保留详细变动记录 | `true` |
| `-window` | 动态检测滑动窗口大小 | `100` |
| `-threshold` | 语言切换阈值 | `5`（窗口小时自动调整） |
| `<日志文件路径>` | 待处理的日志文件路径 | 必填 |

### 语言参数说明

| 值 | 说明 |
| ------ | ------ |
| `dynamic` | 动态检测语言变化，自动切换解析语言（默认） |
| `auto` | 启动时自动检测日志语言 |
| `en` | 英文日志 |
| `tw` | 繁体中文日志 |
| `ja` | 日文日志 |
| `ko` | 韩文日志 |

### 构建和运行

```bash
# 构建
go build -o mmth_etl .

# 构建时注入版本信息
go build -ldflags="-s -w -X main.Version=1.0.0" -o mmth_etl .

# 运行（默认动态检测语言，保留详细记录）
./mmth_etl ./logs/game_log.json

# 运行（繁体中文日志）
./mmth_etl -lang tw ./logs/game_log.json

# 运行（不保留详细记录，减少输出大小）
./mmth_etl -records=false ./logs/game_log.json

# 查看版本
./mmth_etl ./logs/game_log.json
# 输出: MMTH ETL v1.0.0
```

### 运行测试

```bash
go test ./... -v
```

## 配置

输出文件：

- `<output>/diamond_stats.json` - 钻石统计
- `<output>/cave_stats.json` - 洞窟统计
- `<output>/challenge_stats.json` - 战斗统计
- `<output>/rune_ticket_stats.json` - 饼干统计
- `<output>/upgrade_panacea_stats.json` - 红水统计
- `<output>/mmth_etl_state.json` - 检查点文件

## 日志处理规则

### 支持的日志格式

支持两种日志格式：

**Docker JSON 格式**（推荐）：

```json
{"log":"[2026-04-12 15:04:05] [角色名(Lv100)] 日志主体\n","time":"2026-04-12T07:04:05.123456789Z"}
```

**纯文本格式**：

```text
[2026-04-12 15:04:05] [角色名(Lv100)] 日志主体
```

时间戳说明：

- 日志内容中的时间戳为东8区本地时间
- 输出统一使用 ISO 8601 格式带时区偏移：`2026-04-12T15:04:05+08:00`

### 多语言日志格式

物品变动日志格式 `Name: ItemName(Quality) × N` 在各语言中的差异：

| 语言 | Name 标签 | 钻石 | 饼干 | 红水 |
| ------ | --------- | ------ | ------ | ------ |
| 英文 | Name | Diamonds | Rune Ticket | Upgrade Panacea |
| 繁中 | 名称 | 鑽石 | 符石兌換券 | 強化秘藥 |
| 日文 | 名称 | ダイヤ | ルーンチケット | 強化秘薬 |
| 韩文 | 이름 | 다이아 | 룬 티켓 | 강화의 비약 |

**示例对比**：

```text
# 英文
Name: Diamonds(None) × 100

# 繁体中文
名称: 鑽石(None) × 100
```

### 物品变动日志识别

| 日志格式 | 物品类型 | 英文示例 | 繁中示例 |
| ---------- | ---------- | ---------- | ---------- |
| `Name: Diamonds(None) × N` | 钻石 | `Name: Diamonds(None) × 100` | `名称: 鑽石 × 100` |
| `Name: Rune Ticket(Quality) × N` | 饼干 | `Name: Rune Ticket(SR) × 17` | `名称: 符石兌換券(SR) × 17` |
| `Name: Upgrade Panacea(Quality) × N` | 红水 | `Name: Upgrade Panacea(SR) × 38` | `名称: 強化秘藥(SR) × 38` |

**数量识别**：正数表示获取，负数表示消耗

### 洞穴日志识别

| 英文关键字 | 繁中关键字 | 状态 |
| ------------ | ------------ | ------ |
| `Enter Cave of Space-Time` | `进入 時空洞窟` | started |
| `Cave of Space-Time Finished` | `時空洞窟已完成` | finished |
| `KeyNotFoundException` | `KeyNotFoundException` | error |

### 挑战日志识别

| 日志格式 | 英文示例                                 | 繁中示例                                    |
| -------- | ---------------------------------------- | ------------------------------------------- |
| 主线关卡 | `Challenge 36-13 boss`                   | `挑战 36-13 boss`                           |
| 塔挑战   | `Challenge Tower of Crimson 800 layer`   | `挑战 業紅之塔 800 层一次`                  |

**塔名称对照**：

| 英文 | 繁中 |
| ------ | ------ |
| Tower of Infinity | 無窮之塔 |
| Tower of Azure | 憂藍之塔 |
| Tower of Crimson | 業紅之塔 |
| Tower of Emerald | 蒼翠之塔 |
| Tower of Amber | 流金之塔 |

**状态识别**：`triumphed`/`勝利` → 通关，`failed`/`敗北` → 失败

### 来源映射规则

所有物品变动日志统一应用来源映射：

| 原始日志 | 映射结果 | 说明 |
| ---------- | ---------- | ------ |
| `You have triumphed.` | `temple of illusions` | 幻想神殿通关 |
| `Gacha ... N times, ...` | `Gacha ...` | 抽卡日志简化 |
| `Open ... x N` | `Open ...` | 开箱日志简化 |

## 动态语言检测

当使用 `-lang dynamic` 时，ETL 会在运行时持续检测语言变化：

1. 使用 `ScoreAccumulator` 增量累加各语言得分
2. 每行日志识别时，同时检测语言特征并累加得分
3. 定期检查累计得分，当某语言显著领先时触发切换

**参数调优**：

| 参数 | 说明 | 推荐值 | 效果 |
| ------ | ------ | -------- | ------ |
| `-window 1` | 逐行检测 | 快速切换 | 每行都检查，适合语言频繁切换 |
| `-window 10` | 小窗口 | 快速响应 | 较快响应，可能轻微抖动 |
| `-window 100` | 大窗口（默认） | 稳定检测 | 减少抖动，适合语言稳定的日志 |
| `-threshold 1` | 低阈值 | 敏感切换 | 容易触发切换 |
| `-threshold 5` | 默认阈值 | 平衡 | 需要5分优势才切换 |
| `-threshold 10` | 高阈值 | 稳定切换 | 需要明显优势才切换 |

**窗口大小 vs 检查间隔**：

- 窗口大小：决定累积多少行的语言得分
- 检查间隔：`window/2`，即每多少行检查一次是否需要切换

**自适应阈值**：当窗口大小小于阈值时，阈值自动调整为 `max(window/2, 1)`，确保逐行检测时能正常切换。

**使用示例**：

```bash
# 逐行检测（窗口1，阈值自动调整为1）
./mmth_etl -window 1 -lang dynamic ./logs/app.log

# 快速响应切换（小窗口+低阈值）
./mmth_etl -window 10 -threshold 3 -lang dynamic ./logs/app.log

# 稳定检测（大窗口+高阈值）
./mmth_etl -window 100 -threshold 10 -lang dynamic ./logs/app.log
```

**性能优化**：

| 方法     | 额外开销 | 适用场景         |
| -------- | -------- | ---------------- |
| 固定语言 | 0        | 日志语言固定     |
| 动态检测 | O(1)     | 日志语言可能变化 |

## 输出格式

### 物品变动统计格式（钻石/饼干/红水统一）

```json
{
  "角色名": {
    "daily": {
      "2026-04-20": {
        "date": "2026-04-20",
        "records": [
          {
            "character": "角色名",
            "timestamp": "2026-04-20T10:30:00+08:00",
            "amount": 100,
            "source": "temple of illusions"
          }
        ],
        "gain": 100,
        "consume": 50,
        "net_change": 50,
        "sources": {
          "temple of illusions": {"gain": 17, "consume": 0},
          "Gacha 黒葬武具ガチャ": {"gain": 0, "consume": 5}
        }
      }
    },
    "weekly": { ... },
    "monthly": { ... },
    "total": {
      "gain": 500,
      "consume": 100,
      "net_change": 400,
      "sources": {}
    }
  }
}
```

**字段说明**：

| 字段 | 说明 |
| ------ | ------ |
| `records` | 详细变动记录列表（包含时间戳、数量、来源） |
| `gain` | 获取总量 |
| `consume` | 消耗总量 |
| `net_change` | 净变化（获取 - 消耗） |
| `sources` | 按来源分组的统计 |

使用 `-records=false` 可关闭详细记录输出以减少文件大小。

### 洞窟统计格式

```json
{
  "角色名": {
    "2026-04-20": {
      "date": "2026-04-20",
      "records": [...],
      "status": "finished"
    }
  }
}
```

### 挑战统计格式

```json
{
  "角色名": {
    "quest": {
      "43-6": {"level": "43-6", "attempts": 5, "success": true, "last_time": "..."}
    },
    "towers": {
      "Infinity": { "1840": {...} }
    }
  }
}
```

## 增量处理机制

1. 首次运行处理所有日志并创建检查点文件
2. 后续运行使用二分查找定位新增日志位置
3. 流式处理只读取新增日志
4. 新记录与现有统计数据合并

**检查点机制**：

- 使用秒级精度时间戳作为检查点
- 同一秒内的多条记录可能被跳过（优先保证不重复）
- 删除 `mmth_etl_state.json` 将重新处理所有日志

## 注意事项

- 日志文件格式支持 Docker JSON 和纯文本两种格式
- 时间戳使用东8区本地时间，输出带时区偏移
- 输出目录不存在时会自动创建
- 多语言日志解析时，塔名称统一输出为英文（如 `Crimson` 而非 `業紅之塔`）
