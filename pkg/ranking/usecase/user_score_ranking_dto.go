package usecase

// ユーザースコアランキングDTO
type UserScoreRankingDto struct {
	RankingID      int                `json:"ranking_id"`
	RankingName    string             `json:"ranking_name"`
	UserScoreRanks []UserScoreRankDto `json:"user_score_ranks"`
}
