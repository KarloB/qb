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
