package classes

type Metric struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	ResultID uint   `json:"-"`
	Result   Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key      string `json:"key"`
	Value    int64  `json:"value"`
}

type Metadata struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	ResultID uint   `json:"-"`
	Result   Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}
