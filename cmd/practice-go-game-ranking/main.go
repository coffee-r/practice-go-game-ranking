package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server用のドライバ
)

// ユーザー
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// データベース
var db *sql.DB

func main() {
	// json返す練習
	// http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	// 	u := User{
	// 		ID:   1,
	// 		Name: "hoge",
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")

	// 	j, err := json.MarshalIndent(u, "", "    ")
	// 	if err != nil {
	// 		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	_, err = w.Write(j)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}
	// })

	// データベースへの接続
	var err error
	db, err = sql.Open("sqlserver", "sqlserver://sa:r00tP@ss3014@db:1433?database=master&encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// エンドポイントを設定
	http.HandleFunc("/users", usersHandler)

	// サーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// /usersエンドポイント ハンドラー
func usersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		createUser(w, r)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// ユーザー作成処理
func createUser(w http.ResponseWriter, r *http.Request) {
	// HTTPヘッダをセット
	w.Header().Set("Content-Type", "application/json")

	// リクエストボディからデータを取得
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	// データベースにユーザーを挿入しつつ、挿入時のIDを取得
	query := "INSERT INTO users (name) OUTPUT INSERTED.ID VALUES (@p1)"
	var userID int
	err := db.QueryRow(query, sql.Named("p1", u.Name)).Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// レスポンス
	w.WriteHeader(http.StatusCreated)
	u.ID = userID
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Printf("Failed to write response: %v", err)
	}
}
