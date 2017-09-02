package qb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// BulkInsert fast insert for large data set
func BulkInsert(ctx context.Context, query string, rows []interface{}, db *sql.DB) error {
	err := checkInsertRequest(query, rows, db)
	if err != nil {
		return err
	}

	placeholder, fCount, err := createPlaceholder(query, rows[0]) // placeholder create placeholder based on structure. Count fields to determine ideal batch size
	if err != nil {
		return err
	}
	batchSize := len(rows)                             // initial size is length of recieved rows
	maxBatchSize := int(mysqlMaxPlaceholders / fCount) // max batch size can not have over 65536 placeholders. Limitation by MySQL
	if batchSize > maxBatchSize {                      //if it does...
		batchSize = findBatchSize(batchSize, maxBatchSize) // find largest possible batch size that doesn't exceed max number of placeholders
	}

	chunks := ChunkIt(rows, batchSize) // split dataset into chunks

	for i, chunk := range chunks {
		statement, args, err := CreateStatement(query, chunk, placeholder, fCount)
		if err != nil {
			return fmt.Errorf(errors[statementError], err)
		}
		insertInfo(ctx, i, len(chunk))
		_, err = db.Exec(statement, args...)
		if err != nil {
			return fmt.Errorf(errors[insertError], err)
		}
	}
	return nil
}

// QueryBuilder dynamic select query builder with table definition. Returns query and args
func QueryBuilder(query string, definition []Definition) (string, []interface{}) {
	var tableArgs []tableArg
	var requestArgs []interface{}

	query = cleanQueryString(query)

	for _, p := range definition {
		var counter int
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

			switch p.Operator {
			case In:
				h, ok := p.Value.(string)
				if ok {
					values := strings.Split(h, " ")
					values = cleanSlice(values)
					counter = len(values)
					for _, v := range values {
						requestArgs = append(requestArgs, v)
					}
				}

				hs, ok := p.Value.([]string)
				if ok {
					counter = len(hs)
					for i := range hs {
						requestArgs = append(requestArgs, hs[i])
					}
				}
			case Like:
				h, ok := p.Value.(string)
				if ok {
					if p.Operator == Like {
						values := strings.Split(h, " ")
						values = cleanSlice(values)
						counter = len(values)
						for _, v := range values {
							v = fmt.Sprintf("%%%s%%", v)
							requestArgs = append(requestArgs, v)
						}
					}
				}

			default:
				requestArgs = append(requestArgs, p.Value)
			}

			op := buildOperator(p.Operator, counter)

			tableArgs = append(tableArgs, tableArg{value: p.Column, operator: op})
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
