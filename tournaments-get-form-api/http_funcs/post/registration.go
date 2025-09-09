package post

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

func Registration(c *gin.Context, manager *database.DataBase) {
	var err error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := classes.RegDataFromJson(body, c.ClientIP())
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Считаем количество очков
	err = user.CalculateScore(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
