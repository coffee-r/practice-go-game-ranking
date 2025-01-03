package domain

import "time"

// ランキング (エンティティ)
type Ranking struct {
	ID        int
	Name      RankingName
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ランキングを生成する
func NewRanking(id int, name RankingName, createdAt time.Time, updatedAt time.Time) (*Ranking, error) {
	return &Ranking{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
