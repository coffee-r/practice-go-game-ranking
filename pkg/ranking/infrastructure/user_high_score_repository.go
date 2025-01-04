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

// ユーザーハイスコア
type UserHighScore struct {
	RankingID int       `bun:"ranking_id,pk"`
	UserID    int       `bun:"user_id,pk"`
	HighScore int       `bun:"high_score"`
	Timestamp time.Time `bun:"timestamp,nullzero,default:CURRENT_TIMESTAMP"`
}

// ユーザーハイスコアリポジトリ
type UserHighScoreRepository struct {
	db *bun.DB
}

// リポジトリを生成する
func NewUserHighScoreRepository(bun *bun.DB) *UserHighScoreRepository {
	return &UserHighScoreRepository{
		db: bun,
	}
}

// ユーザーハイスコアを取得する
// ランキングをIDをキーとして取得する
func (r *UserHighScoreRepository) Find(ctx context.Context, rankingID int, userID int) (*domain.UserHighScore, error) {
	// ユーザーハイスコア
	userHighScore := new(UserHighScore)

	// クエリ実行
	err := r.db.NewSelect().Model(&userHighScore).Where("ranking_id = ? and user_id = ?", rankingID, userID).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメインのユーザーハイスコアを返す
	return &domain.UserHighScore{
		RankingID: userHighScore.RankingID,
		UserID:    userHighScore.UserID,
		Score:     userHighScore.HighScore,
	}, nil
}

// ユーザーハイスコアを保存する
func (r *UserHighScoreRepository) Store(ctx context.Context, rankingID int, userID int, score int) error {
	// ユーザーハイスコア
	userHighScore := new(UserHighScore)

	// 既にスコア登録されているかを確認するクエリを投げる
	err := r.db.NewSelect().
		Model(&userHighScore).
		Where("ranking_id = ? and user_id = ?", rankingID, userID).
		Scan(ctx)

	// エラーハンドリング
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// データ取得時の予期せぬエラー
		log.Printf("Failed to query user score: %v", err)
		return err
	}

	// データがなければINSERT
	if errors.Is(err, sql.ErrNoRows) {
		userHighScore = &UserHighScore{
			RankingID: rankingID,
			UserID:    userID,
			HighScore: score,
		}

		// スコア登録クエリを実行
		_, err = r.db.NewInsert().Model(userHighScore).Exec(ctx)
		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}
	} else {
		// データがあればUPDATE
		_, err := r.db.NewUpdate().
			Table("user_high_scores").
			Set("score = ?, timestamp = getdate()", score).
			Where("ranking_id = ? AND user_id = ?", rankingID, userID).
			Exec(ctx)

		if err != nil {
			log.Printf("Error occurred: %v", err)
			return err
		}
	}

	return nil
}
