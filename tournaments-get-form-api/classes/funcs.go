package classes

import (
	"encoding/json"
	"errors"
	"time"
)

// формат входных данных {"tournament_id": 1, "username": "Hypoxie","avatar_url": "", "mail": "hypoxie@example.com", "version": "1.6.54s2", "cost": 451, "steam_id": "hypoxie", "metrics":[{"key":"colonists", "value":4}, {"key":"animals", "value":5}]}
func RegDataFromJson(jsonData []byte, ip string) (Result, error) {
	var result Result
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, err
	}

	if result.TournamentID == 0 {
		err := errors.New("error: The tournament_id field is missing")
		return result, err
	} else if result.Username == "" {
		err := errors.New("error: The username field is missing")
		return result, err
	} else if result.GetterMail == "" {
		err := errors.New("error: The mail field is missing")
		return result, err
	} else if result.Version == "" {
		err := errors.New("error: The version field is missing")
		return result, err
	}

	result.SteamID = result.GetterSteamID
	result.Mail = result.GetterMail

	result.Timestamp = uint64(time.Now().Unix())
	result.ID = 0
	result.Status = 0
	result.Score = 0
	result.IP = ip

	return result, nil
}

// {"name": "Test Tournament", "stop_timestamp": 1757099263, "metadata": ["streams", "comment"], "variables": ["animals", "humans"], "formula": "(animals * 2) + (humans * 3)"}
func TournamentDataFromJson(jsonData []byte, ip string) (Tournament, error) {
	var tournament Tournament
	if err := json.Unmarshal(jsonData, &tournament); err != nil {
		return tournament, err
	}

	if tournament.Name == "" {
		err := errors.New("error: The name field is missing")
		return tournament, err
	} else if tournament.StopTimestamp == 0 {
		err := errors.New("error: The stop_timestamp field is missing")
		return tournament, err
	}

	return tournament, nil
}
