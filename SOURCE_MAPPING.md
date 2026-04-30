# 来源映射关系

## Source ID 分配规则

| 范围 | 说明 |
|------|------|
| 0 | 未知/未匹配的来源 |
| 1-99999 | 游戏 TextResource ID |
| 100000+ | helper 自定义 ID |
| `TextResourceID * 1000000 + amount` | helper `RewardMissionMsg` 明细来源 ID，保留游戏文本 ID 与奖励数值特征 |
| `TextResourceID * 1000000` | helper `RewardMissionMsg` 统计聚合来源 ID，按任务类型聚合展示 |

---

## 可识别来源列表

### 游戏内置（部分匹配）

| Source ID | 别名 | 英文 | 繁中 | 简中 | 日文 | 韩文 |
|-----------|------|------|------|------|------|------|
| 140 | Fountain of Prayers | Fountain of Prayers: | 祈願之泉: | 祈愿之泉: | 祈りの泉: | 기원의 샘: |
| 67 | Open | Open | 開啟 | 开启 | 開く | 열기 |
| 719 | Login Bonus | Login | 簽到獎勵: | 签到奖励: | ログイン | ログ인 |
| 138 | Tower of Infinity | Tower of Infinity: | 無窮之塔: | 无穷之塔: | 無窮の塔: | 무한의 탑: |
| 2766 | Temple of Illusions | You have triumphed. | 勝利 | 胜利 | 勝利しました | 승리했습니다. |
| 21308 | Presents Box | Presents Box Claim All | 禮物箱 | 礼物箱 | プレゼントボックス | 선물 상자 |
| 21332 | Monthly Boost | Monthly Boost Already Claimed | 每月強化組合包 | 每月强化组合包 | 月間ブースト | 월간 부스트 |
| 3331 | Total Logins This Month | Total Logins This Month: | 本月累計簽到天數： | 本月累计签到天数： | 今月の合計ログイン日数： | 이번 달 보상 수령: |
| 23277 | World Player Clears | A player in your World | 本世界首次有玩家 | 本世界首次有玩家 | ワールド内のプレイヤーが初めて | 월드 내 플레이어가 최초로 |

### helper自定义（RewardMissionMsg 模式）

| 任务 | TextResource ID | 明细 Source ID 示例 | 统计聚合 Source ID | 前端展示（简中） | 前端展示（繁中） |
|------|-----------------|--------------------|-------------------|----------------|----------------|
| Daily Mission Reward | 23214 | 23214000060 (`23214 * 1000000 + 60`) | 23214000000 | 领取 Daily 任务奖励 | 領取 Daily 任務獎勵 |
| Weekly Mission Reward | 23215 | 23215000080 (`23215 * 1000000 + 80`) | 23215000000 | 领取 Weekly 任务奖励 | 領取 Weekly 任務獎勵 |
| Guild Mission Reward | 111 | 111002000 (`111 * 1000000 + 2000`) | 111000000 | 领取 Guild 任务奖励 | 領取 Guild 任務獎勵 |
| Main Mission Reward | 23213 | 23213 | 23213 | 领取 Main 奖励 | 領取 Main 獎勵 |

### helper自定义（自建 ID）

| Source ID | 别名 | 英文 | 中文 | 日文 | 韩文 |
|-----------|------|------|------|------|------|
| 100002 | Auto Buy Store Items | Auto Buy Store Items | 自动购买商城物品 | 自動購入ストアアイテム | 자동으로 상점 아이템 구매 |
| 100004 | Missions Claim All | Missions Claim All | 见下方多文本映射 | 见下方多文本映射 | 见下方多文本映射 |
| 100005 | Gacha | Gacha | 抽卡 | ガチャ | 가챠 |
| 111002000 | Guild Mission Reward（明细 ID 示例） | Get Guild 's 2000 Reward | 領取 Guild 的 2000 獎勵 | 领取 Guild 的 2000 奖励 | Guild の 2000 の報酬を受け取る | Guild의 2000 보상을 수령합니다 |

**Guild Mission Reward ID 规则**：

- `Guild` 来自游戏 `TextResource`，ID 为 `111`（`[CommonGuildLabel]`）
- `2000` 来自 helper `RewardMissionMsg` 日志参数
- 明细复合 ID 计算：`111 * 1000000 + 2000 = 111002000`
- 统计聚合 ID 计算：`111 * 1000000 = 111000000`
- 同类日志可按相同规则扩展，例如 `Guild` + `3000` 为 `111003000`

**RewardMissionMsg 复合 ID 规则**：

- `Daily` 来自游戏 `TextResource`，ID 为 `23214`（`[MissionTabName2]`），例如 `Daily` + `60` 为 `23214000060`
- `Weekly` 来自游戏 `TextResource`，ID 为 `23215`（`[MissionTabName3]`），例如 `Weekly` + `80` 为 `23215000080`
- `Main` 当前仍保留裸游戏 `TextResource` ID `23213`
- 这些细化来源写入详细记录 JSONL 时，`source` 使用 `id:<source_id>` 形式，例如 `id:23214000060`
- 统计 JSON 的 `sources` 按任务类型聚合，不保留奖励数值：
  - Daily：`id:23214000000`
  - Weekly：`id:23215000000`
  - Guild：`id:111000000`
- 前端展示同样按聚合来源显示，不显示奖励数值，例如 `领取 Daily 任务奖励`

**Missions Claim All (100004) 多文本映射**：

| 来源 | 英文 | 繁中 | 简中 | 日文 | 韩文 |
| ------ | ------ | ------ | ------ | ------ | ------ |
| 游戏 TextResource | You have no more challenges left. | 剩餘挑戰次數不足 | 剩余挑战次数不足 | 残り挑戦回数がありません | - |
| 游戏 TextResource | Cave of Space-TimeFinished | 時空洞窟已完成 | 时空洞窟已完成 | 時空の洞窟完了 | 시공의 동굴 완료 |
| helper ResourceStrings | Nothing to receive | 没有可以领取的 | | 受け取れるものはありません | 수령 가능한 것이 없습니다 |
| helper ResourceStrings | The expected diamond value of the current task is now below 20 | 当前任务的钻石数量期望值已低于20 | | 現在のタスクのダイヤの期待値が20未満になったため | 현재 작업의 다이아몬드 예상 값이 20 미만이므로 |

---

## 特殊来源提取规则

### Gacha 来源

Gacha 日志格式：`<prefix> <gacha_name> <count> <suffix>, ...`

**多语言前缀**：
| 语言 | 前缀 |
|------|------|
| EN | Gacha |
| TW/CN | 抽卡 |
| JA | ガチャ |
| KO | 가챠 |

**提取规则**：截取第二个空格之前的内容
- `Gacha 黒葬武具ガチャ 5 times` → `Gacha 黒葬武具ガチャ`
- `抽卡 黒葬武具ガチャ 5 次` → `抽卡 黒葬武具ガチャ`

### Open 来源

Open 日志格式：`<prefix> <item_name> x <count>`

**多语言前缀**：
| 语言 | 前缀 |
|------|------|
| EN | Open |
| TW | 開啟 |
| CN | 开启 |
| JA | 開く |
| KO | 열기 |

**提取规则**：截取 ` x` 之前的内容
- `Open Gold Sealed Chest x 5` → `Open Gold Sealed Chest`
- `開啟 上級封印寶箱 x 5` → `開啟 上級封印寶箱`

### Auto Refresh Diamonds 来源

自动刷新祈愿之泉任务的 helper 日志来自 `ResourceStrings(.XX).resx`：

| 语言 | 模板 |
|------|------|
| EN | `Current expected value: {0}. Today, {1}/{2} auto-refreshes completed. Refreshing now.` |
| ZH | `当前期望值：{0}，今日已自动刷新 {1}/{2} 次，现在刷新。` |
| JA | `現在の期待値：{0}。本日自動リフレッシュ {1}/{2} 回完了。リフレッシュ中。` |
| KO | `현재 기대값: {0}. 오늘 {1}/{2}회 자동 새로고침 완료. 지금 새로고침합니다.` |

**提取规则**：

- 每条日志视为一次自动刷新，固定写入钻石变动 `amount = -20`
- 来源写入 `Fountain of Prayers`
- 来源 ID 使用游戏内置 `SourceIDFountainOfPrayers = 140`
- 统计 JSON 中会与祈愿之泉获取量聚合到同一个 `id:140`：

```json
{
  "sources": {
    "id:140": {
      "source_id": 140,
      "gain": 26150,
      "consume": 6120
    }
  }
}
```

---

## 代码实现

### 来源映射表 (i18n/sources.go)

```go
// Game built-in source IDs from TextResource
const (
    SourceIDFountainOfPrayers SourceID = 140   // Fountain of Prayers -> 祈願之泉
    SourceIDOpen              SourceID = 67    // Open -> 開啟
    SourceIDLoginBonus        SourceID = 719   // Login Bonus -> 簽到獎勵
    SourceIDTowerInfinity     SourceID = 138   // Tower of Infinity -> 無窮之塔
    SourceIDTempleIllusions   SourceID = 2766  // Temple of Illusions -> 勝利
    SourceIDPresentsBox       SourceID = 21308 // Presents Box -> 禮物箱
    SourceIDMonthlyBoost      SourceID = 21332 // Monthly Boost -> 每月強化組合包
    SourceIDTotalLogins       SourceID = 3331  // Total Logins This Month -> 本月累計簽到天數
    SourceIDWorldClears       SourceID = 23277 // World Player Clears -> 本世界首次有玩家
)

// Helper custom source IDs (starting from 100002)
const (
    SourceIDAutoBuyStore       SourceID = 100002 // Auto Buy Store Items
    SourceIDMissionsClaimed    SourceID = 100004 // Missions Claim All
    SourceIDGacha              SourceID = 100005 // Gacha (抽卡)
    RewardMissionCompositeFactor SourceID = 1000000
    SourceIDDailyMissionReward SourceID = 23214000060 // Daily 明细示例: 23214 * 1000000 + 60
    SourceIDWeeklyMissionReward SourceID = 23215000080 // Weekly 明细示例: 23215 * 1000000 + 80
    SourceIDGuildMissionReward SourceID = 111002000 // Guild 明细示例: 111 * 1000000 + 2000
)

// MissionGroupType IDs from TextResource
const (
    MissionGroupDailyID  SourceID = 23214 // Daily -> 每日
    MissionGroupWeeklyID SourceID = 23215 // Weekly -> 每週
    MissionGroupMainID   SourceID = 23213 // Main -> 主線
)
```

### ChangeRecord 结构

```go
type ChangeRecord struct {
    Character string `json:"character"`
    Timestamp string `json:"timestamp"`
    Amount    int    `json:"amount"`
    Source    string `json:"source,omitempty"`
    SourceID  int    `json:"source_id,omitempty"` // 来源ID
}
```

### 输出示例

详细记录 JSONL 保留数值特征：

```json
{
  "character": "test",
  "timestamp": "2026-04-21T04:55:59+08:00",
  "amount": 100,
  "source": "id:23214000060",
  "source_id": 23214000060
}
```

统计 JSON 的 `sources` 使用聚合 ID：

```json
{
  "sources": {
    "id:23214000000": {
      "source_id": 23214000000,
      "gain": 100,
      "consume": 0
    }
  }
}
```

---

## 多语言支持

来源检测已支持以下语言：
- 英文 (EN)
- 繁中 (TW)
- 简中 (CN)
- 日文 (JA)
- 韩文 (KO)

---
生成时间: 2026-04-24
