package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
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
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name"`
	GameID    int       `json:"game_id" bun:"game_id"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
}

// ユーザースコア構造体
type UserScore struct {
	RankingID int       `json:"ranking_id" bun:"ranking_id,pk"`
	UserID    int       `json:"user_id" bun:"user_id,pk"`
	UserName  string    `json:"user_name" bun:"user_name"`
	Score     int       `json:"score" bun:"score"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,default:CURRENT_TIMESTAMP"`
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
	e.GET("/rankings", getRankings)
	e.POST("/rankings", createRanking)
	e.GET("/rankings/:ranking_id/user_rankings", getUserScoreRanking)
	e.PUT("/rankings/:ranking_id/user_rankings/:user_id", upsertUserScore)

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

// ランキング一覧取得処理
func getRankings(c echo.Context) error {
	// コンテキストからdbを取得
	db := c.Get("db").(*bun.DB)

	// ランキングスライスを定義
	var rankings []Ranking

	// クエリ実行
	err := db.NewSelect().Model(&rankings).Scan(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusOK, rankings)
}

// ランキング作成処理
func createRanking(c echo.Context) error {
	// コンテキストからdbとvalidatorを取得
	db := c.Get("db").(*bun.DB)
	validator := c.Get("validator").(*validator.Validate)

	// リクエストを受ける構造体を定義
	// ※複数代入を避けるためこの構造体を用意してます
	type RankingCreateRequest struct {
		GameID int    `json:"game_id" validate:"required"`
		Name   string `json:"name" validate:"required,max=50"`
	}

	// リクエストを受ける構造体を生成
	rankingCreateRequest := new(RankingCreateRequest)

	// リクエストボディをマッピング
	if err := c.Bind(rankingCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := validator.Struct(rankingCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "バリデーションエラー",
			"message": err.Error(),
		})
	}

	// リクエストパラメタとして受け取ったゲームIDに対応するゲームがあるかをチェック
	game := new(Game)
	exist, err := db.NewSelect().Model(game).Where("id = ?", rankingCreateRequest.GameID).Exists(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ゲーム取得に失敗しました。"})
	}
	if !exist {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ゲームが存在しません。"})
	}

	// ランキング構造体を生成
	ranking := &Ranking{
		GameID: rankingCreateRequest.GameID,
		Name:   rankingCreateRequest.Name,
	}

	// ランキング登録クエリを実行
	_, err = db.NewInsert().Model(ranking).Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング登録に失敗しました。"})
	}

	// 挿入後に ID を基に再取得
	err = db.NewSelect().Model(ranking).Where("id = ?", ranking.ID).Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusCreated, ranking)
}

// ユーザースコアランキング取得処理
func getUserScoreRanking(c echo.Context) error {
	// コンテキストからdbを取得
	db := c.Get("db").(*bun.DB)

	// ランキング構造体を生成
	var ranking Ranking

	// ランキング存在チェックのクエリ実行
	exist, err := db.NewSelect().Model(ranking).Where("id = ?", c.Param("ranking_id")).Exists(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング取得に失敗しました。"})
	}
	if !exist {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキングが存在しません。"})
	}

	// ユーザースコアランキングスライス
	var userScoreRankings []UserScoreRanking

	// ユーザースコアテーブルをベースとして、ユーザー・ランキングテーブルとjoin
	err = db.NewSelect().
		Table("user_scores").
		Join("JOIN users ON users.id = user_scores.user_id").
		Join("JOIN rankings ON rankings.id = user_scores.ranking_id").
		Where("user_scores.ranking_id = ?", c.Param("ranking_id")).
		Order("user_scores.score DESC").
		Scan(context.Background(), &userScoreRankings)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザースコアランキング取得に失敗しました。"})
	}

	// JSON形式でクライアントに返す
	return c.JSON(http.StatusOK, userScoreRankings)
}

// ユーザースコア登録・更新
func upsertUserScore(c echo.Context) error {
	// コンテキストからdbとvalidatorを取得
	db := c.Get("db").(*bun.DB)
	validator := c.Get("validator").(*validator.Validate)

	// パラメータから値を取得し、型変換
	rankingID, err := strconv.Atoi(c.Param("ranking_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ranking_id の変換に失敗しました"})
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id の変換に失敗しました"})
	}

	// ランキング構造体を生成
	var ranking Ranking

	// ランキング存在チェックのクエリ実行
	exist, err := db.NewSelect().Model(ranking).Where("id = ?", rankingID).Exists(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキング取得に失敗しました。"})
	}
	if !exist {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ランキングが存在しません。"})
	}

	// ユーザー構造体を生成
	var user User

	// ユーザー存在チェックのクエリ実行
	exist, err = db.NewSelect().Model(user).Where("id = ?", userID).Exists(c.Request().Context())

	// エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザー取得に失敗しました。"})
	}
	if !exist {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザーが存在しません。"})
	}

	// リクエストを受ける構造体を定義
	// ※複数代入を避けるためこの構造体を用意してます
	type UserScorePutRequest struct {
		Score int `json:"score" validate:"required"`
	}

	// リクエストを受ける構造体を生成
	userScorePutRequest := new(UserScorePutRequest)

	// リクエストボディをマッピング
	if err := c.Bind(userScorePutRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です。"})
	}

	// リクエストパラメタのバリデーション
	if err := validator.Struct(userScorePutRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "バリデーションエラー",
			"message": err.Error(),
		})
	}

	// ユーザースコア構造体を生成
	userScore := new(UserScore)

	// ユーザースコアを取得
	err = db.NewSelect().Model(userScore).Where("ranking_id = ? and user_id = ?", rankingID, userID).Scan(c.Request().Context())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザースコア取得に失敗しました。"})
	}

	// スコアが未登録の場合はINSERT
	if errors.Is(err, sql.ErrNoRows) {
		userScore.RankingID = rankingID
		userScore.UserID = userID
		userScore.Score = userScorePutRequest.Score

		_, err = db.NewInsert().Model(userScore).Exec(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザースコア登録に失敗しました。"})
		}

		// JSON形式でクライアントに返す
		return c.JSON(http.StatusCreated, userScore)
	}

	// リクエストされたスコアの方が大きい場合
	if userScorePutRequest.Score > userScore.Score {
		userScore.RankingID = rankingID
		userScore.UserID = userID
		userScore.Score = userScorePutRequest.Score

		_, err = db.NewUpdate().Model(userScore).Exec(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ユーザースコア更新に失敗しました。"})
		}
	}

	// それ以外は何もしない
	return c.JSON(http.StatusOK, userScore)
}
