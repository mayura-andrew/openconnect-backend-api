package data

import (
	"database/sql"
	"errors"
)








var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Ideas IdeaModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Ideas: IdeaModel{DB: db},
	}
}
