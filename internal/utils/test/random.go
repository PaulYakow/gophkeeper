// Package test вспомогательный модуль для тестирования.
package test

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt генерирует случайное значение типа int в диапазоне [min, max].
func RandomInt(min, max int) int {
	return min + int(rand.Int63n(int64(max-min+1)))
}

// RandomString генерирует случайную последовательность символов из алфавита (alphabet) длиной n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomLogin генерирует случайный логин.
func RandomLogin() string {
	return RandomString(RandomInt(6, 12))
}

// RandomPassword генерирует случайный пароль.
func RandomPassword() string {
	return RandomString(RandomInt(8, 24))
}
