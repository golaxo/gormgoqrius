module gormgoqrius/examples

go 1.25.0

replace github.com/golaxo/gormgoqrius => ../

require (
	github.com/golaxo/goqrius v0.0.5
	github.com/golaxo/gormgoqrius v0.0.3
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.2
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.50 // indirect
	golang.org/x/text v0.41.0 // indirect
)
