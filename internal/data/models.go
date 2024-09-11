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
	Ideas IdeaModel
	Users UserModal
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Ideas: IdeaModel{DB: db},
		Users: UserModal{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
