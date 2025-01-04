package domain

// ユーザーハイスコア (エンティティ)
type UserHighScore struct {
	RankingID int
	UserID    int
	Score     int
}

// ユーザーハイスコアを生成する
func NewUserHighScore(rankingID int, userID int, score int) *UserHighScore {
	return &UserHighScore{
		RankingID: rankingID,
		UserID:    userID,
		Score:     score,
	}
}
