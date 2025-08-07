package main

import (
	"fmt"
	"net/http"

	"github.com/adedaryorh/bookstore-app/pkg/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	routes.RegisterRoutes(r)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("Welcome to the bookStore application!")
}
