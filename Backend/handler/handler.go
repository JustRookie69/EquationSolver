package handler

import (
	"context"
	"encoding/json"
	"grid-api/database"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserGridData struct {
	Message   string `json:"message"`
	GridSize  int    `json:"gridSize"`
	GridType  string `json:"gridType"`
	TimeStamp string `json:"timeStamp"`
}

// Global variables to store the database connection
var (
	dbClient *mongo.Client
	dbCtx    context.Context
)

func Initialize(client *mongo.Client, ctx context.Context) {
	dbClient = client
	dbCtx = ctx
}

// moment system /api/grid-data recieve this we call the handler.function(GridDataHandler)
// so in this case we are using post here as its data post from user
func GridDataHandler(w http.ResponseWriter, r *http.Request) {
	//so we recieve data in json format or handling post

	w.Header().Set("Access-Control-Allow-Origin", "*") // In production, use your frontend's origin
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed) // we are writing back to user updating bad request
		return
	}

	var userMessage UserGridData
	err := json.NewDecoder(r.Body).Decode(&userMessage)
	if err != nil {
		http.Error(w, "error while decoding message", http.StatusBadRequest)
	}

	//here i want send userMessage to database package GetMatrixData
	found, matrixData := database.GetMatrixData(dbClient, dbCtx, userMessage.Message)
	//Boolean, matrix
	// Create a response based on the database result
	var response map[string]interface{}
	if found {
		response = map[string]interface{}{
			"status": "success",
			"data":   matrixData,
		}
	} else {

		// response = map[string]interface{}{
		// 	"status":  "error",
		// 	"message": "Matrix not found",
		// }
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
