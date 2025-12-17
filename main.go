package main

import (
	"net/http"
)

func main() {
	// Serve OpenAPI spec
	http.Handle(
		"/api-docs/",
		http.StripPrefix(
			"/api-docs/",
			http.FileServer(http.Dir("./spec/api-docs")),
		),
	)

	// Serve Swagger UI (jika ada)
	http.Handle(
		"/swagger/",
		http.StripPrefix(
			"/swagger/",
			http.FileServer(http.Dir("./swagger")),
		),
	)

	http.ListenAndServe(":8000", nil)
}
