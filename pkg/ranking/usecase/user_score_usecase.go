package usecase

import (
	"context"
	"practice-go-game-ranking/pkg/ranking/domain"
)

// ユーザースコアユースケース
type UserScoreUseCase struct {
	userScoreRepository   domain.UserScoreRepositoryInterface
	userScoreQueryService UserScoreQueryServiceInterface
}

// ユースケースを生成する
func NewUserScoreUseCase(r domain.UserScoreRepositoryInterface, q UserScoreQueryServiceInterface) *UserScoreUseCase {
	return &UserScoreUseCase{
		userScoreRepository:   r,
		userScoreQueryService: q,
	}
}

// ユーザースコアの一覧を取得する
func (userScoreUseCase *UserScoreUseCase) GetUserScores(ctx context.Context, rankingID int) {
	userScoreUseCase.userScoreQueryService.FetchOrderByScoreDesc(rankingID)
}

// ユーザーのハイスコアを更新する
func (userScoreUseCase *UserScoreUseCase) UpdateUserHighScore(ctx context.Context, rankingID int, userID int, score int) {
	// ランキングを取得

	// エラーハンドリング

	// 当該ランキングが存在しない場合は更新できない

	// ユーザーを取得

	// エラーハンドリング

	// 当該ユーザーが存在しない場合は更新できない

	// ユーザーのスコアを取得

	// エラーハンドリング

	// スコアがない場合、またはハイスコアを更新した場合は永続化する

	// エラーハンドリング

}
