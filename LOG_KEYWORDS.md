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

## ResourceStrings.RewardMissionMsg 对照

helper 奖励任务日志来自 `ResourceStrings(.XX).resx` 的 `RewardMissionMsg` 模板，ETL 按模板识别可归因来源。

- 详细记录 JSONL 使用明细 ID：`TextResourceID * 1000000 + amount`，保留奖励数值特征，例如 `Guild(id=111)` + `2000` 为 `111002000`。
- 统计 JSON 使用聚合 ID：`TextResourceID * 1000000`，按任务类型合并，例如 Guild 聚合为 `111000000`。
- 前端展示使用聚合文案，不显示奖励数值，例如 `领取 Guild 任务奖励`。

| 任务 | helper 模板 | 日志示例 | 明细 Source ID 示例 | 统计聚合 Source ID | 前端展示（简中） | 前端展示（繁中） |
|------|-------------|----------|--------------------|-------------------|----------------|----------------|
| Daily | `Get {0} 's {1} Reward` / `领取 {0} 的 {1} 奖励` | `Get Daily 's 60 Reward` / `领取 Daily 的 60 奖励` | 23214000060 | 23214000000 | 领取 Daily 任务奖励 | 領取 Daily 任務獎勵 |
| Weekly | `Get {0} 's {1} Reward` / `领取 {0} 的 {1} 奖励` | `Get Weekly 's 80 Reward` / `领取 Weekly 的 80 奖励` | 23215000080 | 23215000000 | 领取 Weekly 任务奖励 | 領取 Weekly 任務獎勵 |
| Guild | `Get {0} 's {1} Reward` / `领取 {0} 的 {1} 奖励` | `Get Guild 's 2000 Reward` / `领取 Guild 的 2000 奖励` | 111002000 | 111000000 | 领取 Guild 任务奖励 | 領取 Guild 任務獎勵 |

其它语言示例：

| 语言 | Daily 示例 | Weekly 示例 | Guild 示例 |
|------|------------|-------------|------------|
| 日文 | `Daily の 60 の報酬を受け取る` | `Weekly の 80 の報酬を受け取る` | `Guild の 2000 の報酬を受け取る` |
| 韩文 | `일일의 60 보상을 수령합니다` | `주간의 80 보상을 수령합니다` | `Guild의 2000 보상을 수령합니다` |

## ResourceStrings.Auto refresh diamonds 对照

helper 自动刷新祈愿之泉任务时会输出 `ResourceStrings(.XX).resx` 中的自动刷新模板。每条匹配日志表示消耗 `20` 钻石，ETL 记录为 `amount: -20`，来源归属到游戏内置 `SourceIDFountainOfPrayers`（`140`），因此会和祈愿之泉钻石获取量在同一个来源下聚合。

| 语言 | 模板/示例 | Source ID |
|------|-----------|-----------|
| 英文 | `Current expected value: {0}. Today, {1}/{2} auto-refreshes completed. Refreshing now.` | 140 |
| 简中 | `当前期望值：{0}，今日已自动刷新 {1}/{2} 次，现在刷新。` | 140 |
| 日文 | `現在の期待値：{0}。本日自動リフレッシュ {1}/{2} 回完了。リフレッシュ中。` | 140 |
| 韩文 | `현재 기대값: {0}. 오늘 {1}/{2}회 자동 새로고침 완료. 지금 새로고침합니다.` | 140 |

`{0}` 为当前期望值，可能是整数或小数；`{1}/{2}` 为当日自动刷新次数进度。

## 物品名称

| 英文 (EN)       | 繁中 (TW)   | 简中 (CN)   | 日文 (JA)      | 韩文 (KO)      | 说明        |
| --------------- | ----------- | ----------- | -------------- | -------------- | ----------- |
| Diamonds        | 鑽石        | 钻石        | ダイヤ         | 다이아         | 钻石        |
| Rune Ticket     | 符石兌換券  | 符石兑换券  | ルーンチケット | 룬 티켓        | 饼干/符文票 |
| Upgrade Panacea | 強化秘藥    | 强化秘药    | 強化秘薬       | 강화의 비약    | 红水        |

## 洞穴相关（已验证）

| 英文 (EN) | 繁中 (TW) | 简中 (CN) | 日文 (JA) | 韩文 (KO) |
| --------- | --------- | --------- | --------- | --------- |
| Enter Cave of Space-Time | 进入 時空洞窟 | 进入 时空洞窟 | 時空の洞窟に入る | 시공의 동굴 입장 |
| Cave of Space-Time Finished | 時空洞窟已完成 | 时空洞窟已完成 | 時空の洞窟完了 | 시공의 동굴 완료 |

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
| triumphed | 勝利      | 胜利      | 勝利      | 승리      |
| failed    | 敗北      | 败北      | 敗北      | 패배      |
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
| dynamic    | 启动预热 + 运行时持续检测，自动切换（默认） | 优 | 通用场景           |
| auto       | 启动时检测一次，之后固定             | 优     | 日志语言确定不变   |
| 固定语言   | 全程使用指定语言正则                 | 最优   | 明确知道日志语言   |

### 动态语言检测

**实现原理：**

动态模式分三层处理语言：

1. 启动预热：从断点位置后抽样最多 `window` 条有效日志，先确定初始语言
2. 行级提示：遇到高置信单行特征时，临时使用该语言解析当前行
3. 滑动窗口：使用 `ScoreAccumulator` 增量累加各语言得分，稳定切换全局语言

`ScoreAccumulator` 只记录高置信语言得分：

```go
// i18n/detector.go
type ScoreAccumulator struct {
    scores      map[Language]int     // 当前累计得分
    scoreWindow []languageScore      // 滑动窗口
}

// 每行调用，O(1) 开销
func (a *ScoreAccumulator) AddLine(line string) Language {
    lang, score := detector.DetectSingleLineUnique(line)
    a.scores[lang] += score  // 增量累加
    // 滑动窗口：移除过期得分
    ...
}
```

**性能对比：**

| 方法               | 额外开销              | 总匹配次数（4语言） |
|--------------------|-----------------------|---------------------|
| 无切换             | 0                     | 9                   |
| 动态检测（优化前） | O(W×F) 定期扫描       | ~16                 |
| 动态检测（当前）   | 启动 O(window)，运行期 O(1) | ~9                  |
| 全语言匹配         | 0                     | 36                  |

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
