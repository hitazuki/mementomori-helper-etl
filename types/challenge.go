package types

// ChallengeType 挑战类型
type ChallengeType string

const (
	ChallengeTypeQuest ChallengeType = "quest" // 主线
	ChallengeTypeTower ChallengeType = "tower" // 塔
)

// TowerType 塔类型
type TowerType string

const (
	TowerInfinity TowerType = "Infinity"
	TowerAzure    TowerType = "Azure"
	TowerCrimson  TowerType = "Crimson"
	TowerEmerald  TowerType = "Emerald"
	TowerAmber    TowerType = "Amber"
)

// ChallengeStatus 挑战状态
type ChallengeStatus string

const (
	ChallengeStatusSuccess ChallengeStatus = "success"
	ChallengeStatusFailed  ChallengeStatus = "failed"
)

// ChallengeRecord 挑战记录（解析用，不持久化）
type ChallengeRecord struct {
	Character   string          `json:"-"`
	Timestamp   string          `json:"-"`
	PreciseTime string          `json:"-"`
	Type        ChallengeType   `json:"-"`
	Level       string          `json:"level"`
	TowerType   TowerType       `json:"tower_type,omitempty"`
	Status      ChallengeStatus `json:"-"`
}

// ChallengeLevelStats 单关卡统计
type ChallengeLevelStats struct {
	Level    string `json:"level"`
	Attempts int    `json:"attempts"`            // 尝试次数
	Success  bool   `json:"success"`             // 是否成功过
	LastTime string `json:"last_time,omitempty"` // 最后挑战时间
}

// ChallengeStats 角色挑战统计
type ChallengeStats struct {
	Quest  map[string]*ChallengeLevelStats               `json:"quest"`  // level -> stats
	Towers map[TowerType]map[string]*ChallengeLevelStats `json:"towers"` // tower -> level -> stats
}
