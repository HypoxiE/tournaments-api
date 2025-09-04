package classes

import "github.com/lib/pq"

type Tournament struct {
	ID uint `gorm:"primaryKey"`
	Name  string
	Variables pq.StringArray `gorm:"type:text[]"`
	Formula string
	Results []Result
}