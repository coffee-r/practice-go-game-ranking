package controller

import (
	"fmt"
	"log"
	"net/http"
	"practice-go-game-ranking/pkg/ranking/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ユーザーランキングコントローラー
type UserRankingController struct {
	userRankingQueryService usecase.UserRankingQueryServiceInterface
	validator               *validator.Validate
}

// コントローラーを生成する
func NewUserRankingController(q usecase.UserRankingQueryServiceInterface, v *validator.Validate) *UserRankingController {
	return &UserRankingController{
		userRankingQueryService: q,
		validator:               v,
	}
}

// ユーザーランキングを取得する
func (userRankingController *UserRankingController) GetUserRanking(c echo.Context) error {
	// リクエストを受ける構造体を定義
	type GetUserRankingRequest struct {
		RankingID int    `json:"ranking_id" param:"ranking_id" validate:"required`
		OrderBy   string `json:"order_by" query:"order_by" validate:"required`
		Limit     int    `json:"limit" query:"limit" validate:"required`
	}

	// リクエストを受ける構造体を生成
	getUserRankingRequest := new(GetUserRankingRequest)

	// リクエストボディをマッピング
	if err := c.Bind(getUserRankingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := userRankingController.validator.Struct(getUserRankingRequest); err != nil {
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

	// クエリ構造体を生成
	query := usecase.UserRankingQuery{
		RankingID: getUserRankingRequest.RankingID,
		OrderBy:   getUserRankingRequest.OrderBy,
	}

	// ユーザーランキングを取得
	userRanking, err := userRankingController.userRankingQueryService.FetchUserRanking(c.Request().Context(), query)

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserRankingController.CreateUser] Failed to fetch user ranking: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザーランキングの取得に失敗しました"})
	}

	// ランキングを返却する
	return c.JSON(http.StatusOK, userRanking)
}
