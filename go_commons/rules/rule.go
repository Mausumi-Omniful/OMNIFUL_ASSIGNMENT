package rules

import (
	"errors"
	"fmt"
	"github.com/omniful/go_commons/rules/models"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type DatabaseType uint64

const (
	Postgres DatabaseType = 1
)

type RuleGroup struct {
	Rules []*models.Rule `json:"Rules"`
}

func NewRuleGroup(Rules []*models.Rule) *RuleGroup {
	return &RuleGroup{
		Rules: Rules,
	}
}

func (r *RuleGroup) Scopes(databaseType DatabaseType, RulesName []models.Name, tables map[string]string) (map[DatabaseType]interface{}, error) {
	res := make(map[DatabaseType]interface{}, 0)
	ruleNamesMap := make(map[models.Name]bool, 0)

	for _, v := range RulesName {
		ruleNamesMap[v] = true
	}

	switch databaseType {
	case Postgres:
		ruleQueries := make([]string, 0)
		values := make([]any, 0)
		for _, rule := range r.Rules {
			if _, ok := ruleNamesMap[rule.Name]; !ok {
				continue
			}

			queries := make([]string, 0, len(rule.Conditions))
			for _, condition := range rule.Conditions {
				query, value, err := getQueryAndValues(condition, tables)
				if err != nil {
					return nil, err
				}

				if len(query) > 0 {
					queries = append(queries, query)
					values = append(values, value)
				}
			}

			if len(queries) > 0 {
				switch rule.Operator {
				case models.And:
					ruleQueries = append(ruleQueries, strings.Join(queries, fmt.Sprintf(" %s ", "AND")))
				case models.Or:
					ruleQueries = append(ruleQueries, strings.Join(queries, fmt.Sprintf(" %s ", "OR")))
				default:
					return nil, errors.New(fmt.Sprintf("incorrect seperator ::%d", rule.Operator))
				}
			}
		}

		query := strings.Join(ruleQueries, fmt.Sprintf(") %s (", "AND"))
		if len(ruleQueries) > 1 {
			query = "(" + query
			query = query + ")"
		}

		scope := func(db *gorm.DB) *gorm.DB {
			return db.Where(query, values...)
		}
		res[Postgres] = scope
	}
	return res, nil
}

func (r *RuleGroup) RuleValid(values map[string]string, RulesName []models.Name) (bool, error) {
	ruleNamesMap := make(map[models.Name]bool, 0)
	for _, v := range RulesName {
		ruleNamesMap[v] = true
	}

	isValid := true
	for _, rule := range r.Rules {
		if _, ok := ruleNamesMap[rule.Name]; !ok {
			continue
		}

		isValidAnd := true
		isValidOr := false
		for _, condition := range rule.Conditions {
			filterConditionSatisfied, err := isFilterConditionSatisfied(condition, values)
			if err != nil {
				return false, err
			}

			if rule.Operator == models.And && !filterConditionSatisfied {
				isValidAnd = false
			}

			if rule.Operator == models.Or && filterConditionSatisfied {
				isValidOr = true
			}
		}

		if rule.Operator == models.And {
			isValid = isValid && isValidAnd
		} else if rule.Operator == models.Or {
			isValid = isValid && isValidOr
		}

		if !isValid {
			return false, nil
		}
	}
	return true, nil
}

func isFilterConditionSatisfied(condition models.Condition, params map[string]string) (bool, error) {
	key := condition.Key
	valueReceivedStr, ok := params[key]
	if !ok {
		return false, nil
	}

	if condition.Operator == models.In || condition.Operator == models.NotIn {
		switch condition.DataType {
		case models.Int:
			valueReceived, err := strconv.ParseInt(strings.TrimSpace(valueReceivedStr), 10, 64)
			if err != nil {
				return false, err
			}

			for _, allowedValueStr := range condition.Values {
				allowedValue, cusErr := strconv.ParseInt(strings.TrimSpace(allowedValueStr), 10, 64)
				if cusErr != nil {
					return false, cusErr
				}

				if allowedValue == valueReceived {
					return true, nil
				}
			}
		case models.String:
			for _, allowedValueStr := range condition.Values {
				allowedValue := strings.TrimSpace(allowedValueStr)
				if allowedValue == valueReceivedStr {
					return true, nil
				}
			}
		case models.Float:
			valueReceived, cusErr := strconv.ParseFloat(strings.TrimSpace(valueReceivedStr), 64)
			if cusErr != nil {
				return false, cusErr
			}

			for _, allowedValueStr := range condition.Values {
				allowedValue, parseErr := strconv.ParseFloat(strings.TrimSpace(allowedValueStr), 64)
				if cusErr != nil {
					return false, parseErr
				}

				if allowedValue == valueReceived {
					return true, nil
				}
			}
		case 4:
			valueReceived, cusErr := strconv.ParseBool(strings.TrimSpace(valueReceivedStr))
			if cusErr != nil {
				return false, cusErr
			}

			for _, allowedValueStr := range condition.Values {
				allowedValue, parseErr := strconv.ParseBool(strings.TrimSpace(allowedValueStr))
				if parseErr != nil {
					return false, parseErr
				}

				if allowedValue == valueReceived {
					return true, nil
				}
			}
		default:
			return false, errors.New("invalid datatype")
		}
		return false, nil
	} else if condition.Operator == models.All {
		return true, nil
	} else {
		if len(condition.Values) != 1 {
			return false, errors.New("invalid input")
		}

		switch condition.DataType {
		case models.Int:
			valueReceived, err := strconv.ParseInt(valueReceivedStr, 10, 64)
			if err != nil {
				return false, err
			}

			allowedValue, err := strconv.ParseInt(condition.Values[0], 10, 64)
			if err != nil {
				return false, err
			}

			switch condition.Operator {
			case models.Equals:
				return valueReceived == allowedValue, nil
			case models.NotEquals:
				return valueReceived != allowedValue, nil
			case models.GreaterThan:
				return valueReceived > allowedValue, nil
			case models.LessThan:
				return valueReceived < allowedValue, nil
			default:
				return false, errors.New("invalid operator")
			}

		case models.String:
			switch condition.Operator {
			case models.Equals:
				return valueReceivedStr == condition.Values[0], nil
			case models.NotEquals:
				return valueReceivedStr != condition.Values[0], nil
			case models.GreaterThan:
				return valueReceivedStr > condition.Values[0], nil
			case models.LessThan:
				return valueReceivedStr < condition.Values[0], nil
			default:
				return false, errors.New("invalid operator")
			}

		case models.Float:
			valueReceived, err := strconv.ParseFloat(valueReceivedStr, 64)
			if err != nil {
				return false, err
			}

			allowedValue, err := strconv.ParseFloat(condition.Values[0], 64)
			if err != nil {
				return false, err
			}

			switch condition.Operator {
			case models.Equals:
				return valueReceived == allowedValue, nil
			case models.NotEquals:
				return valueReceived != allowedValue, nil
			case models.GreaterThan:
				return valueReceived > allowedValue, nil
			case models.LessThan:
				return valueReceived < allowedValue, nil
			default:
				return false, errors.New("invalid operator")
			}
		case models.Bool:
			valueReceived, err := strconv.ParseBool(valueReceivedStr)
			if err != nil {
				return false, err
			}

			allowedValue, err := strconv.ParseBool(condition.Values[0])
			if err != nil {
				return false, err
			}

			switch condition.Operator {
			case models.Equals:
				return valueReceived == allowedValue, nil
			case models.NotEquals:
				return valueReceived != allowedValue, nil
			default:
				return false, errors.New("invalid operator")
			}
		default:
			return false, errors.New("invalid datatype")
		}
	}
}

func getQueryAndValues(condition models.Condition, tables map[string]string) (query string, value any, err error) {
	key, ok := tables[condition.Key]
	if !ok {
		key = condition.Key
	}

	if condition.Operator == models.In || condition.Operator == models.NotIn {
		switch condition.DataType {
		case models.Int:
			values := make([]int64, 0)
			for _, allowedValueStr := range condition.Values {
				allowedValue, parseErr := strconv.ParseInt(strings.TrimSpace(allowedValueStr), 10, 64)
				if parseErr != nil {
					err = parseErr
					return
				}
				values = append(values, allowedValue)
				value = values
			}
		case models.String:
			values := make([]string, 0)
			for _, allowedValueStr := range condition.Values {
				values = append(values, strings.TrimSpace(allowedValueStr))
			}
			value = values
		case models.Float:
			values := make([]float64, 0)
			for _, allowedValueStr := range condition.Values {
				allowedValue, parseErr := strconv.ParseFloat(strings.TrimSpace(allowedValueStr), 64)
				if parseErr != nil {
					err = parseErr
					return
				}
				values = append(values, allowedValue)
			}
		case models.Bool:
			values := make([]bool, 0)
			for _, allowedValueStr := range condition.Values {
				allowedValue, parseErr := strconv.ParseBool(strings.TrimSpace(allowedValueStr))
				if parseErr != nil {
					err = parseErr
					return
				}
				values = append(values, allowedValue)
			}
		default:
			err = errors.New("invalid datatype")
			return
		}
		switch condition.Operator {
		case models.In:
			query = fmt.Sprintf("%s IN ?", key)
			return
		case models.NotIn:
			query = fmt.Sprintf("%s NOT IN ?", key)
			return
		}
	} else if condition.Operator == models.All {
		return
	} else {
		if len(condition.Values) != 1 {
			err = errors.New("invalid input")
			return
		}

		switch condition.DataType {
		case models.Int:
			valueShouldBe, parseErr := strconv.ParseInt(condition.Values[0], 10, 64)
			if parseErr != nil {
				err = parseErr
				return
			}
			value = valueShouldBe

		case models.String:
			value = condition.Values

		case models.Float:
			valueShouldBe, parseErr := strconv.ParseFloat(condition.Values[0], 64)
			if parseErr != nil {
				err = parseErr
				return
			}
			value = valueShouldBe

		case models.Bool:
			valueShouldBe, parseErr := strconv.ParseBool(condition.Values[0])
			if parseErr != nil {
				err = parseErr
				return
			}
			value = valueShouldBe
		default:
			err = errors.New("invalid datatype")
			return
		}

		switch condition.Operator {
		case models.Equals:
			query = fmt.Sprintf("%s = ?", key)
			return
		case models.NotEquals:
			query = fmt.Sprintf("%s != ?", key)
			return
		case models.GreaterThan:
			query = fmt.Sprintf("%s > ?", key)
			return
		case models.LessThan:
			query = fmt.Sprintf("%s < ?", key)
			return
		default:
			err = errors.New("invalid operator")
			return
		}
	}

	err = errors.New("invalid operator")
	return
}
