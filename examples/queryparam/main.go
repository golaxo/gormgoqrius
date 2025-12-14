package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gormgoqrius/examples"

	"github.com/golaxo/goqrius"
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

	port := 8475
	h := handler{db: db}
	http.HandleFunc("/users", h.handlerUser)
	fmt.Printf("Starting web server at localhost:%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	// e.g. curl http://localhost:8475/users?q=name%20eq%20%27John%27
}

type handler struct {
	db *gorm.DB
}

func (h handler) handlerUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get a specific parameter by name
	q := query.Get("q")
	e, err := goqrius.Parse(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	whereClause := gormgoqrius.WhereClause(e)

	var users []*examples.User
	tx := h.db
	if whereClause != nil {
		tx = tx.Where(whereClause)
	}
	tx = tx.Find(&users)
	if tx.Error != nil {
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Write the JSON response
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
