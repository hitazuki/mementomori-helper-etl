# MMTH ETL

## 项目概述

MMTH ETL 是 [mementomori-helper](https://github.com/moonheart/mementomori-helper) 的日志转换工具集，用于解析和处理游戏日志数据。

目前支持的功能：

- **物品变动统计**：跟踪钻石/饼干/红水的获取和消耗情况，统一的数据模型和来源映射
- **时空洞窟追踪**：识别洞窟任务执行状态（已执行/已完成/异常），按角色和日期统计
- **战斗日志统计**：识别主线关卡和各种塔的挑战记录，统计尝试次数、通关状态和最后挑战时间

## 功能特点

- **模块化架构**：按功能划分包结构，代码职责清晰
- **统一数据模型**：钻石/饼干/红水使用相同的 `ChangeRecord` 和统计结构
- **流式日志解析**：从JSON日志文件中流式提取记录，支持GB级大文件
- **智能增量处理**：基于时间戳的二分查找定位，只处理新增日志
- **统一来源映射**：所有物品变动日志应用相同的来源映射规则
- **角色隔离**：每个角色的来源独立追踪
- **内存优化**：可选保留详细记录（默认关闭），流式处理不缓存完整数据集

## 架构

```text
mmth-etl/
├── main.go              # 入口：命令行解析、流程协调
├── processor.go         # 日志处理：二分查找、流式读取
├── parser/              # 解析模块
│   ├── parser.go        # 通用日志解析
│   ├── identify.go      # 日志类型识别
│   ├── source.go        # 来源映射
│   └── extract.go       # 记录提取
├── aggregator/          # 聚合模块
│   ├── change.go        # 变动聚合器（钻石/饼干/红水通用）
│   ├── cave.go          # 洞窟聚合器
│   └── challenge.go     # 挑战聚合器
├── types/               # 类型定义
│   ├── base.go          # 基础类型（ChangeRecord, Stats）
│   ├── cave.go          # 洞窟类型
│   ├── challenge.go     # 挑战类型
│   └── regex.go         # 正则表达式
├── storage/             # 存储模块
│   ├── checkpoint.go    # 检查点管理
│   └── file.go          # 文件读写
└── utils/               # 工具模块
    └── time.go          # 时间工具
```

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

日志文件每行应为JSON格式：

```json
{
  "log": "[2026-04-12 15:04:05] [角色名 (Lv100)] 日志主体",
  "time": "2026-04-12T15:04:05Z"
}
```

### 物品变动日志识别（统一格式）

所有物品变动日志格式统一为 `Name: ItemName(Quality) × N`：

| 日志格式 | 物品类型 | 示例 |
| --- | --- | --- |
| `Name: Diamonds(None) × N` | 钻石 | `Name: Diamonds(None) × 100` |
| `Name: Rune Ticket(Quality) × N` | 饼干 | `Name: Rune Ticket(SR) × 17` |
| `Name: Upgrade Panacea(Quality) × N` | 红水 | `Name: Upgrade Panacea(SR) × 38` |

**数量识别**：正数表示获取，负数表示消耗

### 来源映射规则

所有物品变动日志统一应用来源映射：

| 原始日志 | 映射结果 | 说明 |
| --- | --- | --- |
| `You have triumphed.` | `temple of illusions` | 幻想神殿通关 |
| `Gacha ... N times, ...` | `Gacha ...` | 抽卡日志简化 |
| `Open ... x N` | `Open ...` | 开箱日志简化 |

**示例**：
- `Gacha 黒葬武具ガチャ 5 times, used Gold×250000` → `Gacha 黒葬武具ガチャ`
- `Open Silver Sealed Chest x 2` → `Open Silver Sealed Chest`

### 洞穴日志识别

| 日志关键字 | 状态 |
| --- | --- |
| `Enter Cave of Space-Time` | started |
| `Cave of Space-Time Finished` | finished |
| `KeyNotFoundException` | error |

### 挑战日志识别

| 日志格式 | 类型 |
| --- | --- |
| `Challenge X-Y boss ...` | 主线关卡 |
| `Challenge Tower of X N layer ...` | 塔挑战 |

**塔类型**：Infinity, Azure, Crimson, Emerald, Amber

**状态识别**：`triumphed` → 通关，`failed` → 失败

## 输出格式

### 物品变动统计格式（钻石/饼干/红水统一）

```json
{
  "角色名": {
    "daily": {
      "2026-04-20": {
        "date": "2026-04-20",
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

## 注意事项

- 确保日志文件格式正确，每行必须是有效的JSON
- 删除 `mmth_etl_state.json` 将导致重新处理所有日志
- 输出目录不存在时会自动创建
