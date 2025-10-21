package gormgoqrius

import (
	"reflect"
	"testing"

	"github.com/golaxo/goqrius/lexer"
	"github.com/golaxo/goqrius/parser"
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expr := WhereClause(parse(tt.filter))
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

func TestGetClause_Fallback_ColumnToColumn(t *testing.T) {
	t.Parallel()

	expr := WhereClause(parse("age gt otherAge"))
	sql, vars := renderWhere(t, expr)
	// fallback should produce a raw expression comparing the two columns
	wantSQL := "WHERE `age` > `otherAge`"
	if sql != wantSQL {
		t.Fatalf("unexpected SQL. got=%q want=%q", sql, wantSQL)
	}

	if len(vars) != 0 {
		t.Fatalf("unexpected vars: %#v", vars)
	}
}

func TestStatementModifier_AddsWhere(t *testing.T) {
	t.Parallel()

	ast := parse("age ge 21")
	modifier := gormExpression{e: ast}
	// Build statement and apply modifier
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to open gorm sqlite dryrun: %v", err)
	}

	stmt := &gorm.Statement{DB: db}
	stmt.Clauses = map[string]clause.Clause{}
	modifier.ModifyStatement(stmt)

	c, ok := stmt.Clauses["WHERE"]
	if !ok {
		t.Fatalf("WHERE clause was not added")
	}

	c.Build(stmt)

	wantSQL := "WHERE `age` >= ?"
	if got := stmt.SQL.String(); got != wantSQL {
		t.Fatalf("unexpected SQL. got=%q want=%q", got, wantSQL)
	}

	wantVars := []any{int64(21)}
	if !reflect.DeepEqual(stmt.Vars, wantVars) {
		t.Fatalf("unexpected vars. got=%#v want=%#v", stmt.Vars, wantVars)
	}
}

// parse is a small helper to get an AST from a filter string.
func parse(input string) parser.Expression {
	l := lexer.New(input)
	p := parser.New(l)

	return p.Parse()
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
