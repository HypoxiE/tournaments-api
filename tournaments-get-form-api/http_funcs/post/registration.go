package post

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"tournaments-api/classes"
	"tournaments-api/database"
	"tournaments-api/funcs"

	"github.com/gin-gonic/gin"
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
type CreateMetricInput struct {
	ResultID uint   `json:"result_id"`
	Key      string `json:"key"`
	Value    *int64 `json:"value"`
}
type CreateMetadataInput struct {
	ResultID uint    `json:"result_id"`
	Key      string  `json:"key"`
	Value    *string `json:"value"`
}

func (input CreateResultInput) NewResultFromInput(ip string) classes.Result {
	var result = classes.Result{
		TournamentID: input.TournamentID,

		Username:  input.Username,
		Avatar:    funcs.OrElsePtr(input.Avatar, ""),
		Version:   input.Version,
		Score:     funcs.GetPtr(0),
		Penalty:   funcs.GetPtr(0),
		Cost:      funcs.OrElsePtr(input.Cost, 0),
		Status:    funcs.GetPtr(1),
		Timestamp: uint64(time.Now().Unix()),

		PublicSteamID: funcs.OrElsePtr(input.PublicSteamID, ""),
		SteamID:       funcs.OrElsePtr(input.PublicSteamID, ""),
		PublicMail:    input.PublicMail,
		Mail:          input.PublicMail,
		PublicIP:      ip,
		IP:            ip,
	}

	for _, v := range input.Metrics {
		result.Metrics = append(result.Metrics, v.NewMetricFromInput())
	}
	for _, v := range input.Metadata {
		result.Metadata = append(result.Metadata, v.NewMetadataFromInput())
	}

	return result
}
func (input CreateMetricInput) NewMetricFromInput() classes.Metric {
	return classes.Metric{
		ResultID: input.ResultID,
		Key:      input.Key,
		Value:    funcs.OrElsePtr(input.Value, 0),
	}
}
func (input CreateMetadataInput) NewMetadataFromInput() classes.Metadata {
	return classes.Metadata{
		ResultID: input.ResultID,
		Key:      input.Key,
		Value:    funcs.OrElsePtr(input.Value, ""),
	}
}

// формат входных данных {"tournament_id": 1, "username": "Hypoxie","avatar_url": "", "mail": "hypoxie@example.com", "version": "1.6.54s2", "cost": 451, "steam_id": "hypoxie", "metrics":[{"key":"colonists", "value":4}, {"key":"animals", "value":5}]}
func RegDataFromJson(jsonData []byte, ip string) (classes.Result, error) {
	var raw_result CreateResultInput
	if err := json.Unmarshal(jsonData, &raw_result); err != nil {
		return raw_result.NewResultFromInput(ip), err
	}

	var result classes.Result = raw_result.NewResultFromInput(ip)

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

func Registration(c *gin.Context, manager *database.DataBase) {
	var err error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := RegDataFromJson(body, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tournament classes.Tournament
	err = manager.DataBase.First(&tournament, user.TournamentID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка на сроки отправки

	if user.Timestamp < tournament.StartTimestamp {
		err = fmt.Errorf("error: the tournament has not started yet")
	}
	if user.Timestamp > tournament.StopTimestamp {
		err = fmt.Errorf("error: the tournament has already ended")
	}

	// Проверка метрик
	find_metric_key := func(metric classes.Metric) bool {
		for _, variable := range tournament.Variables {
			if variable == metric.Key {
				return true
			}
		}
		return false
	}
	var used_metrics []string
	for _, metric := range user.Metrics {
		if !find_metric_key(metric) {
			err = fmt.Errorf("error: unknown metric %v", metric.Key)
		}
		for _, v := range used_metrics {
			if metric.Key == v {
				err = fmt.Errorf("error: metric %v duplicated", metric.Key)
			}
		}
		used_metrics = append(used_metrics, metric.Key)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка метаданных
	find_meta_key := func(meta classes.Metadata) bool {
		for _, variable := range tournament.Metadata {
			if variable == meta.Key {
				return true
			}
		}
		return false
	}
	var used_metadata []string
	for _, meta := range user.Metadata {
		if !find_meta_key(meta) {
			err = fmt.Errorf("error: unknown metadata %v", meta.Key)
		}
		for _, v := range used_metadata {
			if meta.Key == v {
				err = fmt.Errorf("error: metadata %v duplicated", meta.Key)
			}
		}
		used_metadata = append(used_metrics, meta.Key)
	}
	find_key_in_meta := func(meta string) bool {
		for _, variable := range user.Metadata {
			if variable.Key == meta {
				return true
			}
		}
		return false
	}
	for _, meta := range tournament.Metadata {
		if !find_key_in_meta(meta) {
			err = fmt.Errorf("error: missing metadata with the key %v", meta)
		}
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Считаем количество очков
	err = user.CalculateScore(tournament)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := manager.DataBase.Create(&user)
	if result.Error != nil {
		log.Fatalf("[ERROR] failed to insert user: %v", result.Error)
	}

	jsonBytes, err := user.ToJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[INFO] New registration: %v!", user.Username)

	c.Data(http.StatusOK, "application/json", jsonBytes)
}
