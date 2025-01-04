package domain

import "context"

// ユーザーハイスコアリポジトリ (インターフェース)
type UserHighScoreRepositoryInterface interface {
	// ユーザーハイスコアを取得する
	Find(ctx context.Context, rankingID int, userID int) (*UserHighScore, error)

	// ユーザーハイスコアを保存する
	Store(ctx context.Context, rankingID int, userID int, score int) error
}
