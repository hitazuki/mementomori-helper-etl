# MMTH ETL

## 项目概述

MMTH ETL 是 [mementomori-helper](https://github.com/moonheart/mementomori-helper) 的日志转换工具集，用于解析和处理游戏日志数据。

目前支持的功能：

- **钻石记录处理**：跟踪游戏中钻石的获取和消耗情况，为每个角色提供每日、每周、每月和总计的统计信息
- **时空洞窟追踪**：识别洞窟任务执行状态（已执行/已完成/异常），按角色和日期统计
- **战斗日志统计**：识别主线关卡和各种塔的挑战记录，统计尝试次数、通关状态和最后挑战时间
- **物品日志统计**：跟踪饼干(Rune Ticket)和红水(Upgrade Panacea)的获取/消耗情况，支持来源追踪

## 功能特点

- **流式日志解析**：从JSON日志文件中流式提取钻石交易记录，支持GB级大文件
- **智能增量处理**：基于时间戳的二分查找定位，只处理新增日志，避免重复处理
- **多维度统计**：计算每日、每周、每月和总计的钻石统计信息
- **来源追踪**：按来源统计钻石获取和消耗情况（自动排除物品变动、挑战记录、错误日志）
- **时空洞窟追踪**：识别洞窟日志关键字，统计每日洞窟任务执行状态
- **战斗日志统计**：识别主线关卡和塔挑战记录，支持5种塔类型
- **物品日志统计**：跟踪饼干和红水的获取/消耗，支持来源映射（Gacha/Open日志简化显示）
- **角色隔离**：每个角色的钻石来源和洞窟状态独立追踪
- **内存优化**：可选保留详细记录（默认关闭以节省内存），流式处理不缓存完整数据集
- **自动目录创建**：输出目录不存在时自动创建
- **日志类型识别**：一次扫描识别日志类型，分发到对应处理函数

## 架构

应用程序遵循清晰的分离关注点架构：

1. **数据解析** - `log_parser.go` 处理日志解析、二分查找定位、流式读取、日志类型识别、来源映射
2. **钻石统计** - `aggregator.go` 执行每日、每周、每月和总计的聚合计算
3. **洞窟统计** - `cave_aggregator.go` 处理时空洞窟状态统计
4. **战斗统计** - `challenge_aggregator.go` 处理主线关卡和塔挑战统计
5. **物品统计** - `item_aggregator.go` 通用物品聚合器，`rune_ticket_aggregator.go` 和 `upgrade_panacea_aggregator.go` 处理饼干和红水统计
6. **检查点管理** - `checkpoint.go` 管理断点状态（上次处理时间戳）
7. **类型定义** - `types.go` 定义数据结构和正则表达式模式
8. **主协调** - `main.go` 协调整合处理管道，支持可配置输出目录

## 使用方法

### 命令行参数

```bash
./mmth_etl [-output <输出目录>] <日志文件路径>
```

| 参数 | 说明 | 默认值 |
| --- | --- | --- |
| `-output` | 输出目录路径 | `./data` |
| `<日志文件路径>` | 待处理的日志文件路径 | 必填 |

### 构建和运行

```bash
# 构建
go build -o mmth_etl .

# 运行（使用默认输出目录）
./mmth_etl ./logs/game_log.json

# 运行（指定输出目录）
./mmth_etl -output ./output ./logs/game_log.json

# 或者直接运行
go run . ./logs/game_log.json
```

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定测试文件
go test ./log_parser_test.go -v
```

## 文件结构

```text
mmth_etl/
├── main.go                      # 主程序入口，接收命令行参数，协调处理流程
├── types.go                     # 数据结构和正则表达式模式定义
├── log_parser.go                # 日志解析、二分查找定位、流式读取
├── aggregator.go                # 钻石统计聚合计算（日/周/月/总计）
├── cave_aggregator.go           # 时空洞窟状态统计
├── challenge_aggregator.go      # 战斗日志统计
├── item_aggregator.go           # 物品统计通用聚合器
├── rune_ticket_aggregator.go    # 饼干统计聚合器
├── upgrade_panacea_aggregator.go # 红水统计聚合器
├── checkpoint.go                # 断点状态管理
├── CLAUDE.md                    # Claude Code 开发指南
├── README.md                    # 项目文档
├── go.mod                       # Go 模块定义
├── logs/                        # 测试日志文件目录
│   └── test*.log                # 测试用日志文件
├── data/                        # 输出数据目录
│   ├── diamond_stats.json       # 钻石统计结果
│   ├── cave_stats.json          # 洞窟统计结果
│   ├── challenge_stats.json     # 战斗日志统计结果
│   ├── rune_ticket_stats.json   # 饼干统计结果
│   ├── upgrade_panacea_stats.json # 红水统计结果
│   └── mmth_etl_state.json      # 检查点文件
└── scripts/                     # 工具脚本目录
    └── extract_by_date.py       # 按日期提取日志的脚本
```

## 配置

输出目录通过命令行参数 `-output` 指定：

```bash
# 默认输出到 ./data 目录
./mmth_etl ./logs/game_log.json

# 指定自定义输出目录
./mmth_etl -output /path/to/output ./logs/game_log.json
```

输出文件：

- `<output>/diamond_stats.json` - 钻石统计结果
- `<output>/cave_stats.json` - 时空洞窟统计结果
- `<output>/challenge_stats.json` - 战斗日志统计结果
- `<output>/rune_ticket_stats.json` - 饼干统计结果
- `<output>/upgrade_panacea_stats.json` - 红水统计结果
- `<output>/mmth_etl_state.json` - 检查点文件

## 命名约定

### 文件命名（snake_case）

- `main.go` - 主程序入口
- `types.go` - 数据模型定义
- `log_parser.go` - 日志解析器
- `aggregator.go` - 钻石聚合统计逻辑
- `cave_aggregator.go` - 洞窟聚合统计逻辑
- `challenge_aggregator.go` - 战斗日志聚合统计逻辑
- `checkpoint.go` - 检查点/状态管理

### 结构体命名（PascalCase）

- `LogProcessor` - 主处理器
- `Aggregator` - 统计聚合器
- `DiamondRecord` - 钻石记录
- `DailyStats` / `WeeklyStats` / `MonthlyStats` / `TotalStats` - 统计结构体

### 变量/函数命名（camelCase）

- `inputLogPath` - 输入日志路径
- `outputJSONPath` - 输出JSON路径
- `stateFilePath` - 状态文件路径
- `checkpoint` - 上次处理时间戳
- `loadCheckpoint()` / `saveCheckpoint()` - 检查点操作

## 日志处理规则

### 支持的日志格式

日志文件中的每行应为JSON格式，包含 `log` 和 `time` 字段：

```json
{
  "log": "[2026-04-12 15:04:05] [角色名 (Lv100)] 日志主体",
  "time": "2026-04-12T15:04:05Z"
}
```

### 日志主体格式

日志主体必须符合 `"[时间] [名字] 日志主体"` 格式：

- 时间格式：`YYYY-MM-DD HH:MM:SS`
- 名字格式：`角色名 (Lv等级)`，程序会自动提取角色名

### 特殊日志类型处理

| 日志前缀 | 处理方式 |
| --- | --- |
| `Name:` | 物品变动日志，不作为钻石来源 |
| `Challenge` | 挑战记录，不作为钻石来源 |
| `OnError` | 错误日志，不作为钻石来源 |
| `Diamonds(None) × N` | 钻石获取记录 |
| `Diamonds(None) × -N` | 钻石消耗记录 |

### 时空洞窟日志识别

| 日志关键字 | 状态 |
| --- | --- |
| `Enter Cave of Space-Time` | 已执行 (started) |
| `Cave of Space-Time Finished` | 已完成 (finished) |
| `KeyNotFoundException` | 异常 (error) |

每日状态优先级：异常 > 已完成 > 未完成

### 战斗日志识别

| 日志格式 | 类型 | 示例 |
| --- | --- | --- |
| `Challenge X-Y boss ...` | 主线关卡 | `Challenge 43-6 boss one time：You have triumphed.` |
| `Challenge Tower of X N layer ...` | 塔挑战 | `Challenge Tower of Amber 1303 layer one time：You have triumphed.` |

**塔类型**：Infinity, Azure, Crimson, Emerald, Amber

**状态识别**：

- `triumphed` → 已通关
- `failed` → 未通关

**统计内容**：

- 尝试次数（attempts）
- 是否通关（success）
- 最后挑战时间（last_time）

### 物品日志识别

| 日志格式 | 物品类型 | 示例 |
| --- | --- | --- |
| `Name: Rune Ticket(Quality) × N` | 饼干 | `Name: Rune Ticket(SR) × 17` |
| `Name: Upgrade Panacea(Quality) × N` | 红水 | `Name: Upgrade Panacea(SR) × 38` |

**数量识别**：正数表示获取，负数表示消耗

### 来源映射规则

物品和钻石来源取自同一角色的最近非物品变动日志，并应用以下映射规则：

| 原始日志 | 映射结果 | 说明 |
| --- | --- | --- |
| `You have triumphed.` | `temple of illusions` | 幻想神殿通关 |
| `Gacha ... N times, ...` | `Gacha ...` | 抽卡日志简化，去掉次数和消耗 |
| `Open ... x N` | `Open ...` | 开箱日志简化，去掉数量 |

**示例**：
- `Gacha 黒葬武具ガチャ 5 times, used Gold×250000` → `Gacha 黒葬武具ガチャ`
- `Open Silver Sealed Chest x 2` → `Open Silver Sealed Chest`

### 来源追踪规则

- 钻石来源取自同一角色的最近非物品变动日志
- 不同角色之间的来源相互隔离
- 只使用时间戳大于上次处理的日志

## 输出格式

### 钻石统计格式

处理完成后，统计结果保存到 `diamond_stats.json`，格式如下：

```json
{
  "角色名": {
    "daily": {
      "2026-04-12": {
        "date": "2026-04-12",
        "gain": 100,
        "consume": 50,
        "net_change": 50,
        "records": [
          {
            "character": "角色名",
            "timestamp": "2026-04-12 15:04:05",
            "amount": 100,
            "type": "gain",
            "source": "任务奖励"
          }
        ],
        "sources": {
          "任务奖励": {"gain": 100, "consume": 0},
          "商店购买": {"gain": 0, "consume": 50}
        }
      }
    },
    "weekly": {
      "2026-W15": {
        "week": "2026-W15",
        "gain": 500,
        "consume": 300,
        "net_change": 200,
        "sources": {}
      }
    },
    "monthly": {
      "2026-04": {
        "month": "2026-04",
        "gain": 2000,
        "consume": 1500,
        "net_change": 500,
        "sources": {}
      }
    },
    "total": {
      "gain": 10000,
      "consume": 6000,
      "net_change": 4000,
      "sources": {}
    }
  }
}
```

**注意**：`records` 字段仅在启用 `keepRecords` 模式时出现（默认关闭以节省内存）。

### 时空洞窟统计格式

处理完成后，洞窟统计结果保存到 `cave_stats.json`，格式如下：

```json
{
  "角色名": {
    "2026-04-20": {
      "date": "2026-04-20",
      "records": [
        {"character": "角色名", "timestamp": "2026-04-20 10:30:00", "status": "started"},
        {"character": "角色名", "timestamp": "2026-04-20 10:45:00", "status": "finished"}
      ],
      "status": "finished"
    }
  }
}
```

**状态说明**：

- `started` - 已执行但未完成
- `finished` - 已完成
- `error` - 异常

### 战斗日志统计格式

处理完成后，战斗统计结果保存到 `challenge_stats.json`，格式如下：

```json
{
  "角色名": {
    "quest": {
      "43-6": {
        "level": "43-6",
        "attempts": 5,
        "success": true,
        "last_time": "2026-04-20 10:30:00"
      }
    },
    "towers": {
      "Infinity": {
        "1840": {
          "level": "1840",
          "attempts": 1,
          "success": true,
          "last_time": "2026-04-20 11:00:00"
        }
      },
      "Amber": {
        "1303": {
          "level": "1303",
          "attempts": 3,
          "success": false,
          "last_time": "2026-04-20 12:00:00"
        }
      }
    }
  }
}
```

**字段说明**：

- `quest` - 主线关卡统计
- `towers` - 塔挑战统计（按塔类型分组）
- `attempts` - 尝试次数
- `success` - 是否已通关
- `last_time` - 最后挑战时间

### 物品统计格式

处理完成后，物品统计结果保存到 `rune_ticket_stats.json` 和 `upgrade_panacea_stats.json`，格式如下：

```json
{
  "角色名": {
    "daily": {
      "2026-04-20": {
        "date": "2026-04-20",
        "gain": 17,
        "consume": 0,
        "net_change": 17,
        "sources": {
          "temple of illusions": {"gain": 17, "consume": 0},
          "Gacha 黒葬武具ガチャ": {"gain": 0, "consume": 5}
        }
      }
    },
    "weekly": {
      "2026-W16": {
        "week": "2026-W16",
        "gain": 100,
        "consume": 20,
        "net_change": 80,
        "sources": {}
      }
    },
    "monthly": {
      "2026-04": {
        "month": "2026-04",
        "gain": 200,
        "consume": 50,
        "net_change": 150,
        "sources": {}
      }
    },
    "total": {
      "gain": 500,
      "consume": 100,
      "net_change": 400,
      "sources": {
        "temple of illusions": {"gain": 300, "consume": 0},
        "Gacha 黒葬武具ガチャ": {"gain": 0, "consume": 80}
      }
    }
  }
}
```

**字段说明**：

- `gain` - 获取数量
- `consume` - 消耗数量
- `net_change` - 净变化（获取 - 消耗）
- `sources` - 按来源分组的统计（仅在 daily 和 total 中记录）

## 增量处理机制

1. 首次运行时，处理所有日志并创建检查点文件
2. 后续运行时，使用**二分查找**快速定位第一个时间戳大于上次处理的日志位置（约20次读取即可定位百MB文件）
3. 从定位位置开始流式处理，只读取新增日志
4. 新记录与现有统计数据合并
5. 更新检查点文件为最新处理时间

### 流式处理流程

- 使用256KB缓冲区逐行扫描，支持GB级大文件
- 时间戳二次校验，确保不处理已处理过的记录
- 即时聚合，记录处理后可立即被GC回收

## 注意事项

- 确保日志文件格式正确，每行必须是有效的JSON
- 首次运行建议备份现有 `data/diamond_stats.json`（如果存在）
- 检查点文件 `data/mmth_etl_state.json` 存储上次处理时间，删除它将导致重新处理所有日志
- 默认不保留详细记录以节省内存，如需保留可在 `main.go` 中将 `NewAggregator(false)` 改为 `NewAggregator(true)`
- 输出目录不存在时会自动创建

## 示例工作流

```bash
# 1. 构建程序
go build -o mmth_etl .

# 2. 首次运行（处理所有日志）
./mmth_etl ./logs/test2w-json.log

# 3. 查看输出
cat data/diamond_stats.json

# 4. 追加新日志后再次运行（只处理新日志）
./mmth_etl ./logs/test2w-json.log

# 5. 查看更新后的统计
cat data/diamond_stats.json
```
