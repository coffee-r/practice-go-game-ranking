package usecase

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
)

// ユーザーハイスコアユースケース
type UserHighScoreUseCase struct {
	rankingRepository       domain.RankingRepositoryInterface
	userRepository          domain.UserRepositoryInterface
	userHighScoreRepository domain.UserHighScoreRepositoryInterface
}

// ユースケースを生成する
func NewUserHighScoreUseCase(rankingRepo domain.RankingRepositoryInterface, userRepo domain.UserRepositoryInterface, userHighScoreRepo domain.UserHighScoreRepositoryInterface) *UserHighScoreUseCase {
	return &UserHighScoreUseCase{
		rankingRepository:       rankingRepo,
		userRepository:          userRepo,
		userHighScoreRepository: userHighScoreRepo,
	}
}

// ユーザーのハイスコアを更新する
func (userHighScoreUseCase *UserHighScoreUseCase) UpdateUserHighScore(ctx context.Context, rankingID int, userID int, newScore int) error {
	// ランキングの存在チェック
	_, err := userHighScoreUseCase.rankingRepository.FindByID(ctx, rankingID)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserHighScoreUseCase.UpdateUserHighScore] Failed to fetch ranking: %v", err)
		return err
	}

	// 当該ランキングが存在しない場合は更新できない

	// ユーザーの存在チェック
	_, err = userHighScoreUseCase.userRepository.FindByID(ctx, userID)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserHighScoreUseCase.UpdateUserHighScore] Failed to fetch user: %v", err)
		return err
	}

	// 当該ユーザーが存在しない場合は更新できない

	// ユーザーのハイスコアを取得
	// userHighScore, err := userHighScoreUseCase.userHighScoreRepository.Find(ctx, rankingID, userID)

	// エラーハンドリング

	// スコアがない場合、またはハイスコアを更新した場合は永続化する

	// エラーハンドリング

	return nil
}
