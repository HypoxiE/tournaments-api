package classes

import "github.com/lib/pq"

type Tournament struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Metadata  pq.StringArray `gorm:"type:text[]"`
	Variables pq.StringArray `gorm:"type:text[]"`
	Formula   string
	Results   []Result
}
