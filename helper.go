package qb

import (
	"database/sql"
	"fmt"
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

	statement := fmt.Sprintf("%s values %s", query, strings.Join(placeholders, ","))

	return statement, args, nil
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
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// createPlaceholder create placeholder for one insert on structure
func createPlaceholder(a interface{}) (string, int) {

	instance := reflect.TypeOf(a)
	fCount := instance.NumField()

	ph := make([]string, fCount)
	for i := 0; i < fCount; i++ {
		ph[i] = "?"
	}

	placeholder := fmt.Sprintf("(%s)", strings.Join(ph, ","))

	return placeholder, fCount
}

func insertInfo(i int) {
	fmt.Printf("Inserting batch %d\n", i+1)
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
