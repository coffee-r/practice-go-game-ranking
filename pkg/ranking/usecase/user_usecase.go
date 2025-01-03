package usecase

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
)

// ユーザーユースケース
type UserUseCase struct {
	userRepository domain.UserRepositoryInterface
}

// ユースケースを生成する
func NewUserUseCase(r domain.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		userRepository: r,
	}
}

// ユーザー一覧を取得する
func (userUseCase *UserUseCase) GetUsers(ctx context.Context) ([]UserDto, error) {
	// ユーザー一覧をリポジトリから取得する
	users, err := userUseCase.userRepository.FindAll(ctx)
	if err != nil {
		log.Printf("[UserUseCase.GetUsers] Failed to fetch users: %v", err)
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
func (userUseCase *UserUseCase) CreateUser(ctx context.Context, name string) (*UserDto, error) {
	// ユーザー名
	userName, err := domain.NewUserName(name)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserUseCase.CreateUser] invalid user_name: %v", err)
		return nil, err
	}

	// リポジトリを使ってユーザーを登録する
	user, err := userUseCase.userRepository.Create(ctx, userName)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserUseCase.CreateUser] Failed to create new user: %v", err)
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
