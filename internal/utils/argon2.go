package utils

import "github.com/matthewhartstonge/argon2"

func Hash(password string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))

	return string(encoded), err
}

func Verify(password, encoded string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(encoded))

	return ok, err
}
