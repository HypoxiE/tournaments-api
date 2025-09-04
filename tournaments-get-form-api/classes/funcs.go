package classes

import (
	"github.com/Knetic/govaluate"
)

func (result *Result) CalculateScore(tournament Tournament) (error) {
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