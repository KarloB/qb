package qb

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestQueryBuilder(t *testing.T) {

	query := "select u.id, u.name, u.email, u.registered, u.active_from from user u"

	request := []Definition{
		{[]string{"John", "Milkovocih", "Pimpek"}, "u.name", In},
	}

	result, args := QueryBuilder(query, request)

	fmt.Println(result)
	fmt.Println(args)
}

type TestStruct struct {
	Id    int
	Name  string
	Other string
}

const query = "insert into test (col1, col2, col3)"

// createRequest with large data set
func createRequest(size int) []interface{} {
	result := make([]interface{}, size)
	for i := 0; i < size; i++ {
		result[i] = TestStruct{Id: i, Name: "a", Other: "haha"}
	}
	return result
}

func TestBulkInsert(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	size := 500000
	request := createRequest(size) // create dummy request with large data set

	// find out how many batches will be created
	instance := reflect.TypeOf(request[0])
	fCount := instance.NumField()
	maxBatchSize := int(mysqlMaxPlaceholders / fCount)
	if size > maxBatchSize {
		size = findBatchSize(size, maxBatchSize)
	}
	chunks := ChunkIt(request, size)

	for i := 0; i < len(chunks); i++ {
		mock.ExpectExec(escape(query)).WillReturnResult(sqlmock.NewResult(0, int64(size)))
	}

	ctx := context.WithValue(context.Background(), "grah", "kupus")

	err = BulkInsert(ctx, query, request, db)
	if err != nil {
		panic(err)
	}
}

func TestBulkInsertWithCustomPlaceholder(t *testing.T) {

	type testStruct struct {
		Id    string
		Name  string `qb:"placeholder:uuid_to_bin(?,true)"`
		Value string `qb:"placeholder:uuid_to_bin(?,true)"`
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	size := 10000

	createRequest := func(size int) []interface{} {
		result := make([]interface{}, size)
		for i := 0; i < size; i++ {
			result[i] = testStruct{Id: uuid.New().String(), Name: "a", Value: "haha"}
		}
		return result
	}

	request := createRequest(size) // create dummy request with large dataset

	// find out how many batches will be created
	instance := reflect.TypeOf(request[0])
	fCount := instance.NumField()
	maxBatchSize := int(mysqlMaxPlaceholders / fCount)
	if size > maxBatchSize {
		size = findBatchSize(size, maxBatchSize)
	}
	chunks := ChunkIt(request, size)

	for i := 0; i < len(chunks); i++ {
		mock.ExpectExec(escape(query)).WillReturnResult(sqlmock.NewResult(0, int64(size)))
	}

	ctx := context.WithValue(context.Background(), "grah", "kupus")

	err = BulkInsert(ctx, query, request, db)
	if err != nil {
		panic(err)
	}
}

func TestCreateStatement(t *testing.T) {

	request := createRequest(50)

	statement, args, err := CreateStatement(query, request, "", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)
	fmt.Println(args)
}
func escape(query string) string {
	chars := []string{`(`, `)`, `$`, `+`, `?`, `.`}
	for _, r := range chars {
		query = strings.Replace(query, r, `\`+r, -1)
	}
	return query
}
