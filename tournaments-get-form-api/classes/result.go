package classes

import (
	"encoding/json"
)

type Metric struct {
	ID uint `gorm:"primaryKey"`
	ResultID uint
	Result Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Key string `json:"key"`
	Value int64 `json:"value"`
}

type Result struct {
	ID uint `gorm:"primaryKey"`
	TournamentID uint `json:"tournament_id"`
	Tournament Tournament `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Username  string `json:"username"`
	Email string `json:"email"`
	SteamID string `json:"steam_id"`
	Score int
	Metrics []Metric `json:"metrics"`
}

// формат входных данных {"tournament_id": 1, "username": "HypoxiE", "email": "hypoxie@example.com", "metrics":[{"key":"humans", "value":1}, {"key":"animals", "value":38}]}
func From_Json(jsonData []byte) (Result, error) {
	var result Result
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (user Result) To_Json() ([]byte, error) {
	jsonBytes, err := json.Marshal(user)
	return jsonBytes, err
}