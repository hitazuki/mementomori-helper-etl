package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

func BenchmarkGetSourceID(b *testing.B) {
	types.InitI18n(i18n.NewManager())
	
	testCases := []string{
		"Login",                                           // exact match
		"Fountain of Prayers: 100 diamonds",              // partial match (early)
		"Get Daily 's 60 Reward",                          // reward mission pattern
		"Unknown source that doesn't match anything",     // no match (worst case)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			GetSourceID(tc)
		}
	}
}

func BenchmarkGetSourceID_ExactMatch(b *testing.B) {
	types.InitI18n(i18n.NewManager())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSourceID("Login")
	}
}

func BenchmarkGetSourceID_PartialMatch(b *testing.B) {
	types.InitI18n(i18n.NewManager())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSourceID("Fountain of Prayers: 100 diamonds")
	}
}

func BenchmarkGetSourceID_RewardMission(b *testing.B) {
	types.InitI18n(i18n.NewManager())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSourceID("Get Daily 's 60 Reward")
	}
}

func BenchmarkGetSourceID_NoMatch(b *testing.B) {
	types.InitI18n(i18n.NewManager())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSourceID("Something that doesn't match any pattern at all")
	}
}
