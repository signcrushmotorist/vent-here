package domain

import "time"

type User struct {
	ID              int
	PublicID        string
	Email           string
	PasswordHash    string
	Username        string
	PublicAlias     string
	AliasChanges    int
	LastAliasChange *time.Time
	CreatedAt       time.Time
}
