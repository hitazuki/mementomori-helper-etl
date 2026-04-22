# 日志解析语言对照表

## 日志格式来源

mementomori-helper 输出的日志格式：`{Name}: {itemName}({itemRarity}) × {itemCount}`

- `{Name}` 标签来自 `ResourceStrings.Name`（本地化）
- `{itemName}` 来自 MB 文件的 `TextResourceTable`（根据语言设置）

**各语言日志示例**：

| 语言 | 日志示例                            |
|------|-------------------------------------|
| 英文 | `Name: Diamonds(None) × -100`       |
| 繁中 | `名称: 鑽石(None) × -100`           |
| 简中 | `名称: 钻石(None) × -100`           |
| 日文 | `名称: ダイヤ(None) × -100`         |
| 韩文 | `이름: 다이아(None) × -100`         |

## ResourceStrings.Name 对照

| 语言 | 值   |
|------|------|
| 英文 | Name |
| 繁中 | 名称 |
| 简中 | 名称 |
| 日文 | 名称 |
| 韩文 | 이름 |

## 物品名称

| 英文 (EN)       | 繁中 (TW)   | 简中 (CN)   | 日文 (JA)      | 韩文 (KO)      | 说明        |
| --------------- | ----------- | ----------- | -------------- | -------------- | ----------- |
| Diamonds        | 鑽石        | 钻石        | ダイヤ         | 다이아         | 钻石        |
| Rune Ticket     | 符石兌換券  | 符石兑换券  | ルーンチケット | 룬 티켓        | 饼干/符文票 |
| Upgrade Panacea | 強化秘藥    | 强化秘药    | 強化秘薬       | 강화의 비약    | 红水        |

## 洞穴相关（已验证）

| 英文 (EN) | 繁中 (TW) | 简中 (CN) | 日文 (JA) | 韩文 (KO) |
| --------- | --------- | --------- | --------- | --------- |
| Enter Cave of Space-Time | 进入 時空洞窟 | 进入 时空洞窟 | (待确认) | (待确认) |
| Cave of Space-Time Finished | 時空洞窟已完成 | 时空洞窟已完成 | (待确认) | (待确认) |

## 塔名称

| 英文 (EN)         | 繁中 (TW) | 简中 (CN) | 日文 (JA) | 韩文 (KO)  |
| ----------------- | --------- | --------- | --------- | ---------- |
| Tower of Infinity | 無窮之塔  | 无穷之塔  | 無窮の塔  | 무한의 탑  |
| Tower of Azure    | 憂藍之塔  | 忧蓝之塔  | 藍の塔    | 남청의 탑  |
| Tower of Crimson  | 業紅之塔  | 业红之塔  | 紅の塔    | 홍염의 탑  |
| Tower of Emerald  | 蒼翠之塔  | 苍翠之塔  | 翠の塔    | 비취의 탑  |
| Tower of Amber    | 流金之塔  | 流金之塔  | 黄の塔    | 황철의 탑  |

## 挑战结果（已验证）

| 英文 (EN) | 繁中 (TW) | 简中 (CN) | 日文 (JA) | 韩文 (KO) |
| --------- | --------- | --------- | --------- | --------- |
| triumphed | 勝利      | 胜利      | (待确认)  | (待确认)  |
| failed    | 敗北      | 败北      | (待确认)  | (待确认)  |
| Challenge | 挑战      | 挑战      | 挑戦      | 도전      |

---

## ETL 多语言支持

### 代码结构

```text
mmth-etl/
├── i18n/
│   ├── i18n.go       # 核心类型：Language, PatternSet, Manager
│   ├── patterns.go   # 多语言正则定义（EN/TW/JA/KO）
│   └── detector.go   # 语言检测器 + ScoreAccumulator
├── types/
│   └── regex.go      # 正则访问器（通过 i18n.Manager 获取）
├── parser/
│   ├── identify.go   # 日志类型识别
│   └── extract.go    # 数据提取（含塔名规范化）
├── processor.go      # 主处理流程（集成动态检测）
└── main.go           # 入口：-lang 参数处理
```

### 使用方法

```bash
# 默认使用动态检测（推荐）
./mmth_etl ./logs/app.log

# 指定固定语言
./mmth_etl -lang tw ./logs/app.log    # 繁体中文
./mmth_etl -lang en ./logs/app.log    # 英文

# 启动时一次性检测
./mmth_etl -lang auto ./logs/app.log
```

### 语言模式对比

| 模式       | 说明                                 | 性能   | 适用场景           |
|------------|--------------------------------------|--------|--------------------|
| dynamic    | 运行时持续检测，自动切换（默认）     | 优     | 通用场景           |
| auto       | 启动时检测一次，之后固定             | 优     | 日志语言确定不变   |
| 固定语言   | 全程使用指定语言正则                 | 最优   | 明确知道日志语言   |

### 动态语言检测

**实现原理：**

使用 `ScoreAccumulator` 增量累加各语言得分：

```go
// i18n/detector.go
type ScoreAccumulator struct {
    scores      map[Language]int    // 当前累计得分
    scoreWindow []Language          // 滑动窗口
}

// 每行调用，O(1) 开销
func (a *ScoreAccumulator) AddLine(line string) Language {
    lang, score := detector.DetectSingleLine(line)
    a.scores[lang] += score  // 增量累加
    // 滑动窗口：移除过期得分
    ...
}
```

**性能对比：**

| 方法               | 额外开销         | 总匹配次数（4语言） |
|--------------------|------------------|---------------------|
| 无切换             | 0                | 9                   |
| 动态检测（优化前） | O(W×F) 定期扫描  | ~16                 |
| 动态检测（优化后） | O(1) 累加        | ~9                  |
| 全语言匹配         | 0                | 36                  |

**参数配置：**

```go
// processor.go
type DynamicLanguageConfig struct {
    Enabled         bool  // 启用动态检测
    WindowSize      int   // 滑动窗口大小（默认 100）
    SwitchThreshold int   // 切换阈值（默认 5）
}
```

---

## 实际日志样本（繁体中文）

### 物品变动

```text
[2026-01-01 00:00:00] [PlayerName(Lv100)] 名称: 鑽石(None) × 100
[2026-01-01 00:00:00] [PlayerName(Lv100)] 名称: 符石兌換券(SR) × 2
[2026-01-01 00:00:00] [PlayerName(Lv100)] 名称: 強化秘藥(SR) × 14
```

### 洞穴日志

```text
[2026-01-01 00:00:00] [PlayerName(Lv100)] 进入 時空洞窟
[2026-01-01 00:00:00] [PlayerName(Lv100)] 時空洞窟已完成
```

### 塔挑战

```text
[2026-01-01 00:00:00] [PlayerName(Lv100)] 挑战 業紅之塔 800 层一次: 勝利
```

### 主线挑战

```text
[2026-01-01 00:00:00] [PlayerName(Lv100)] 挑战 36-13 boss 一次：敗北
```

---

生成时间: 2026-04-21
