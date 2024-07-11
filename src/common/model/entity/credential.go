package entity

import "time"

type Credential struct {
	Email    string `json:"email" gorm:"primary_key;column:email"`
	Password string `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" grom:"column:updated_at;autoUpdateTime"`
}

func (c *Credential) TableName() string {
	return "credentials"
}
