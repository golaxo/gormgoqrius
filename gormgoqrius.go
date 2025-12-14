package gormgoqrius

import (
	"strconv"

	"github.com/golaxo/goqrius"
	"gorm.io/gorm/clause"
)

// WhereClause transform a goqrius.Expression into a GORM clause.
// The return value can be null, meaning there was no query.
func WhereClause(e goqrius.Expression) clause.Expression {
	return toClauseExpression(e)
}

// toClauseExpression converts a parser.Expression into a GORM clause.Expression
// supported operations: and, or, not, eq, ne, gt, ge, lt, le.
//
//nolint:gocognit,funlen // refactor later.
func toClauseExpression(e goqrius.Expression) clause.Expression {
	switch v := e.(type) {
	case *goqrius.AndExpr:
		left := toClauseExpression(v.Left)
		right := toClauseExpression(v.Right)
		return clause.And(left, right)
	case *goqrius.OrExpr:
		left := toClauseExpression(v.Left)
		right := toClauseExpression(v.Right)
		return clause.Or(left, right)
	case *goqrius.FilterExpr:
		//nolint:exhaustive // no need to check every case.
		switch v.Operator {
		case goqrius.Eq:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Eq{Column: col, Value: val}
			}

			return nil
		case goqrius.NotEq:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Neq{Column: col, Value: val}
			}

			return nil
		case goqrius.GreaterThan:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Gt{Column: col, Value: val}
			}

			return nil
		case goqrius.GreaterThanOrEqual:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Gte{Column: col, Value: val}
			}

			return nil
		case goqrius.LessThan:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Lt{Column: col, Value: val}
			}

			return nil
		case goqrius.LessThanOrEqual:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Lte{Column: col, Value: val}
			}

			return nil
		default:
			return nil
		}
	case *goqrius.NotExpr:
		inner := toClauseExpression(v.Right)
		if inner == nil {
			return nil
		}

		return clause.Not(inner)
	case *goqrius.Identifier:
		return clause.Expr{SQL: "?", Vars: []any{clause.Column{Name: v.Value}}}
	case *goqrius.IntegerLiteral:
		return clause.Expr{SQL: "?", Vars: []any{parseInt(v.Value)}}
	case *goqrius.StringLiteral:
		return clause.Expr{SQL: "?", Vars: []any{v.Value}}
	case *goqrius.Null:
		return clause.Expr{SQL: "NULL"}
	default:
		return nil
	}
}

//nolint:nonamedreturns // named returns make sense here
func splitComparisonSides(lid *goqrius.Identifier, right goqrius.Value) (column, value any, ok bool) {
	column = clause.Column{Name: lid.Value}
	// Right can be int, string or null
	switch rv := right.(type) {
	case *goqrius.IntegerLiteral:
		value = parseInt(rv.Value)

		return column, value, true
	case *goqrius.StringLiteral:
		value = rv.Value

		return column, value, true
	case *goqrius.Null:
		// Let GORM render IS NULL / IS NOT NULL
		value = nil

		return column, value, true
	default:
		return nil, nil, false
	}
}

func parseInt(s string) any {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return u
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	return s
}
