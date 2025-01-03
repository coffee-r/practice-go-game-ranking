package controller

import (
	"fmt"
	"log"
	"net/http"
	"practice-go-game-ranking/pkg/ranking/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

// ランキングコントローラー
type RankingController struct {
	rankingUseCase *usecase.RankingUseCase
	validator      *validator.Validate
	db             *bun.DB
}

// コントローラーを生成する
func NewRankingController(u *usecase.RankingUseCase, v *validator.Validate, d *bun.DB) *RankingController {
	return &RankingController{
		rankingUseCase: u,
		validator:      v,
		db:             d,
	}
}

// ランキング一覧を取得する
func (rankingController *RankingController) GetRankings(c echo.Context) error {
	// ランキング一覧を取得
	rankings, err := rankingController.rankingUseCase.GetRankings(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserController.GetUsers] Failed to fetch users: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング一覧の取得に失敗しました。"})
	}

	// ランキング一覧を返却する
	return c.JSON(http.StatusOK, rankings)
}

// ランキングを新規登録する
func (rankingController *RankingController) CreateRanking(c echo.Context) error {
	// リクエストを受ける構造体を定義
	type CreateRankingRequest struct {
		Name string `json:"name" validate:"required,max=50"`
	}

	// リクエストを受ける構造体を生成
	createRankingRequest := new(CreateRankingRequest)

	// リクエストボディをマッピング
	if err := c.Bind(createRankingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := rankingController.validator.Struct(createRankingRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, vErr := range validationErrors {
			messages = append(messages, fmt.Sprintf("フィールド '%s' の値が不正です: %s", vErr.Field(), vErr.Tag()))
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "バリデーションエラー",
			"message": messages,
		})
	}

	// トランザクションを管理する
	tx, err := rankingController.db.BeginTx(c.Request().Context(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "トランザクション開始に失敗しました"})
	}

	// ランキングを新規登録
	ranking, err := rankingController.rankingUseCase.CreateRanking(c.Request().Context(), createRankingRequest.Name)

	// エラーハンドリング
	if err != nil {
		tx.Rollback()
		log.Printf("[RankingController.CreateRanking] Failed to create ranking: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング登録に失敗しました。"})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "コミットに失敗しました"})
	}

	// 登録したランキングを返却する
	return c.JSON(http.StatusCreated, ranking)
}
