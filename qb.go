package qb

import (
	"database/sql"
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

// BulkInsert fast insert for large data set
func BulkInsert(query string, rows []interface{}, db *sql.DB) error {
	var err error
	if len(rows) == 0 {
		err = fmt.Errorf("No rows in request")
		return err
	}

	fmt.Println("db", db)

	placeholder, fCount := createPlaceholder(rows[0])  // placeholder create placeholder based on structure. Count fields to determine ideal batch size
	batchSize := len(rows)                             // initial size is length of recieved rows
	maxBatchSize := int(mysqlMaxPlaceholders / fCount) // max batch size can not have over 65536 placeholders. Limitation by MySQL
	if batchSize > maxBatchSize {                      //if it does...
		batchSize = findBatchSize(batchSize, maxBatchSize) // find largest possible batch size that doesn't exceed max number of placeholders
	}

	chunks := ChunkIt(rows, batchSize) // split dataset into chunks

	for _, chunk := range chunks {
		statement, args, err := CreateStatement(query, chunk, placeholder, fCount)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(statement, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateStatement create insert statement for large data set
func CreateStatement(query string, rows []interface{}, placeholder string, count int) (string, []interface{}, error) {
	var err error

	if len(placeholder) == 0 && count == 0 {
		placeholder, count = createPlaceholder(rows[0])
	}

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

	statement := fmt.Sprintf("%s VALUES %s", query, strings.Join(placeholders, ","))

	return statement, args, err
}
