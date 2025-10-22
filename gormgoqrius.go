package gormgoqrius

import (
	"strconv"

	"github.com/golaxo/goqrius"
	"github.com/golaxo/goqrius/token"
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
	case *goqrius.FilterExpr:
		left := toClauseExpression(v.Left)

		right := toClauseExpression(v.Right)
		//nolint:exhaustive // no need to check every case.
		switch v.Operator {
		case token.And:
			return clause.And(left, right)
		case token.Or:
			return clause.Or(left, right)
		case token.Eq:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Eq{Column: col, Value: val}
			}

			return fallbackBinary(token.Eq, v.Left, v.Right)
		case token.NotEq:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Neq{Column: col, Value: val}
			}

			return fallbackBinary(token.NotEq, v.Left, v.Right)
		case token.GreaterThan:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Gt{Column: col, Value: val}
			}

			return fallbackBinary(token.GreaterThan, v.Left, v.Right)
		case token.GreaterThanOrEqual:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Gte{Column: col, Value: val}
			}

			return fallbackBinary(token.GreaterThanOrEqual, v.Left, v.Right)
		case token.LessThan:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Lt{Column: col, Value: val}
			}

			return fallbackBinary(token.LessThan, v.Left, v.Right)
		case token.LessThanOrEqual:
			col, val, ok := splitComparisonSides(v.Left, v.Right)
			if ok {
				return clause.Lte{Column: col, Value: val}
			}

			return fallbackBinary(token.LessThanOrEqual, v.Left, v.Right)
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
	default:
		return nil
	}
}

// splitComparisonSides extracts a left column and right value if possible.
// Returns ok=false when it cannot be represented using standard GORM comparison clauses.
//
//nolint:nonamedreturns // named returns make sense here
func splitComparisonSides(left, right goqrius.Expression) (column, value any, ok bool) {
	// Left must be an identifier (column)
	lid, ok := left.(*goqrius.Identifier)
	if !ok {
		return nil, nil, false
	}

	column = clause.Column{Name: lid.Value}
	// Right can be int, string, or identifier (column-to-column)
	switch rv := right.(type) {
	case *goqrius.IntegerLiteral:
		value = parseInt(rv.Value)

		return column, value, true
	case *goqrius.StringLiteral:
		value = rv.Value

		return column, value, true
	case *goqrius.Identifier:
		// Use raw expression for column-to-column comparisons
		return nil, nil, false
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

// fallbackBinary builds a raw SQL expression for complex cases (e.g., column-to-column).
func fallbackBinary(op token.Type, left, right goqrius.Expression) clause.Expression {
	l := toClauseExpression(left)

	r := toClauseExpression(right)
	if l == nil || r == nil {
		return nil
	}

	sqlOp := operatorSQL(op)

	return clause.Expr{SQL: "? " + sqlOp + " ?", Vars: []any{l, r}, WithoutParentheses: false}
}

//nolint:exhaustive // no need to check every case.
func operatorSQL(op token.Type) string {
	switch op {
	case token.Eq:
		return "="
	case token.NotEq:
		return "<>"
	case token.GreaterThan:
		return ">"
	case token.GreaterThanOrEqual:
		return ">="
	case token.LessThan:
		return "<"
	case token.LessThanOrEqual:
		return "<="
	case token.And:
		return "AND"
	case token.Or:
		return "OR"
	default:
		return string(op)
	}
}
