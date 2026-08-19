package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/pets", petsHandler)
	http.HandleFunc("/appointments", appointmentsHandler)
	http.HandleFunc("/appointments/session", sessionHandler)

	fmt.Println("Servidor corriendo en http://localhost:8080")

	log.Fatal(
		http.ListenAndServe(
			":8080",
			nil,
		),
	)
}
