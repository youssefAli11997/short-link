package shortener

import "strings"

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

// TODO(ya): consider adding DecodeBase62 and testing the round-trip
// Note: DecodeBase62 is not needed so far, so this TODO is not a priority
