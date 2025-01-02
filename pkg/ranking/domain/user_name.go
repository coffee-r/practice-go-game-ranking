package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ユーザー名 (値オブジェクト)
type UserName struct {
	Value string
}

// ユーザー名を生成する
func NewUserName(name string) (UserName, error) {
	// ユーザー名の前後の空白を取り除く
	trimmedName := strings.TrimSpace(name)

	// ブランク文字、空白文字のみは許容しない
	if trimmedName == "" {
		return UserName{}, fmt.Errorf("ユーザー名は空にできません。入力された名前: %q", name)
	}

	// 30文字を超えたユーザー名を許容しない
	if utf8.RuneCountInString(trimmedName) > 30 {
		return UserName{}, fmt.Errorf("ユーザー名は30文字以内である必要があります。入力された名前: %q", name)
	}

	// ユーザー名を返却する
	return UserName{Value: trimmedName}, nil
}
