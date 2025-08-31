package mybaits

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Statement struct {
	sql string
}

func (s *Statement) formatSQL() (formatted string, err error) {
	var stmt sqlparser.Statement

	if stmt, err = sqlparser.Parse(s.sql); err != nil {
		return
	}
	formatted = sqlparser.String(stmt)

	return
}
