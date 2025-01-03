package main

import (
	"database/sql"
	"log"
	"os"
	"practice-go-game-ranking/pkg/ranking/controller"
	"practice-go-game-ranking/pkg/ranking/infrastructure"
	"practice-go-game-ranking/pkg/ranking/usecase"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server用のドライバ
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
)

func main() {
	// .envファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数からポート番号を取得
	port := os.Getenv("PORT")
	if port == "" {
		// デフォルトポート
		port = "8080"
	}

	// データベース接続情報を環境変数から取得
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbEncrypt := os.Getenv("DB_ENCRYPT")

	// データベース接続
	dsn := "sqlserver://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "?database=" + dbName + "&encrypt=" + dbEncrypt
	sqldb, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer sqldb.Close()
	db := bun.NewDB(sqldb, mssqldialect.New())

	// Echo
	e := echo.New()

	// バリデーター
	validator := validator.New()

	// 依存関係のセットアップ
	userRepository := infrastructure.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userController := controller.NewUserController(userUseCase, validator, db)
	rankingRepository := infrastructure.NewRankingRepository(db)
	rankingUseCase := usecase.NewRankingUseCase(rankingRepository)
	rankingController := controller.NewRankingController(rankingUseCase, validator, db)

	// エンドポイント定義とControllerのマッピング
	e.GET("/users", userController.GetUsers)
	e.POST("/users", userController.CreateUser)
	e.GET("/rankings", rankingController.GetRankings)
	e.POST("/rankings", rankingController.CreateRanking)
	e.GET("/rankings/:ranking_id/user_high_scores")
	e.PUT("/rankings/:ranking_id/user_high_scores/:user_id")

	// サーバを起動
	e.Logger.Fatal(e.Start(":" + port))
}
