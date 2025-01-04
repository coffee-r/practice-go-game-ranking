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

// ユーザーハイスコアコントローラー
type UserHighScoreController struct {
	userHighScoreUseCase *usecase.UserHighScoreUseCase
	validator            *validator.Validate
	db                   *bun.DB
}

// コントローラーを生成する
func NewUserHighScoreController(u *usecase.UserHighScoreUseCase, v *validator.Validate, d *bun.DB) *UserHighScoreController {
	return &UserHighScoreController{
		userHighScoreUseCase: u,
		validator:            v,
		db:                   d,
	}
}

// ハイスコアを登録する
func (userHighScoreController *UserHighScoreController) StoreHighScore(c echo.Context) error {
	// リクエストを受ける構造体を定義
	type CreateUserHighScoreRequest struct {
		RankingID int `json:"ranking_id" param:"ranking_id" validate:"required`
		UserID    int `json:"user_id" param:"user_id" validate:"required`
		Score     int `json:"score" validate:"required`
	}

	// リクエストを受ける構造体を生成
	createUserHighScoreRequest := new(CreateUserHighScoreRequest)

	// リクエストボディをマッピング
	if err := c.Bind(createUserHighScoreRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := userHighScoreController.validator.Struct(createUserHighScoreRequest); err != nil {
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
	tx, err := userHighScoreController.db.BeginTx(c.Request().Context(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "トランザクション開始に失敗しました"})
	}

	// ハイスコアを登録
	err = userHighScoreController.userHighScoreUseCase.UpdateUserHighScore(c.Request().Context(), createUserHighScoreRequest.RankingID, createUserHighScoreRequest.UserID, createUserHighScoreRequest.Score)

	// エラーハンドリング
	if err != nil {
		tx.Rollback()
		log.Printf("[UserController.CreateUser] Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ハイスコア更新に失敗しました。"})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "コミットに失敗しました"})
	}

	// 更新結果を返却する
	return c.JSON(http.StatusOK, nil)
}
