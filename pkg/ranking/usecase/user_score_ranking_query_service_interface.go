package usecase

type UserScoreRankingQueryServiceInterface interface {
	FetchOrderByScoreDesc(rankingID int) ([]UserScoreRankingDto, error)
}
