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

// BulkInsert create insert statement
func BulkInsert(query string, rows []interface{}) (string, []interface{}, error) {
	var err error

	if len(rows) == 0 {
		err = fmt.Errorf("No rows in request")
		return query, nil, err
	}

	placeholder, count := createPlaceholder(rows[0])

	placeholders := make([]string, len(rows))
	args := make([]interface{}, (len(rows) * count))

	var argCount int
	for i, entry := range rows {
		placeholders[i] = placeholder
		v := reflect.ValueOf(entry)
		for y := 0; y < v.NumField(); y++ {
			args[argCount] = v.Field(y).Interface()
			argCount++
		}
	}

	statement := fmt.Sprintf("%s %s", query, strings.Join(placeholders, ","))

	return statement, args, err
}

// isZero check if interface equals zero value of its type
func isZero(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// createPlaceholder create placeholder for one insert on structure
func createPlaceholder(a interface{}) (string, int) {

	instance := reflect.TypeOf(a)
	fCount := instance.NumField()

	var ph []string
	for i := 1; i <= fCount; i++ {
		ph = append(ph, "?")
	}

	placeholder := fmt.Sprintf("(%s)", strings.Join(ph, ","))

	return placeholder, fCount
}
