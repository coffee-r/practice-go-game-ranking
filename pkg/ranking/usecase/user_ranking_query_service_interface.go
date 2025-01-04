package usecase

import "context"

// ユーザーランキングのクエリサービス
type UserRankingQueryServiceInterface interface {
	// ユーザーランキングを取得する
	FetchUserRanking(ctx context.Context, query UserRankingQuery) (*UserRankingDto, error)
}
