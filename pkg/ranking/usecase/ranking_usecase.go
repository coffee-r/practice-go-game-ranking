package usecase

import (
	"context"
	"log"
	"practice-go-game-ranking/pkg/ranking/domain"
)

// ランキングユースケース
type RankingUsecase struct {
	rankingRepository domain.RankingRepositoryInterface
}

// ユースケースを生成する
func NewRankingUseCase(r domain.RankingRepositoryInterface) *RankingUsecase {
	return &RankingUsecase{
		rankingRepository: r,
	}
}

// ランキング一覧を取得する
func (u *RankingUsecase) GetRankings(ctx context.Context) ([]RankingDto, error) {
	// ランキング一覧をリポジトリから取得する
	rankings, err := u.rankingRepository.FindAll(ctx)
	if err != nil {
		log.Printf("[RankingUsecase.GetRankings] Failed to fetch rankings: %v", err)
		return nil, err
	}

	// ユースケース層の構造体にマッピング (スライスの容量を事前に確保)
	rankingDtos := make([]RankingDto, 0, len(rankings))
	for _, r := range rankings {
		// ユーザーDTOにマッピング
		rankingDtos = append(rankingDtos, RankingDto{
			ID:        r.ID,
			Name:      r.Name.Value,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	// ユースケースのランキングを返す
	return rankingDtos, nil
}

// ランキングを新規登録する
func (u *RankingUsecase) CreateRanking(ctx context.Context, name string) (*RankingDto, error) {
	// ランキング名
	rankingName, err := domain.NewRankingName(name)

	// エラーハンドリング
	if err != nil {
		log.Printf("[RankingUsecase.CreateRanking] invalid ranking_name: %v", err)
		return nil, err
	}

	// ランキング名が既に登録されているか確認
	ranking, err := u.rankingRepository.FindByName(ctx, rankingName)

	// エラーハンドリング
	if err != nil {
		log.Printf("[RankingUsecase.CreateRanking] Failed to fetch ranking: %v", err)
		return nil, err
	}
	if ranking != nil {
		log.Printf("[RankingUsecase.CreateRanking] ranking name %v already used", err)
		return nil, err
	}

	// リポジトリを使ってランキングを登録する
	user, err := u.rankingRepository.Create(ctx, rankingName)

	// エラーハンドリング
	if err != nil {
		log.Printf("[RankingUsecase.CreateRanking] Failed to create new ranking: %v", err)
		return nil, err
	}

	// ユースケースのユーザーを返す
	return &RankingDto{
		ID:        user.ID,
		Name:      user.Name.Value,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
