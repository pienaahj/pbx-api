package model

import (
	"github.com/google/uuid"
)

// User defines domain model and its json and db representations
type User struct {
	UID      uuid.UUID `db:"uid" json:"uid"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"`
	Username string    `db:"username" json:"username"`
}

// UserCreds represents the log in credetials for a user
type UserCreds struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}
