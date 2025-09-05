package classes

import (
	"encoding/json"

	"github.com/Knetic/govaluate"
)

type Result struct {
	ID           uint       `gorm:"primaryKey" json:"result_id"`
	TournamentID uint       `json:"tournament_id"`
	Tournament   Tournament `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Username     string     `json:"username"`
	Avatar       string     `json:"avatar_url"`
	Version      string     `json:"version"`
	Score        int        `json:"score"`
	Cost         int        `json:"cost"`
	Status       int        `json:"status"` // 0 - не проверено; 1 - проверено, разрешено; -1 - проверено, заблокировано; -2 - автоматическая блокировка
	Timestamp    uint64     `json:"timestamp"`

	//confident
	GetterSteamID string `json:"steam_id" gorm:"-"`
	SteamID       string `json:"-"`
	GetterMail    string `json:"mail" gorm:"-"`
	Mail          string `json:"-"`
	IP            string `gorm:"type:inet" json:"-"`

	Metrics  []Metric   `json:"metrics"`
	Metadata []Metadata `json:"metadata"`
}

func (result *Result) CalculateScore(tournament Tournament) error {
	expr, err := govaluate.NewEvaluableExpression(tournament.Formula)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})
	for _, m := range result.Metrics {
		params[m.Key] = m.Value
	}

	score, err := expr.Evaluate(params)
	if err != nil {
		return err
	}

	result.Score = int(score.(float64))

	return nil
}

func (result Result) ToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(result)
	return jsonBytes, err
}
