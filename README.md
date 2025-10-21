# GormGoQrius

[![version](https://img.shields.io/github/v/release/golaxo/gormgoqrius)](https://img.shields.io/github/v/release/golaxo/gormgoqrius)
[![PR checks](https://github.com/golaxo/gormgoqrius/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/golaxo/gormgoqrius/actions/workflows/pr-checks.yml)

> [!WARNING]
> GormGoQrius is under heavy development.

Implementation of [GoQrius][goqrius] for [GORM][gorm].

## Getting Started

GormGoQrius is using [GORM clauses](https://gorm.io/gen/clause.html) to leverage the `WHERE` clause.
So it can be used as simple as:

```go
filter := "name eq 'John'"
p := parser.New(lexer.New(filter))
e := p.Parse()

whereClause := gormgoqrius.WhereClause(e)

var users []*User
db.Clauses(whereClause).Find(&users)
```

[goqrius]: https://github.com/golaxo/goqrius
[gorm]: https://gorm.io/
