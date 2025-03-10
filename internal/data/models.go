package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Ideas       IdeaModel
	Permissions PermissionModel
	Users       UserModal
	UserProfile ProfileModel
	Tokens      TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Ideas:       IdeaModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users:       UserModal{DB: db},
		UserProfile: ProfileModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
