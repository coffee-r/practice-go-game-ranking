package usecase

// ユーザーランキング
type UserRankingDto struct {
	RankingID   int           `json:"ranking_id"`
	RankingName string        `json:"ranking_name"`
	UserRanks   []UserRankDto `json:"user_ranks"`
}
