package qb

const mysqlMaxPlaceholders = 65536

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
	In
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

type tableArg struct {
	value    string
	operator string
}
