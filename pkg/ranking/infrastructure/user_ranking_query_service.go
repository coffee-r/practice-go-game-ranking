package infrastructure

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/usecase"

	"github.com/uptrace/bun"
)

// ユーザーランク
type UserRank struct {
	UserID   int    `bun:"user_id"`
	UserName string `bun:"user_name"`
	Rank     int    `bun:"rank"`
	Score    int    `bun:"score"`
}

// ユーザーランキングクエリサービス
type UserRankingQueryService struct {
	db *bun.DB
}

// リポジトリを生成する
func NewUserRankingQueryService(bun *bun.DB) *UserRankingQueryService {
	return &UserRankingQueryService{
		db: bun,
	}
}

// ユーザーランキングを取得する
func (userRankingQueryService *UserRankingQueryService) FetchUserRanking(ctx context.Context, query usecase.UserRankingQuery) (*usecase.UserRankingDto, error) {
	// ランキング
	ranking := new(Ranking)

	// ランキング取得クエリ実行
	err := userRankingQueryService.db.NewSelect().Model(&ranking).Where("id = ?", query.RankingID).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ユーザーランクスライス
	var userRanks []UserRank

	// ユーザーハイスコアランキング取得クエリ実行
	err = userRankingQueryService.db.NewSelect().
		Table("user_high_scores").
		Join("JOIN users ON user_high_scores.user_id = users.id").
		Column("users.id AS user_id", "users.name AS user_name", "user_high_scores.score", "ROW_NUMBER() OVER (ORDER BY user_high_scores.ここはqueryから撮りたい) AS rank").
		Order("user_high_scores.rank DESC").
		Where("ranking_id = ?", query.RankingID).
		Scan(ctx, userRanks)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ユースケース層のユーザーランク構造体にマッピング
	var usecaseUserRanks []usecase.UserRankDto
	for _, userRank := range userRanks {
		usecaseUserRanks = append(usecaseUserRanks, usecase.UserRankDto{
			UserID:   userRank.UserID,
			UserName: userRank.UserName,
			Rank:     userRank.Rank,
			Score:    userRank.Score,
		})
	}

	// ユースケース層のユーザーランキング構造体にマッピング
	usecaseUserRanking := &usecase.UserRankingDto{
		RankingID:   ranking.ID,
		RankingName: ranking.Name,
		UserRanks:   usecaseUserRanks,
	}

	// ユーザーランキングを返却する
	return usecaseUserRanking, nil
}
