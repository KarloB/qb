package qb

import (
	"fmt"
	"reflect"
	"strings"
)

// ChunkIt split slice into slices of slice based on batch size!
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
