package main

import (
	"fmt"
	"grid-api/database"
	"grid-api/handler"
	"log"
	"net/http"
)

func main() {
	fmt.Println("starting server main")
	fmt.Println("connecting db")

	client, ctx, cancel := database.DBConnect()
	defer cancel()               // Cancel the context when done
	defer client.Disconnect(ctx) // Properly disconnect from MongoDB
	if !database.TestConnection(client, ctx) {
		log.Fatal("Could not connect to MongoDB. Exiting...")
	}
	// Initialize the handler with the database connection
	handler.Initialize(client, ctx)

	http.HandleFunc("/api/grid-data", handler.GridDataHandler)

	log.Println("initiating Server on  http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
