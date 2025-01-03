package usecase

type UserScoreQueryServiceInterface interface {
	FetchOrderByScoreDesc(rankingID int) ([]UserScoreRankingDto, error)
}
