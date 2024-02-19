package result

import (
	"math/rand"
)

type Scorer func(map[string]string) int

type Score struct {
	Scorer         string   `mapstructure:"scorer"`
	MetadataFields []string `mapstructure:"metadata_fields"`
	Rules          []Rule   `mapstructure:"rules"`
}

var (
	scorerMap = map[string]Scorer{
		"random": scoreRandom,
		"rule":   scoreByRule,
	}
)

// getScorer return a default scorer if no scorer is specified.
func getScorer() Scorer {
	if sr, ok := scorerMap[conf.Score.Scorer]; ok {
		return sr
	}
	return scoreRandom
}

// scoreRandom return random score for each result.
func scoreRandom(data map[string]string) int {
	return rand.Intn(100)
}

// scoreByRule return total rules score for each result.
func scoreByRule(data map[string]string) int {
	var sum int
	for _, rule := range rules {
		if rule.match(data) {
			sum += rule.Score
		}
	}
	return sum
}
