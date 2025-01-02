package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server用のドライバ
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
)

// ユーザー構造体
type User struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ゲーム構造体
type Game struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ランキング構造体
type Ranking struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GameID    int    `json:"game_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ユーザースコア構造体
type UserScore struct {
	RankingID int    `json:"ranking_id"`
	UserID    int    `json:"user_id"`
	Score     int    `json:"score"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ユーザースコアランキング構造体
type UserScoreRanking struct {
	RankingID int `json:"ranking_id"`
	UserID    int `json:"user_id"`
	UserName  int `json:"user_name"`
	Score     int `json:"score"`
	Rank      int `json:"rank"`
}

func main() {
	// データベース (Bun)
	sqldb, err := sql.Open("sqlserver", "sqlserver://sa:r00tP@ss3014@db:1433?database=master&encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer sqldb.Close()
	db := bun.NewDB(sqldb, mssqldialect.New())

	// バリデーター
	validator := validator.New()

	// Echo
	e := echo.New()

	// ミドルウェアで `db` と `validator` を Context に登録
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			c.Set("validator", validator)
			return next(c)
		}
	})

	// エンドポイントを設定
	e.GET("/users", getUsers)
	e.POST("/users", createUser)
	e.GET("/games", getGames)
	e.POST("/games", createGame)

	// サーバを起動
	e.Logger.Fatal(e.Start(":8080"))
}

// ユーザー一覧取得処理
func getUsers(c echo.Context) error {
	// コンテキストからdbを取得
	db := c.Get("db").(*bun.DB)

	// ユーザースライスを定義
	var users []User

	// クエリ実行
	err := db.NewSelect().Model(&users).Scan(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusOK, users)
}

// ユーザー作成処理
func createUser(c echo.Context) error {
	// コンテキストからdbとvalidatorを取得
	db := c.Get("db").(*bun.DB)
	validator := c.Get("validator").(*validator.Validate)

	// リクエストを受ける構造体を定義
	// ※複数代入を避けるためこの構造体を用意してます
	type UserCreateRequest struct {
		Name string `json:"name" validate:"required,max=30"`
	}

	// リクエストを受ける構造体を生成
	userCreateRequest := new(UserCreateRequest)

	// リクエストボディをマッピング
	if err := c.Bind(userCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := validator.Struct(userCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "バリデーションエラー",
			"message": err.Error(),
		})
	}

	// ユーザー構造体を生成
	user := &User{
		Name: userCreateRequest.Name,
	}

	// ユーザー登録クエリを実行
	_, err := db.NewInsert().Model(user).Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー登録に失敗しました。"})
	}

	// 挿入後に ID を基に再取得
	err = db.NewSelect().Model(user).Where("id = ?", user.ID).Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusCreated, user)
}

// ゲーム一覧取得処理
func getGames(c echo.Context) error {
	// コンテキストからdbを取得
	db := c.Get("db").(*bun.DB)

	// ゲームスライスを定義
	var games []Game

	// クエリ実行
	err := db.NewSelect().Model(&games).Scan(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ゲーム取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusOK, games)
}

// ゲーム作成処理
func createGame(c echo.Context) error {
	// コンテキストからdbとvalidatorを取得
	db := c.Get("db").(*bun.DB)
	validator := c.Get("validator").(*validator.Validate)

	// リクエストを受ける構造体を定義
	// ※複数代入を避けるためこの構造体を用意してます
	type GameCreateRequest struct {
		Name string `json:"name" validate:"required,max=30"`
	}

	// リクエストを受ける構造体を生成
	gameCreateRequest := new(GameCreateRequest)

	// リクエストボディをマッピング
	if err := c.Bind(gameCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := validator.Struct(gameCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "バリデーションエラー",
			"message": err.Error(),
		})
	}

	// ゲーム構造体を生成
	game := &Game{
		Name: gameCreateRequest.Name,
	}

	// ゲーム登録クエリを実行
	_, err := db.NewInsert().Model(game).Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ゲーム登録に失敗しました。"})
	}

	// 挿入後に ID を基に再取得
	err = db.NewSelect().Model(game).Where("id = ?", game.ID).Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ゲーム取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusCreated, game)
}
