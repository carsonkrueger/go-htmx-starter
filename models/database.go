package models

import "github.com/go-jet/jet/v2/postgres"

type SearchParams struct {
	Where   postgres.BoolExpression
	OrderBy []postgres.OrderByClause
	Offset  *int64
	Limit   *int64
}
