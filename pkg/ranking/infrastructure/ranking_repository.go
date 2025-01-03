package infrastructure

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
	"time"

	"github.com/uptrace/bun"
)

// ランキング
type Ranking struct {
	ID        int       `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ランキングリポジトリ
type RankingRepository struct {
	db *bun.DB
}

// リポジトリを生成する
func NewRankingRepository(bun *bun.DB) *RankingRepository {
	return &RankingRepository{
		db: bun,
	}
}

// ランキングをIDをキーとして取得する
func (r *RankingRepository) FindByID(ctx context.Context, id int) (*domain.Ranking, error) {
	// ランキング
	ranking := new(Ranking)

	// クエリ実行
	err := r.db.NewSelect().Model(&ranking).Where("id = ?", id).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ランキング名
	rankingName, err := domain.NewRankingName(ranking.Name)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメインのランキングを返す
	return &domain.Ranking{
		ID:        ranking.ID,
		Name:      rankingName,
		CreatedAt: ranking.CreatedAt,
		UpdatedAt: ranking.UpdatedAt,
	}, nil
}

// ランキングを名前をキーとして取得する
func (r *RankingRepository) FindByName(ctx context.Context, name domain.RankingName) (*domain.Ranking, error) {
	// ランキング
	ranking := new(Ranking)

	// クエリ実行
	err := r.db.NewSelect().Model(&ranking).Where("name = ?", name.Value).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメインのランキングを返す
	return &domain.Ranking{
		ID:        ranking.ID,
		Name:      name,
		CreatedAt: ranking.CreatedAt,
		UpdatedAt: ranking.UpdatedAt,
	}, nil
}

// ランキング一覧を取得する
func (r *RankingRepository) FindAll(ctx context.Context) ([]domain.Ranking, error) {
	// ランキングスライス
	var rankings []Ranking

	// クエリ実行
	err := r.db.NewSelect().Model(&rankings).Scan(ctx)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメイン層のランキング構造体にマッピング
	var domainRankings []domain.Ranking
	for _, ranking := range rankings {
		// ランキング名
		rankingName, err := domain.NewRankingName(ranking.Name)

		// エラーハンドリング
		if err != nil {
			log.Printf("Error occurred: %v", err)
			return nil, err
		}

		// domainRankingsに詰める
		domainRankings = append(domainRankings, domain.Ranking{
			ID:        ranking.ID,
			Name:      rankingName,
			CreatedAt: ranking.CreatedAt,
			UpdatedAt: ranking.UpdatedAt,
		})
	}

	// ドメインのランキングスライスを返す
	return domainRankings, nil
}

// ランキングを登録する
func (r *RankingRepository) Create(ctx context.Context, name domain.RankingName) (*domain.Ranking, error) {
	// ランキング構造体を生成
	ranking := &Ranking{
		Name: name.Value,
	}

	// ランキング名を生成
	rankingName, err := domain.NewRankingName(name.Value)

	// エラーハンドリング
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ランキング登録クエリを実行
	_, err = r.db.NewInsert().Model(ranking).Exec(ctx)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// 挿入後に ID を基に再取得
	err = r.db.NewSelect().Model(ranking).Where("id = ?", ranking.ID).Scan(ctx)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil, err
	}

	// ドメインのランキングを返す
	return &domain.Ranking{
		ID:        ranking.ID,
		Name:      rankingName,
		CreatedAt: ranking.CreatedAt,
		UpdatedAt: ranking.UpdatedAt,
	}, nil
}
