package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassWord(password string) (string, error) {
	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed_pw), nil
}

func CheckPassWordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
