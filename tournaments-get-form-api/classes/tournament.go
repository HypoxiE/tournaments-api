package classes

type Tournament struct {
	ID uint `gorm:"primaryKey"`
	Name  string
	Formula string
	MainFields string
	Results []Result
}