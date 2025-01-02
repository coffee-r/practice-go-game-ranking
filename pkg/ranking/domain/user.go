package domain

import "time"

// ユーザー (エンティティ)
type User struct {
	ID        int
	Name      UserName
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ユーザーを生成する
func NewUser(id int, name UserName, createdAt time.Time, updatedAt time.Time) (*User, error) {
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
