package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ActivityLog struct {
	UserID string
	Message string
	Success bool
}

func (s *server) MigrateDb() error {
	s.db.AutoMigrate(&ActivityLog{})
	return nil
}