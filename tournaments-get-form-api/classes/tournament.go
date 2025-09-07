package classes

import (
	"encoding/json"

	"github.com/lib/pq"
)

type Tournament struct {
	ID             uint           `gorm:"primaryKey" json:"tournament_id"`
	Name           string         `json:"name"`
	StartTimestamp uint64         `json:"start_timestamp"`
	StopTimestamp  uint64         `json:"stop_timestamp"`
	Metadata       pq.StringArray `gorm:"type:text[]" json:"metadata"`
	Variables      pq.StringArray `gorm:"type:text[]" json:"variables"`
	Formula        string         `json:"formula"`
	Results        []Result       `json:"results"`
}

func (tournament Tournament) ToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(tournament)
	return jsonBytes, err
}
