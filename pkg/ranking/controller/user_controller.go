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

// ユーザーコントローラー
type UserController struct {
	userUseCase *usecase.UserUsecase
	validator   *validator.Validate
	db          *bun.DB
}

// ユーザーコントローラーを生成する
func NewUserController(u *usecase.UserUsecase, v *validator.Validate, d *bun.DB) *UserController {
	return &UserController{
		userUseCase: u,
		validator:   v,
		db:          d,
	}
}

// ユーザー一覧を取得する
func (u *UserController) GetUsers(c echo.Context) error {
	// ユーザー一覧を取得
	users, err := u.userUseCase.GetUsers(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		log.Printf("[UserController.GetUsers] Failed to fetch users: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー一覧の取得に失敗しました。"})
	}

	// ユーザー一覧を返却する
	return c.JSON(http.StatusOK, users)
}

// ユーザーを新規登録する
func (u *UserController) CreateUser(c echo.Context) error {
	// リクエストを受ける構造体を定義
	type CreateUserRequest struct {
		Name string `json:"name" validate:"required,max=30"`
	}

	// リクエストを受ける構造体を生成
	createUserRequest := new(CreateUserRequest)

	// リクエストボディをマッピング
	if err := c.Bind(createUserRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := u.validator.Struct(createUserRequest); err != nil {
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
	tx, err := u.db.BeginTx(c.Request().Context(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "トランザクション開始に失敗しました"})
	}

	// ユーザーを新規登録
	user, err := u.userUseCase.CreateUser(c.Request().Context(), createUserRequest.Name)

	// エラーハンドリング
	if err != nil {
		tx.Rollback()
		log.Printf("[UserController.CreateUser] Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー登録に失敗しました。"})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "コミットに失敗しました"})
	}

	// 登録したユーザーを返却する
	return c.JSON(http.StatusCreated, user)
}
