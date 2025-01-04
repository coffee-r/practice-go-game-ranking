package domain

import "time"

// ユーザー (エンティティ)
type User struct {
	ID        int
	Name      UserName
	CreatedAt time.Time
	UpdatedAt time.Time
}
