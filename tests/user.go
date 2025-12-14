package tests

import (
	"fmt"
)

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Age     int     `json:"age"`
	Email   *string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s %s, %d years old, email: %v", u.ID, u.Name, u.Surname, u.Age, u.Email)
}

type DefaultTest struct {
	Input       string
	ToMigrate   []*User
	ExpectedIDs []int
}

func DefaultScenarios() map[string]DefaultTest {
	email := "test@test.com"
	tdt := map[string]DefaultTest{
		"no filter": {
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{1, 2, 3, 4},
		},
		"name eq 'John'": {
			Input: "name eq 'John'",
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{1},
		},
		"name ne 'John'": {
			Input: "name ne 'John'",
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{2, 3, 4},
		},
		"not surname eq 'Doe'": {
			Input: "not surname eq 'Doe'",
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{3, 4},
		},
		"email eq null": {
			Input: "email eq null",
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{2, 3, 4},
		},
		"email ne null": {
			Input: "email ne null",
			ToMigrate: []*User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			ExpectedIDs: []int{1},
		},
	}

	return tdt
}
