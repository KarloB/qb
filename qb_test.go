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
	query := "insert into test (col1, col2, col3) values"

	request = append(request,
		TestStruct{Id: 1, Name: "a", Other: "haha"},
		TestStruct{Other: "Hehe"},
		TestStruct{Id: 3, Name: "c", Other: "hihi"},
		TestStruct{Id: 4, Other: "hoho"},
		TestStruct{Id: 5, Name: "e", Other: "huhu"},
		TestStruct{Id: 6, Name: "f", Other: "=U="},
	)

	q, args, err := BulkInsert(query, request)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query: ", q)

	for i, a := range args {
		fmt.Println(i+1, a)
	}
}
