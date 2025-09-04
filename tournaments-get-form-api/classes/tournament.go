package classes

type Tournament struct {
	ID uint `gorm:"primaryKey"`
	Name  string
	Formula string
	Results []Result
}