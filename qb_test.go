package qb

import (
	"context"
	"fmt"
	"testing"
)

func TestQueryBuilder(t *testing.T) {

	query := "select u.id, u.name, u.email, u.registered, u.active_from from user u join photo p on (a+b where id in select a + b where c = 4) where id = ?"

	request := []Definition{
		{[]string{"John", "Milkovocih"}, "u.name", Or},
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

	ctx := context.WithValue(context.Background(), "grah", "kupus")

	err := BulkInsert(ctx, query, request, nil)
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
