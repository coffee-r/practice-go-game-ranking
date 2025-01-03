package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 文字数の境界値テスト
func TestUserNameWordCount(t *testing.T) {
	// アルファベット30文字はOK
	_, err := NewUserName("abcdeabcdeabcdeabcdeabcdeabcde")
	assert.NoError(t, err, "Expected no error for 30 characters (alphabet)")

	// 日本語30文字はOK
	_, err = NewUserName("あいうえおあいうえおあいうえおあいうえおあいうえおあいうえお")
	assert.NoError(t, err, "Expected no error for 30 characters (Japanese)")

	// アルファベット31文字はNG
	_, err = NewUserName("abcdeabcdeabcdeabcdeabcdeabcdea")
	assert.Error(t, err, "Expected error for 31 characters (alphabet)")

	// 日本語31文字はNG
	_, err = NewUserName("あいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあ")
	assert.Error(t, err, "Expected error for 31 characters (Japanese)")
}

// 空白文字の取り扱い
func TestUserNameWhiteTrim(t *testing.T) {
	// 空白文字だけ
	_, err := NewUserName(" ")
	assert.Error(t, err, "Expected error for only blank char")

	// 前後の空白はトリミングされる
	userName, err := NewUserName(" a ")
	assert.NoError(t, err, "Expected no error for trimmed username")
	assert.Equal(t, "a", userName.Value, "Expected trimmed username to be 'a'")

	// 文字の間にある空白はトリミングされない
	userName, err = NewUserName("a b")
	assert.NoError(t, err, "Expected no error for trimmed username")
	assert.Equal(t, "a b", userName.Value, "Expected trimmed username to be 'a b'")
}
