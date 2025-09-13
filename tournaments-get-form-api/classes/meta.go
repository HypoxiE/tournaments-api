package classes

type CreateMetricInput struct {
	ResultID uint   `json:"result_id"`
	Key      string `json:"key"`
	Value    *int64 `json:"value"`
}
type Metric struct {
	ID       uint   `gorm:"column:metric_id;primaryKey" json:"metric_id"`
	ResultID uint   `gorm:"column:result_id;uniqueIndex:metric_idx_result_key" json:"-"`
	Result   Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key      string `gorm:"column:metric_key;uniqueIndex:metric_idx_result_key" json:"key"`
	Value    *int64 `gorm:"column:value" json:"value"`
}

func (input CreateMetricInput) NewMetricFromInput() Metric {
	return Metric{
		ResultID: input.ResultID,
		Key:      input.Key,
		Value:    OrElsePtr(input.Value, 0),
	}
}

type CreateMetadataInput struct {
	ResultID uint    `json:"result_id"`
	Key      string  `json:"key"`
	Value    *string `json:"value"`
}
type Metadata struct {
	ID       uint    `gorm:"column:metadata_id;primaryKey" json:"metadata_id"`
	ResultID uint    `gorm:"column:result_id;uniqueIndex:metadata_idx_result_key" json:"-"`
	Result   Result  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key      string  `gorm:"column:metadata_key;uniqueIndex:metadata_idx_result_key" json:"key"`
	Value    *string `gorm:"column:value" json:"value"`
}

func (input CreateMetadataInput) NewMetadataFromInput() Metadata {
	return Metadata{
		ResultID: input.ResultID,
		Key:      input.Key,
		Value:    OrElsePtr(input.Value, ""),
	}
}
