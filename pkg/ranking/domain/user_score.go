package domain

// ユーザースコア (エンティティ)
type UserScore struct {
	RankingID int
	UserID    int
	Score     int
}

// ユーザースコアを生成する
func NewUserScore(rankingID int, userID int, score int) (*UserScore, error) {
	return &UserScore{
		RankingID: rankingID,
		UserID:    userID,
		Score:     score,
	}, nil
}
