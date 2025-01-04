package usecase

// ユーザーランキングのクエリ条件
type UserRankingQuery struct {
	RankingID int
	OrderBy   string
}
