## 概要

Go言語で作るゲームランキングWebAPIの習作です。

## 作業記録

* 2024/12/30 作業開始

## 要件

* ユーザーを登録できる
* ゲームを登録できる
* ランキングを作成できる
    * ランキング名
    * 対応するゲーム
* ランキングに対してユーザーと紐付けてスコアを登録できる
* あるランキングに対して登録できるスコアは1ユーザーにつき1つまでで、同じユーザーがスコアを登録しようとした時にはスコアが高い方の値を優先して更新する・ないしは更新しない
* スコアの値が同じ場合は登録日時が古い方を優先してランク付けする
* ランキングのスコアは整数型とし、昇順でランク付けする
* ランキングを指定してスコア一覧をソート済みで取得できる、ランクがつく
* ランキングとユーザーを指定して現在のランクを取得できる
* ランキングは無効化ができる (ソフトデリート)、無効化したランキングに対してスコアを登録することはできない
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
    * enable (BOOLEAN)
    * created_at
    * updated_at
    * ※rankingsはupdateではname、enableしか変えれない
* user_rankings
    * ranking_id (primary)
    * user_id (primary、foreign key)
    * score ※integer
    * created_at
    * updated_at

## REST API設計

* POST /users
* PATCH /users/{user_id}
* POST /games
* PATCH /games/{game_id}
* POST /rankings
* PATCH /rankings/{ranking_id}
* GET /rankings/{ranking_id}/user_rankings
* GET /rankings/{ranking_id}/user_rankings/{user_id}
* PUT /rankings/{ranking_id}/user_rankings/{user_id}

## struct/interface図
goplantuml