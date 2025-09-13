package classes

import (
	"encoding/json"
	"tournaments-api/funcs"

	"github.com/Knetic/govaluate"
)

type Result struct {
	ID           uint       `gorm:"column:result_id;primaryKey" json:"result_id"`
	TournamentID uint       `gorm:"column:tournament_id" json:"tournament_id"`
	Tournament   Tournament `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Username     string     `gorm:"column:username" json:"username"`
	Avatar       *string    `gorm:"column:avatar_url" json:"avatar_url"`
	Version      string     `gorm:"column:version" json:"version"`
	Score        *int       `gorm:"column:score" json:"score"`
	Penalty      *int       `gorm:"column:penalty" json:"penalty"`
	Cost         *int       `gorm:"column:cost" json:"cost"`
	// 1 - не проверено; 2 - проверено, разрешено; -1 - проверено, заблокировано; -2 - автоматическая блокировка
	Status    *int   `gorm:"column:status" json:"status"`
	Timestamp uint64 `gorm:"column:timestamp" json:"timestamp"`

	//confident
	PublicSteamID *string `gorm:"-" json:"steam_id"`
	SteamID       *string `gorm:"column:steam_id" json:"-"`
	PublicMail    string  `gorm:"-" json:"mail"`
	Mail          string  `gorm:"column:mail" json:"-"`
	PublicIP      string  `gorm:"-" json:"ip"`
	IP            string  `gorm:"column:ip;type:inet" json:"-"`

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
		params[m.Key] = *m.Value
	}

	score, err := expr.Evaluate(params)
	if err != nil {
		return err
	}

	result.Score = funcs.GetPtr(int(score.(float64)) - *result.Penalty)

	return nil
}

func (result Result) ToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(result)
	return jsonBytes, err
}

func (result *Result) PublicToPrivate() {
	result.SteamID = result.PublicSteamID
	result.Mail = result.PublicMail
}
