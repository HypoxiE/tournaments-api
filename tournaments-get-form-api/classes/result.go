package classes

import (
	"encoding/json"
	"errors"
	"net"
)

type Metric struct {
	ID uint `gorm:"primaryKey" json:"-"`
	ResultID uint `json:"-"`
	Result Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key string `json:"key"`
	Value int64 `json:"value"`
}

type Result struct {
	ID uint `gorm:"primaryKey" json:"-"`
	IP net.IP `gorm:"type:inet" json:"-"`
	TournamentID uint `json:"tournament_id"`
	Tournament Tournament `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Username  string `json:"username"`
	Avatar string `json:"avatar_url"`
	Mail string `json:"mail"`
	SteamID string `json:"steam_id"`
	Version string `json:"version"`

	Score int `json:"score"`
	Cost int `json:"cost"`
	Status int `json:"status"` // 0 - не проверено; 1 - проверено, разрешено; -1 - проверено, заблокировано; -2 - автоматическая блокировка
	Metrics []Metric `json:"metrics"`
}

// формат входных данных {"tournament_id": 1, "username": "HypoxiE", "mail": "hypoxie@example.com", "metrics":[{"key":"humans", "value":1}, {"key":"animals", "value":38}]}
func From_Json(jsonData []byte, ip string) (Result, error) {
	var result Result
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, err
	}

	if (result.TournamentID == 0) {
		err := errors.New("error: The tournament_id field is missing")
		return result, err
	} else if (result.Username == "") {
		err := errors.New("error: The username field is missing")
		return result, err
	} else if (result.Mail == "") {
		err := errors.New("error: The mail field is missing")
		return result, err
	} else if (result.Version == "") {
		err := errors.New("error: The version field is missing")
		return result, err
	}

	result.Score = 0
	result.IP = net.ParseIP(ip)

	return result, nil
}

func (user Result) To_Json() ([]byte, error) {
	jsonBytes, err := json.Marshal(user)
	return jsonBytes, err
}