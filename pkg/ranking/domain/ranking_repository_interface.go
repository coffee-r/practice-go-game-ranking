package domain

import "context"

// ランキングリポジトリ (インターフェース)
type RankingRepositoryInterface interface {
	// ランキングをIDをキーとして取得する
	FindByID(ctx context.Context, id int) (*Ranking, error)

	// ランキングを名前をキーとして取得する
	FindByName(ctx context.Context, name RankingName) (*Ranking, error)

	// ランキング一覧を取得する
	FindAll(ctx context.Context) ([]Ranking, error)

	// ランキングを登録する
	Create(ctx context.Context, name RankingName) (*Ranking, error)
}
