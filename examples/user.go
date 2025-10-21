package examples

import "fmt"

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s %s, %d years old", u.ID, u.Name, u.Surname, u.Age)
}
