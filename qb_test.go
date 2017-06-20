package qb

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryBuilder(t *testing.T) {

	query := "select u.id, u.name, u.email, u.registered, u.active_from from user u join photo p on (p.user_id = u.id)"

	request := []Definition{
		{"John", "u.name", Like},
		{1, "u.id", Equals},
		{"", "u.email", Equals},
		{time.Now(), "u.registered", Lesser},
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

	request := createRequest(500000) // create dummy request with large data set

	err := BulkInsert(query, request, nil)
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
