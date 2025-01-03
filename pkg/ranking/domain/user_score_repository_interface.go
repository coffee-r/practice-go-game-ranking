package domain

import "context"

// ユーザースコアリポジトリ (インターフェース)
type UserScoreRepositoryInterface interface {
	// ユーザースコアを保存する
	Store(ctx context.Context, userScore UserScore) error
}
