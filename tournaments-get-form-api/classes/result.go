package classes

import (
	"encoding/json"
	"time"

	"github.com/Knetic/govaluate"
)

type CreateResultInput struct {
	TournamentID uint    `json:"tournament_id"`
	Username     string  `json:"username"`
	Avatar       *string `json:"avatar_url"`
	Version      string  `json:"version"`
	Cost         *int    `json:"cost"`

	PublicSteamID *string `json:"steam_id"`
	PublicMail    string  `json:"mail"`
	PublicIP      *string `json:"ip"`

	Metrics  []CreateMetricInput   `json:"metrics"`
	Metadata []CreateMetadataInput `json:"metadata"`
}
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
	PublicIP      *string `gorm:"-" json:"ip"`
	IP            *string `gorm:"column:ip;type:inet" json:"-"`

	Metrics  []Metric   `json:"metrics"`
	Metadata []Metadata `json:"metadata"`
}

func (input CreateResultInput) NewResultFromInput(ip string) Result {
	var result = Result{
		TournamentID: input.TournamentID,

		Username:  input.Username,
		Avatar:    OrElsePtr(input.Avatar, ""),
		Version:   input.Version,
		Score:     GetPtr(0),
		Penalty:   GetPtr(0),
		Cost:      OrElsePtr(input.Cost, 0),
		Status:    GetPtr(1),
		Timestamp: uint64(time.Now().Unix()),

		PublicSteamID: OrElsePtr(input.PublicSteamID, ""),
		SteamID:       OrElsePtr(input.PublicSteamID, ""),
		PublicMail:    input.PublicMail,
		Mail:          input.PublicMail,
		PublicIP:      GetPtr(ip),
		IP:            GetPtr(ip),
	}

	for _, v := range input.Metrics {
		result.Metrics = append(result.Metrics, v.NewMetricFromInput())
	}
	for _, v := range input.Metadata {
		result.Metadata = append(result.Metadata, v.NewMetadataFromInput())
	}

	return result
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

	result.Score = GetPtr(int(score.(float64)) - *result.Penalty)

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
