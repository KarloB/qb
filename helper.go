package qb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// ChunkIt split slice into slices of slice based on batch size
func ChunkIt(rows []interface{}, insertBatchSize int) [][]interface{} {
	var result [][]interface{}

	rowLen := len(rows)
	rowChunk := insertBatchSize

	if rowLen > rowChunk {

		for i := 0; i < len(rows); i += rowChunk {

			end := i + rowChunk
			if end > len(rows) {
				end = len(rows)
			}

			result = append(result, rows[i:end])
		}
	} else {
		result = append(result, rows)
	}

	return result
}

// CreateStatement create insert statement for large data set
func CreateStatement(query string, rows []interface{}, placeholder string, count int) (string, []interface{}, error) {
	var err error
	if len(placeholder) == 0 && count == 0 {
		placeholder, count, err = createPlaceholder(query, rows[0])
		if err != nil {
			return "", nil, err
		}
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

	query = queryValues(query)
	statement := fmt.Sprintf("%s %s", query, strings.Join(placeholders, ","))

	return statement, args, nil
}

func queryValues(query string) string {
	query = strings.ToLower(query)
	valuesIndex := strings.Index(query, "values")
	if valuesIndex > 0 {
		query = query[:valuesIndex] // delete placeholders if any exist
	}

	query = fmt.Sprintf("%s values", query)
	return query
}

func findBatchSize(a int, limit int) int {
	var result int

	i := 1
	for {
		result = int(a / i)
		if result < limit {
			break
		}
		i++
	}

	return result
}

// isZero check if interface equals zero value of its type
func isZero(x interface{}) bool {
	var result bool
	switch x.(type) {
	case []string:
		h, ok := x.([]string)
		if ok {
			if h == nil {
				return true
			}
			if len(h) == 0 {
				return true
			}
		}
	case []int:
		h, ok := x.([]int)
		if ok {
			if h == nil {
				return true
			}
			if len(h) == 0 {
				return true
			}
		}
	default:
		result = (x == reflect.Zero(reflect.TypeOf(x)).Interface())
	}
	return result
}

// createPlaceholder create placeholder for one insert on structure
func createPlaceholder(query string, a interface{}) (string, int, error) {

	instance := reflect.TypeOf(a)
	fCount := instance.NumField()

	columns := extractQueryColumns(query)
	if len(columns) != fCount {
		return "", 0, fmt.Errorf("Structure type doesn't match column count")
	}

	ph := make([]string, fCount)
	for i := 0; i < fCount; i++ {
		ph[i] = "?"
	}

	placeholder := fmt.Sprintf("(%s)", strings.Join(ph, ","))

	return placeholder, fCount, nil
}

func extractQueryColumns(query string) []string {
	columnsStart := strings.Index(query, "(")
	columnsEnd := strings.Index(query, ")")
	columnsString := query[columnsStart+1 : columnsEnd]
	columnsString = strings.Replace(columnsString, " ", "", -1)
	columns := strings.Split(columnsString, ",")
	return columns
}

func insertInfo(ctx context.Context, i, chunk int) {
	switch i {
	case 0:
		log.Printf("[BulkInsert] %v Inserting %d entries", ctx.Value(""), chunk)
	default:
		log.Printf("[BulkInsert] %v Insert batch %d. Entries: %d", ctx.Value(""), i, chunk)
	}
}

func checkInsertRequest(query string, rows []interface{}, db *sql.DB) error {
	if len(rows) == 0 {
		return fmt.Errorf(errors[requestEmpty])
	}
	if len(query) == 0 {
		return fmt.Errorf(errors[queryError])
	}
	if db == nil {
		return fmt.Errorf(errors[databaseError])
	}

	return nil
}

// cleanSlice remove empty strings from string slice
func cleanSlice(a []string) []string {
	var result []string
	for _, b := range a {
		if len(strings.Replace(b, " ", "", -1)) == 0 {
			continue
		}
		result = append(result, b)
	}
	return result
}

// buildOperator "in"" operator can have multiple argumens as value
func buildOperator(operator Operator, counter int) string {
	op := operator.String()

	if operator == In && counter > 1 {
		var newOperator []string
		for i := 1; i <= counter; i++ {
			newOperator = append(newOperator, "?")
		}
		op = fmt.Sprintf("in (%s)", strings.Join(newOperator, ","))
	}

	if operator == Like && counter > 1 {
		var ors []string
		for i := 0; i < counter; i++ {
			ors[i] = "(?)"
		}
		newOperator := "like " + strings.Join(ors, " or ")
		op = newOperator
	}

	if operator == Or {
		var ors []string
		for i := 0; i < counter; i++ {
			ors[i] = "(?)"
		}
		newOperator := "= " + strings.Join(ors, " or ")
		op = newOperator
	}

	return op
}

func cleanQueryString(query string) string {
	query = strings.ToLower(query)
	index := strings.Index(query, " where ")
	if index < 0 {
		return query
	}

	query = query[0:index]
	return query
}
