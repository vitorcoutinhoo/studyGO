package service

import "golang.org/x/crypto/bcrypt"

func BcryptHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func BcryptCheckPassword(passwordHash string, rawPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), rawPassword)

	return err == nil
}
