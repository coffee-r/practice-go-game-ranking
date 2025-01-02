package usecase

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
	"time"
)

// ユーザーユースケース
type UserUsecase struct {
	UserRepository domain.UserRepositoryInterface
}

// ユーザーDTO
type UserDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ユーザー一覧を取得する
func (u *UserUsecase) GetUsers(ctx context.Context) ([]UserDto, error) {
	// ユーザー一覧をリポジトリから取得する
	users, err := u.UserRepository.FindAll(ctx)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ユースケース層の構造体にマッピング
	// スライスの容量を事前に確保
	userDtos := make([]UserDto, 0, len(users))
	for _, u := range users {
		// ユーザーDTOにマッピング
		userDtos = append(userDtos, UserDto{
			ID:        u.ID,
			Name:      u.Name.Value,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	// ユースケースのユーザーを返す
	return userDtos, nil
}

// ユーザーを新規登録する
func (u *UserUsecase) CreateUser(ctx context.Context, name string) (*UserDto, error) {
	// ユーザー名
	userName, err := domain.NewUserName(name)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// リポジトリを使ってユーザーを登録する
	user, err := u.UserRepository.Create(ctx, userName)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ユースケースのユーザーを返す
	return &UserDto{
		ID:        user.ID,
		Name:      user.Name.Value,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
