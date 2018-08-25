package qb

import "fmt"

const (
	mysqlMaxPlaceholders = 65536
	defaultPlaceholder   = "?"
)

type e int

const (
	requestEmpty e = iota + 1
	insertError
	statementError
	databaseError
	queryError
)

var errors = map[e]string{
	requestEmpty:   "Request is empty",
	insertError:    "Error inserting to database. %v",
	statementError: "Error creating statement. %v",
	databaseError:  "Error connecting to database",
	queryError:     "Query template not provided",
}

// Definition define value, column name and statement type for table columns
type Definition struct {
	Value       interface{}
	Column      string
	Operator    Operator
	Placeholder string
}

// ColStatement define which sql statement will each column use
type Operator int

const (
	Equals Operator = iota + 1
	NotEquals
	Like
	Between
	Greater
	Lesser
	// support string, []string and []int
	In
	// Or operator for slice of string or slice of int
	Or
)

func (t Operator) String() string {
	switch t {
	case Equals:
		return "= ?"
	case NotEquals:
		return "!= ?"
	case Like:
		return "like ?"
	case Between:
		return "between"
	case Greater:
		return ">= ?"
	case Lesser:
		return "<= ?"
	case In:
		return "in (?)"
	case Or:
		return ""
	default:
		return "= ?"
	}
}

func (t Operator) WithPlaceholder(placeholder string) string {
	if len(placeholder) == 0 {
		placeholder = defaultPlaceholder
	}
	switch t {
	case Equals:
		return fmt.Sprintf("= %s", placeholder)
	case NotEquals:
		return fmt.Sprintf("!= %s", placeholder)
	case Like:
		return fmt.Sprintf("like %s", placeholder)
	case Between:
		return "between"
	case Greater:
		return fmt.Sprintf(">= %s", placeholder)
	case Lesser:
		return fmt.Sprintf("<= %s", placeholder)
	case In:
		return fmt.Sprintf("in (%s)", placeholder)
	case Or:
		return ""
	default:
		return fmt.Sprintf("= %s", placeholder)
	}
}

type tableArg struct {
	value    string
	operator string
}
