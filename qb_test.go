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

func TestBulkInsert(t *testing.T) {

	var request []interface{}
	query := "insert into test (col1, col2, col3)"

	for i := 0; i < 25000; i++ {
		request = append(request,
			TestStruct{Id: i, Name: "a", Other: "haha"},
		)
	}

	err := BulkInsert(query, request, nil)
	if err != nil {
		panic(err)
	}

}
