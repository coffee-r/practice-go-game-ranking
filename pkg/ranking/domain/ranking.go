package domain

import "time"

// ランキング (エンティティ)
type Ranking struct {
	ID        int
	Name      RankingName
	CreatedAt time.Time
	UpdatedAt time.Time
}
