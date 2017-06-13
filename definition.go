package qb

// Definition define value, column name and statement type for table columns
type Definition struct {
	Value    interface{}
	Column   string
	Operator Operator
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
	default:
		return "= ?"
	}
}

type tableArg struct {
	value    string
	operator string
}
