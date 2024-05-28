package user

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`
	Created       time.Time `bun:",notnull"`
	Updated       bun.NullTime
	ResetToken    string `bun:",nullzero"`
	Id            string `bun:",pk"`
	Email         string
	Password      string
	Name          string
}

func CreateNewUser(node *snowflake.Node, passwordHasher PasswordHasher, email string, password string, name string) (User, error) {
	pass, err := passwordHasher.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:         node.Generate().String(),
		Email:      email,
		Name:       name,
		Created:    time.Now(),
		Password:   pass,
		ResetToken: "",
	}, nil
}

func UserFromData(id string, email string, password string, name string, resetToken string, created time.Time, updated bun.NullTime) User {
	return User{
		Id:         id,
		Email:      email,
		Password:   password,
		Name:       name,
		Created:    created,
		Updated:    updated,
		ResetToken: resetToken,
	}
}

func (u *User) VerifyPassword(hasher PasswordHasher, password string) bool {
	return hasher.VerifyPassword(u.Password, password)
}

type UserRepository interface {
	GetById(id string) (*User, error)
	FindByEmailAndPassword(email string, password string) (*User, error)
	Save()
}
