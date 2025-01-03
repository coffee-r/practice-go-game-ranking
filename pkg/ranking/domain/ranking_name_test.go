package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 文字数の境界値テスト
func TestRankingNameWordCount(t *testing.T) {
	// 日本語50文字はOK
	_, err := NewRankingName("あいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえお")
	assert.NoError(t, err, "Expected no error for 50 characters (Japanese)")

	// 日本語51文字はNG
	_, err = NewRankingName("あいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあいうえおあ")
	assert.Error(t, err, "Expected error for 51 characters (Japanese)")
}

// 空白文字の取り扱い
func TestRankingNameWhiteTrim(t *testing.T) {
	// 空白文字だけ
	_, err := NewRankingName(" ")
	assert.Error(t, err, "Expected error for only blank char")

	// 前後の空白はトリミングされる
	rankingName, err := NewRankingName(" a ")
	assert.NoError(t, err, "Expected no error for trimmed ranking name")
	assert.Equal(t, "a", rankingName.Value, "Expected trimmed ranking name to be 'a'")

	// 文字の間にある空白はトリミングされない
	rankingName, err = NewRankingName("a b")
	assert.NoError(t, err, "Expected no error for trimmed ranking name")
	assert.Equal(t, "a b", rankingName.Value, "Expected trimmed ranking name to be 'a b'")
}
