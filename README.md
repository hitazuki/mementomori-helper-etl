# MMTH ETL

MMTH ETL 是 [mementomori-helper](https://github.com/moonheart/mementomori-helper)
的日志转换工具集，用于解析和处理游戏日志数据。

## 项目概述

目前支持的功能：

- **钻石记录处理**：跟踪游戏中钻石的获取和消耗情况，为每个角色提供每日、每周、
  每月和总计的统计信息

## 功能特点

- **日志解析**：从JSON日志文件中提取钻石交易记录
- **统计计算**：计算每日、每周、每月和总计的钻石统计信息
- **增量处理**：跟踪上次处理的时间戳，只处理新日志
- **来源追踪**：按来源统计钻石获取和消耗情况
- **角色隔离**：每个角色的钻石来源独立追踪
- **数据去重**：自动去重，支持重复运行

## 架构

应用程序遵循清晰的分离关注点架构：

1. **数据解析** - `log_parser.go` 处理从JSON日志中提取钻石交易记录
2. **统计分析** - `aggregator.go` 执行每日、每周、每月和总计的计算
3. **检查点管理** - `checkpoint.go` 管理处理器检查点状态（上次处理时间戳）
4. **类型定义** - `types.go` 定义数据结构和正则表达式模式
5. **主协调** - `main.go` 协整整个处理管道

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
├── main.go          # 主程序入口，接收命令行参数
├── types.go         # 数据结构和正则表达式模式定义
├── log_parser.go    # 从日志中提取钻石交易记录，处理来源追踪
├── aggregator.go    # 计算每日/每周/每月/总计统计
├── checkpoint.go    # 检查点文件管理（上次处理时间戳）
├── CLAUDE.md        # Claude Code 开发指南
├── README.md        # 项目文档
├── go.mod           # Go 模块定义
├── logs/            # 测试日志文件目录
│   └── test*.log    # 测试用日志文件
├── data/            # 输出数据目录
│   ├── diamond_stats.json          # 统计结果输出
│   └── mmth_etl_state.json         # 检查点文件
└── scripts/         # 工具脚本目录
    └── extract_by_date.py          # 按日期提取日志的脚本
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

- `<output>/diamond_stats.json` - 统计结果
- `<output>/mmth_etl_state.json` - 检查点文件

## 命名约定

### 文件命名（snake_case）

- `main.go` - 主程序入口
- `types.go` - 数据模型定义
- `log_parser.go` - 日志解析器
- `aggregator.go` - 聚合统计逻辑
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
| `Diamonds(None) × N` | 钻石获取记录 |
| `Diamonds(None) × -N` | 钻石消耗记录 |

### 来源追踪规则

- 钻石来源取自同一角色的最近非物品变动日志
- 不同角色之间的来源相互隔离
- 只使用时间戳大于上次处理的日志

## 输出格式

处理完成后，统计结果保存到 `data/diamond_stats.json`，格式如下：

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

## 增量处理机制

1. 首次运行时，处理所有日志并创建检查点文件
2. 后续运行时，从检查点文件读取上次处理的时间戳
3. 只处理时间戳大于上次处理的日志记录
4. 新记录与现有统计数据合并，自动去重
5. 更新检查点文件为最新处理时间

## 注意事项

- 确保日志文件格式正确，每行必须是有效的JSON
- 首次运行建议备份现有 `data/diamond_stats.json`（如果存在）
- 检查点文件 `data/mmth_etl_state.json` 存储上次处理时间，
  删除它将导致重新处理所有日志
- 程序使用 `timestamp|character|amount` 作为去重键，
  相同时间、角色、数量的记录视为重复

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
