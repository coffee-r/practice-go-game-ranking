## 概要

Go言語で作るゲームランキングWebAPIの習作です。

## 作業記録

* 2024/12/30 作業開始

## 要件

* ユーザーを登録できる
* ユーザー一覧を取得できる
* ゲームを登録できる
* ゲーム一覧を取得できる
* ランキングを作成できる
    * ランキング名
    * 対応するゲーム
* ランキング一覧を取得できる
* ランキングに対してユーザーと紐付けてユーザースコアを登録できる
* あるランキングに対して登録できるスコアは1ユーザーにつき1つまでで、同じユーザーがスコアを登録しようとした時にはスコアが高い方の値を優先して更新する・ないしは更新しない
* スコアの値が同じ場合は登録日時が古い方を優先してランク付けする
* ランキングのスコアは整数型とし、昇順でランク付けする
* ランキングを指定してスコア一覧をソート済みで取得できる、ランクがつく
* ランキングとユーザーを指定して現在のランクを取得できる
* Delete系の機能は一旦実装対象外
* 認可については一旦実装対象外
* チート対策については一旦実装対象外
* 小数点などをケアするスコアについては一旦実装対象外
* ページネーションは一旦保留 カーソルベースのページネーション?
* パフォーマンスは一旦保留 実用的にするなら100万件のランク付けを高速にしたい

## データベーステーブル設計

* users
    * id (primary)
    * name
    * created_at
    * updated_at
* games
    * id (primary)
    * name
    * created_at
    * updated_at
* rankings
    * id (primary)
    * game_id (foreign key)
    * name
    * created_at
    * updated_at
* user_rankings
    * ranking_id (primary)
    * user_id (primary、foreign key)
    * score ※integer
    * created_at
    * updated_at

## REST API設計

* GET /users
* POST /users
* GET /games
* POST /games
* GET /rankings
* POST /rankings
* GET /rankings/{ranking_id}/user_rankings
* GET /rankings/{ranking_id}/user_rankings/{user_id}
* PUT /rankings/{ranking_id}/user_rankings/{user_id}

## ライブラリ

* Echo (ルーティング、パラメタのやり取り、jsonレスポンスを楽にしたいので)
* Bun (SQL操作を楽にしたいので)

## struct/interface図

* オニオンアーキテクチャっぽい構造とする
goplantuml