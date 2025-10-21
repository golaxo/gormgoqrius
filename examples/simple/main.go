package main

import (
	"fmt"

	"gormgoqrius/examples"

	"github.com/golaxo/goqrius/lexer"
	"github.com/golaxo/goqrius/parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/golaxo/gormgoqrius"
)

func main() {
	db, err := gorm.Open(sqlite.Open("file:mem?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&examples.User{})

	migrateUsers := []*examples.User{
		{ID: 1, Name: "John", Surname: "Doe", Age: 20},
		{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
		{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
		{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
	}
	db.CreateInBatches(&migrateUsers, len(migrateUsers))
	fmt.Printf("%d users created\n", len(migrateUsers))

	filters := []string{
		"name eq 'John'",
		"name ne 'John'",
		"not name eq 'John'",
		"age gt 18 and age lt 65",
		"age le 18 or age gt 65",
	}
	for _, filter := range filters {
		fmt.Printf("filter: %s\n", filter)
		p := parser.New(lexer.New(filter))
		e := p.Parse()

		whereClause := gormgoqrius.WhereClause(e)

		var users []*examples.User
		db.Clauses(whereClause).Find(&users)
		fmt.Printf("%d user(s) found\n", len(users))
		for _, u := range users {
			fmt.Printf("\t%+v\n", *u)
		}
	}
}
