package shortener

import (
	"fmt"
	"strings"

	"url-shortener/internal/model"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// populated once when the server runs
var decodeMap = func() map[rune]int64 {
	m := make(map[rune]int64)

	for i, c := range charset {
		m[c] = int64(i)
	}

	return m
}()

func EncodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	base := int64(len(charset))
	var result strings.Builder

	for num > 0 {
		remainder := num % base
		result.WriteByte(charset[remainder])
		num = num / base
	}

	// reverse string
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func DecodeBase62(s string) (int64, error) {
	if s == "" {
		return 0, model.ErrEmptyBase62String
	}

	var result int64
	base := int64(len(charset))

	for _, c := range s {
		value, ok := decodeMap[c]
		if !ok {
			return 0, fmt.Errorf("invalid base62 character: %c", c)
		}

		result = result*base + value
	}

	return result, nil
}
