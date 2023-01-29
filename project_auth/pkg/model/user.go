package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64
	UserName       string
	PasswordDigest string
	Phone          string
}

// SetPassword 密码加密
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
