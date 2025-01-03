package domain

import "context"

// ユーザーリポジトリ (インターフェース)
type UserRepositoryInterface interface {
	// ユーザーを取得する
	FindByID(ctx context.Context, id int) (*User, error)

	// ユーザー一覧を取得する
	FindAll(ctx context.Context) ([]User, error)

	// ユーザーを登録する
	Create(ctx context.Context, name UserName) (*User, error)
}
