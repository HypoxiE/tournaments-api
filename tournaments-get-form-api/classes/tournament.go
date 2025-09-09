package classes

import (
	"encoding/json"

	"github.com/lib/pq"
)

type Tournament struct {
	ID             uint           `gorm:"column:tournament_id;primaryKey" json:"tournament_id"`
	Name           string         `gorm:"column:name" json:"name"`
	StartTimestamp uint64         `gorm:"column:start_timestamp" json:"start_timestamp"`
	StopTimestamp  uint64         `gorm:"column:stop_timestamp" json:"stop_timestamp"`
	Metadata       pq.StringArray `gorm:"column:metadata;type:text[]" json:"metadata"`
	Variables      pq.StringArray `gorm:"column:variables;type:text[]" json:"variables"`
	Formula        string         `gorm:"column:formula" json:"formula"`

	Results []Result `json:"results"`
}

func (tournament Tournament) ToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(tournament)
	return jsonBytes, err
}
