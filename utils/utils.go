package utils

import "math/rand"

/**
*@Author lyer
*@Date 2/20/21 16:27
*@Describe
**/

func RandInt(min int, max int) int {
	return rand.Intn((max - min) + min)
}

func Rand5Str() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
