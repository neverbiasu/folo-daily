package main

import (
	"fmt"
	"net/http"

	"github.com/neverbiasu/folo-daily/handlers"
)

func main() {
	http.HandleFunc("/webhook", handlers.HandleWebhook)
	port := "8080"
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
