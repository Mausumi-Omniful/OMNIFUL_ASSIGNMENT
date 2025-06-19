package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "IMS Running - Connected to DB host: %s", os.Getenv("DB_HOST"))
	})
	fmt.Println("Service running on :8080")
	http.ListenAndServe(":8080", nil)
}
