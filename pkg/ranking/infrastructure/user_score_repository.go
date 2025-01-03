package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
	"time"

	"github.com/uptrace/bun"
)

// ユーザースコア
type UserScore struct {
	RankingID int       `bun:"ranking_id,pk"`
	UserID    int       `bun:"user_id,pk"`
	Score     int       `bun:"score"`
	CreatedAt time.Time `bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ユーザースコアリポジトリ
type UserScoreRepository struct {
	db *bun.DB
}

// ユーザースコアを保存する
func (r *UserScoreRepository) Store(ctx context.Context, domainUserScore domain.UserScore) error {
	// ユーザースコア
	userScore := new(UserScore)

	// 既にスコア登録されているかを確認するクエリを投げる
	err := r.db.NewSelect().
		Model(&userScore).
		Where("ranking_id = ? and user_id = ?", domainUserScore.RankingID, domainUserScore.UserID).
		Scan(ctx)

	// エラーハンドリング
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// データ取得時の予期せぬエラー
		log.Printf("Failed to query user score: %v", err)
		return err
	}

	// データがなければINSERT
	if errors.Is(err, sql.ErrNoRows) {
		userScore = &UserScore{
			RankingID: domainUserScore.RankingID,
			UserID:    domainUserScore.Score,
			Score:     domainUserScore.Score,
		}

		// スコア登録クエリを実行
		_, err = r.db.NewInsert().Model(userScore).Exec(ctx)
		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}
	} else {
		// データがあればUPDATE
		_, err := r.db.NewUpdate().
			Model(userScore).
			Set("score = ?, updated_at = ?", userScore.Score, time.Now()).
			Where("ranking_id = ? AND user_id = ?", domainUserScore.RankingID, domainUserScore.UserID).
			Exec(ctx)

		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}
	}

	return nil
}
