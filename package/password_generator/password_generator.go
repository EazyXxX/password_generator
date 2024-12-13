package password_generator

import "math/rand"

func GenerateRandomPasssword() string {
	allCharacters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz@#$%^&*!,-+_"
	length := 12

	password := make([]byte, length)

	for i := 0; i < length; i++ {
		index := rand.Intn(len(allCharacters))
		password[i] = allCharacters[index]
	}

	return string(password)
}
