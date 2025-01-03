package usecase

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
)

// ユーザーユースケース
type UserUsecase struct {
	userRepository domain.UserRepositoryInterface
}

// ユースケースを生成する
func NewUserUseCase(r domain.UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{
		userRepository: r,
	}
}

// ユーザー一覧を取得する
func (u *UserUsecase) GetUsers(ctx context.Context) ([]UserDto, error) {
	// ユーザー一覧をリポジトリから取得する
	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		log.Printf("[UserUsecase.GetUsers] Failed to fetch users: %v", err)
		return nil, err
	}

	// ユースケース層の構造体にマッピング (スライスの容量を事前に確保)
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
		log.Printf("[UserUsecase.CreateUser] invalid user_name: %v", err)
		return nil, err
	}

	// リポジトリを使ってユーザーを登録する
	user, err := u.userRepository.Create(ctx, userName)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserUsecase.CreateUser] Failed to create new user: %v", err)
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
