package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/url"
)

type ActivityLog struct {
	UserID string
	Message string
	Success bool
}

type Repo struct {
	db *gorm.DB
}

func newRepo(engine string, config string) (*Repo, error){
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&ActivityLog{})

	return &Repo{
		db: db,
	}, nil
}

func (repo *Repo) createLog(r *url.Values, success bool) error {
	userId := r.Get("user_id")
	text := r.Get("text")

	a := &ActivityLog{
		UserID: userId,
		Message: text,
		Success: success,
	}
	repo.db.Create(a)

	return nil
}