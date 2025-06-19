package utils

import "math/rand"

func RandomPicker(list []string) string {
	return list[rand.Intn(len(list))]
}
