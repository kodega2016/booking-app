package main

import (
	"fmt"
	"net/http"

	"booking-app/pkg/handlers"
)

const port = 8080

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("server is running on port %d\n", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
