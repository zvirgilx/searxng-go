package result

import (
	"slices"
	"strings"
)

// Rule is a way that scores search results depends on conditions.
type Rule struct {
	Name       string      `mapstructure:"name"`
	Enable     bool        `mapstructure:"enable"`
	Score      int         `mapstructure:"score"`
	Conditions []Condition `mapstructure:"conditions"`
}

// Condition is to determine whether the search results meet some requirements.
type Condition struct {
	Field    string   `mapstructure:"field"`    // result field name, like title, description,...
	Operator string   `mapstructure:"operator"` // operator, like in, containAny
	Expects  []string `mapstructure:"expects"`  // expect values
}

var rules []Rule

func loadRule() {
	rs := conf.Score.Rules
	rules = make([]Rule, 0, len(rs))
	for _, rule := range rs {
		if !rule.Enable {
			continue
		}
		rules = append(rules, rule)
	}
}

// matched only if all conditions under this rule are met.
func (r *Rule) match(data map[string]string) bool {
	for _, cond := range r.Conditions {
		if !cond.match(data) {
			return false
		}
	}
	return true
}

func (c *Condition) match(data map[string]string) bool {
	// if data don't have the field, just return false.
	if data[c.Field] == "" {
		return false
	}

	// replace the variable in condition values.
	replaceVariable(c.Expects, data)

	switch c.Operator {
	case "in": // whether v exists in condition.expects. if true when "a" in ["a", "b", "c"].
		return slices.Contains(c.Expects, data[c.Field])
	case "containAny": // whether condition.expects is a substring of v. if true when "aBC" containAny ["a", "b"].
		for _, sub := range c.Expects {
			if strings.Contains(data[c.Field], sub) {
				return true
			}
		}
	}
	return false
}
