package qb

import (
	"fmt"
	"reflect"
	"strings"
)

// QueryBuilder dynamic select query builder with table definition. Returns query and args
func QueryBuilder(query string, definition []Definition) (string, []interface{}) {

	var tableArgs []tableArg
	var requestArgs []interface{}

	for _, p := range definition {
		res := isZero(p.Value)
		if !res {

			switch p.Value.(type) {
			case string:
				h, ok := p.Value.(string)
				if ok {
					if p.Operator == Like {
						p.Value = fmt.Sprintf("%%%s%%", h)
					}
				}
			}

			requestArgs = append(requestArgs, p.Value)
			tableArgs = append(tableArgs, tableArg{value: p.Column, operator: p.Operator.String()})
		}
	}

	if len(tableArgs) > 0 {
		buildArgs := []string{}
		for i, ta := range tableArgs {
			if i == 0 {
				buildArgs = append(buildArgs, "where", ta.value, ta.operator)
				continue
			}
			buildArgs = append(buildArgs, "and", ta.value, ta.operator)
		}
		query = fmt.Sprintf("%s %s", query, strings.Join(buildArgs, " "))
	}

	return query, requestArgs
}

// isZero check if interface equals zero value of its type
func isZero(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
