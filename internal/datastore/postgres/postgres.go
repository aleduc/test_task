package postgres

import (
	"gorm.io/gorm"
)

const (
	equalStmt = "= ?"
)

func StringEqFilterScope(field string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(value) > 0 {
			db = db.Where(field+equalStmt, value)
		}
		return db
	}
}
