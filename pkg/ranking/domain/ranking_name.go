package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ランキング名 (値オブジェクト)
type RankingName struct {
	Value string
}

// ランキング名を生成する
func NewRankingName(name string) (RankingName, error) {
	// ランキング名の前後の空白を取り除く
	trimmedName := strings.TrimSpace(name)

	// ブランク文字、空白文字のみは許容しない
	if trimmedName == "" {
		return RankingName{}, fmt.Errorf("ランキング名は空にできません。入力された名前: %q", name)
	}

	// 50文字を超えたランキング名を許容しない
	if utf8.RuneCountInString(trimmedName) > 50 {
		return RankingName{}, fmt.Errorf("ランキング名は50文字以内である必要があります。入力された名前: %q", name)
	}

	// ランキング名を返却する
	return RankingName{Value: trimmedName}, nil
}
