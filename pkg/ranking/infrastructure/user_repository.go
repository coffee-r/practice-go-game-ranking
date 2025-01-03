package infrastructure

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
	"time"

	"github.com/uptrace/bun"
)

// ユーザー
type User struct {
	ID        int       `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ユーザーリポジトリ
type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(bun *bun.DB) *UserRepository {
	return &UserRepository{
		db: bun,
	}
}

// ユーザー一覧を取得する
func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	// ユーザースライス
	var users []User

	// クエリ実行
	err := r.db.NewSelect().Model(&users).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメイン層のユーザー構造体にマッピング
	var domainUsers []domain.User
	for _, u := range users {
		// ユーザー名
		userName, err := domain.NewUserName(u.Name)

		// エラーハンドリング
		if err != nil {
			log.Printf("Error occurred: %v", err)
			return nil, err
		}

		// domainUsersに詰める
		domainUsers = append(domainUsers, domain.User{
			ID:        u.ID,
			Name:      userName,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	// ドメインのユーザースライスを返す
	return domainUsers, nil
}

// ユーザーを登録する
func (r *UserRepository) Create(ctx context.Context, name domain.UserName) (*domain.User, error) {
	// ユーザー構造体を生成
	user := &User{
		Name: name.Value,
	}

	// ユーザー名を生成
	userName, err := domain.NewUserName(name.Value)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ユーザー登録クエリを実行
	_, err = r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// 挿入後に ID を基に再取得
	err = r.db.NewSelect().Model(user).Where("id = ?", user.ID).Scan(ctx)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメインのユーザーを返す
	return &domain.User{
		ID:        user.ID,
		Name:      userName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
