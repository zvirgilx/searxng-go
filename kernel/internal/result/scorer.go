package result

import (
	"math/rand"
	"regexp"
	"slices"
	"strings"

	"github.com/zvirgilx/searxng-go/kernel/config"
)

type scorer func(data Data) int

var enableScorer scorer
var rules []config.Rule

var scorerMap = map[string]scorer{
	"random": ScoreRandom,
	"weight": ScoreByWeight,
	"rule":   ScoreByRule,
}

func InitScorer() {
	enableScorer = ScoreRandom
	if sr, ok := scorerMap[config.Conf.Result.Score.Scorer]; ok {
		enableScorer = sr
	}

	rs := config.Conf.Result.Score.Rules
	rules = make([]config.Rule, 0, len(rs))
	for _, rule := range rs {
		if !rule.Enable {
			continue
		}
		rules = append(rules, rule)
	}
}

func getScore(data Data) int {
	return enableScorer(data)
}

// ScoreByWeight return weight as score for each result.
// if weight map is null, use random score.
func ScoreByWeight(data Data) int {
	weightMap := config.Conf.Result.Score.Weight
	if weightMap == nil {
		return ScoreRandom(data)
	}
	return weightMap[data.Engine]
}

// ScoreRandom return random score for each result.
func ScoreRandom(data Data) int {
	return rand.Intn(100)
}

// ScoreByRule return total rules score for each result.
func ScoreByRule(data Data) int {
	var sum int
	for _, rule := range rules {
		if matchConditions(data.Metadata, rule.Conditions) {
			sum += rule.Score
		}
	}
	return sum
}

func matchConditions(metadata map[string]string, conditions []config.Condition) bool {
	for _, cond := range conditions {
		if metadata[cond.Field] == "" {
			return false
		}

		// replace the variable in condition values.
		replaceVariable(cond.Values, metadata)

		if !match(metadata[cond.Field], cond) {
			return false
		}
	}
	return true
}

func match(value string, condition config.Condition) bool {
	switch condition.Operator {

	case "in": // whether value exists in condition.values. if true when "a" in ["a", "b", "c"].
		return slices.Contains(condition.Values, value)
	case "containAny": // whether condition.values is a substring of value. if true when "aBC" containAny ["a", "b"].
		for _, sub := range condition.Values {
			if strings.Contains(value, sub) {
				return true
			}
		}
	}
	return false
}

// match the variable in conditions values.
var variableMatcher = regexp.MustCompile(`^\$[A-Z_][A-Z0-9_]*$`)

// replace variable with real value.
func replaceVariable(origin []string, metadata map[string]string) {
	for i := range origin {
		// the variable value will replace by real value.
		// e.g. $QUERY -> query(query from search).
		if variableMatcher.MatchString(origin[i]) {
			origin[i] = metadata[origin[i]]
		}
	}
}
