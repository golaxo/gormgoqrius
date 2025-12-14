package tests

import "fmt"

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
