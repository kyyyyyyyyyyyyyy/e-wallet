package domain

import "github.com/google/uuid"

type User struct {
	ID        string `gorm:"type:varchar(250);primaryKey" json:"id"`
	PublicKey string `gorm:"type:varchar(250)" json:"public_key"`
	Password  string `gorm:"type:varchar(250)" json:"password"`
	Balance   int64  `gorm:"type:bigint" json:"balance"`
}

func NewUser(publicKey string, password string, balance int64) *User {
	u := new(User)
	u.ID = uuid.NewString()
	u.PublicKey = publicKey
	u.Password = password
	u.Balance = balance
	return u
}
