package gormgoqrius

import (
	"reflect"
	"testing"

	"github.com/golaxo/goqrius"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestGetClause_SQLLite_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter   string
		wantSQL  string
		wantVars []any
	}{
		"simple eq": {
			filter:   "name eq 'John'",
			wantSQL:  "WHERE `name` = ?",
			wantVars: []any{"John"},
		},
		"simple not": {
			filter:   "not name eq 'John'",
			wantSQL:  "WHERE `name` <> ?",
			wantVars: []any{"John"},
		},
		"simple and condition": {
			filter:   "age gt 18 and age lt 65",
			wantSQL:  "WHERE `age` > ? AND `age` < ?",
			wantVars: []any{int64(18), int64(65)},
		},
		"simple or condition": {
			filter:   "name eq 'John' or age lt 21",
			wantSQL:  "WHERE (`name` = ? OR `age` < ?)",
			wantVars: []any{"John", int64(21)},
		},
		"eq null": {
			filter:   "name eq null",
			wantSQL:  "WHERE `name` IS NULL",
			wantVars: nil,
		},
		"ne null": {
			filter:   "name ne null",
			wantSQL:  "WHERE `name` IS NOT NULL",
			wantVars: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e, err := goqrius.Parse(tt.filter)
			if err != nil {
				t.Fatalf("error not expected: %s", err)
			}

			expr := WhereClause(e)
			if expr == nil {
				t.Fatalf("expected non-nil clause expression")
			}

			sql, vars := renderWhere(t, expr)
			// sqlite quotes identifiers with backticks
			if sql != tt.wantSQL {
				t.Fatalf("unexpected SQL. got=%q want=%q", sql, tt.wantSQL)
			}

			if !reflect.DeepEqual(vars, tt.wantVars) {
				t.Fatalf("unexpected vars. got=%#v want=%#v", vars, tt.wantVars)
			}
		})
	}
}

func TestGetClause_EmptyString(t *testing.T) {
	t.Parallel()

	e, err := goqrius.Parse("")
	if err != nil {
		t.Fatalf("error not expected: %s", err)
	}

	expr := WhereClause(e)
	if expr != nil {
		t.Fatalf("expected nil clause expression")
	}
}

// renderWhere renders the WHERE clause SQL and Vars for a given clause.Expression.
func renderWhere(t *testing.T, expr clause.Expression) (string, []any) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to open gorm sqlite dryrun: %v", err)
	}

	stmt := &gorm.Statement{DB: db}
	stmt.Clauses = map[string]clause.Clause{}
	stmt.AddClause(clause.Where{Exprs: []clause.Expression{expr}})
	c := stmt.Clauses["WHERE"]
	c.Build(stmt)

	return stmt.SQL.String(), stmt.Vars
}
