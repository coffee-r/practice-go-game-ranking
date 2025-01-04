package usecase

// ユーザーランク
type UserRankDto struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}
