package user

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword string, password string) bool
}
type bcryptPasswordHasher struct{}

func (ph *bcryptPasswordHasher) HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func (ph *bcryptPasswordHasher) VerifyPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

var BcryptPasswordHasher *bcryptPasswordHasher = &bcryptPasswordHasher{}
