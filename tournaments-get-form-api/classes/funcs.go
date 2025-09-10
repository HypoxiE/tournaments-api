package classes

import (
	"encoding/json"
	"errors"
)

// формат входных данных {"tournament_id": 1, "username": "Hypoxie","avatar_url": "", "mail": "hypoxie@example.com", "version": "1.6.54s2", "cost": 451, "steam_id": "hypoxie", "metrics":[{"key":"colonists", "value":4}, {"key":"animals", "value":5}]}
func RegDataFromJson(jsonData []byte, ip string) (Result, error) {
	var raw_result CreateResultInput
	if err := json.Unmarshal(jsonData, &raw_result); err != nil {
		return raw_result.NewResultFromInput(ip), err
	}

	var result Result = raw_result.NewResultFromInput(ip)

	if result.TournamentID == 0 {
		err := errors.New("error: The tournament_id field is missing")
		return result, err
	} else if result.Username == "" {
		err := errors.New("error: The username field is missing")
		return result, err
	} else if result.PublicMail == "" {
		err := errors.New("error: The mail field is missing")
		return result, err
	} else if result.Version == "" {
		err := errors.New("error: The version field is missing")
		return result, err
	}

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

	tournament.ID = 0

	return tournament, nil
}
