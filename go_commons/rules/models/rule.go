package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type RuleOperator uint64
type ConditionOperator uint64
type DataType uint64
type Name uint64
type EntityType uint64
type Rules []*Rule

type (
	Rule struct {
		ID         uint64 `gorm:"primaryKey"`
		TenantID   null.Int
		EntityID   uint64       `gorm:"type:not null"`
		EntityType EntityType   `gorm:"type: not null"`
		Name       Name         `gorm:"type: not null"`
		Conditions []Condition  `gorm:"serializer:json"`
		Operator   RuleOperator `gorm:"type: not null"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  null.Time
	}

	Condition struct {
		Key      string            `gorm:"not null" json:"key"`
		Operator ConditionOperator `gorm:"not null" json:"operator"`
		DataType DataType          `gorm:"not null" json:"data_type"`
		Values   []string          `json:"values"`
	}
)

const (
	And RuleOperator = 1
	Or  RuleOperator = 2
)

const (
	Equals      ConditionOperator = 1
	NotEquals   ConditionOperator = 2
	GreaterThan ConditionOperator = 3
	LessThan    ConditionOperator = 4
	In          ConditionOperator = 5
	NotIn       ConditionOperator = 6
	All         ConditionOperator = 7
)

const (
	Int    DataType = 1
	String DataType = 2
	Float  DataType = 3
	Bool   DataType = 4
)

const (
	UserHub        Name = 1
	OrderStatus    Name = 2
	Seller         Name = 3
	Tenant         Name = 4
	ShippingClient Name = 5
	SortingHub     Name = 6
)

const (
	TenantUser  EntityType = 1
	OmnifulUser EntityType = 2
)

//var OperatorToRuleOperator = map[string]RuleOperator{
//	"AND": 1,
//	"OR":  2,
//}
//
//var OperatorToConditionOperator = map[string]ConditionOperator{
//	"EQUAL":        1,
//	"NOT_EQUAL":    2,
//	"GREATER_THAN": 3,
//	"LESS_THAN":    4,
//	"IN":           5,
//	"NOT_IN":       6,
//	"ALL":          7,
//}
//
//var OperatorToDataType = map[string]DataType{
//	"INT":    1,
//	"STRING": 2,
//	"FLOAT":  3,
//	"BOOL":   4,
//}

func (rules Rules) GroupRulesByRuleIDs() (ruleIDToRuleMap map[uint64]*Rule) {
	ruleIDToRuleMap = make(map[uint64]*Rule)
	for _, rule := range rules {
		ruleIDToRuleMap[rule.ID] = rule
	}
	return
}
