package database

import (
	"context"
	genAi "grid-api/AISearch"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI = "........................................................................tls=true" //add your own mogoURI

// defining database and collections
const mongoDB = "YourDBName"
const dbCollection = "YourCollectionName"

func DBConnect() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second) //adjust time as needed make sure continue request before code slips to cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx, cancel
}

// fetch matrix information by name
// fetch matrix information by name
func GetMatrixData(client *mongo.Client, ctx context.Context, name string) (bool, map[string]interface{}) {
	collection := client.Database(mongoDB).Collection(dbCollection)

	filter := bson.M{"matrixId": name}

	var matrixData bson.M

	err := collection.FindOne(ctx, filter).Decode(&matrixData)
	if err == mongo.ErrNoDocuments {
		// Matrix not found, generate it and save to DB
		return SaveEquationResult(client, ctx, name)
	} else if err != nil {
		log.Println("err fetching", err)
		return false, nil
	}
	return true, formatMatrixData(matrixData)
}

func SaveEquationResult(client *mongo.Client, ctx context.Context, equation string) (bool, map[string]interface{}) {
	// Call the AI to solve the equation
	matrixResp, err := genAi.SolveEquation(equation)
	if err != nil {
		log.Printf("Error solving equation: %v", err)
		return false, nil
	}

	// Convert MatrixResponse to a map for MongoDB
	matrixData := map[string]interface{}{
		"matrixId": matrixResp.MatrixId,
		"rows":     matrixResp.Rows,
		"columns":  matrixResp.Columns,
		"cells":    matrixResp.Cells,
	}

	// Insert the data into MongoDB
	collection := client.Database(mongoDB).Collection(dbCollection)
	_, err = collection.InsertOne(ctx, matrixData)
	if err != nil {
		log.Printf("Failed to insert matrix data: %v", err)
		return false, nil
	}

	// Return true and the formatted data
	return true, matrixData
}
func formatMatrixData(data bson.M) map[string]interface{} {
	formatted := map[string]interface{}{
		"matrixId": data["matrixId"],
		"rows":     data["rows"],
		"columns":  data["columns"],
		"cells":    data["cells"], // This keeps the cell data structure intact
	}
	return formatted
}

func TestConnection(client *mongo.Client, ctx context.Context) bool {
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return false
	}
	log.Println("Successfully connected to MongoDB")
	return true
}
//test
// func main() {
// 	client, ctx, cancel := dbConnect()
// 	defer cancel()
// 	defer client.Disconnect(ctx)
// 	err := client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal("Failed to connect:", err)
// 	}
// 	fmt.Println("Connected to MongoDB")
// 	// Example test case
// 	matrixName := "A"
// 	found, matrixData := getMatrixData(client, ctx, matrixName)
// 	if found {
// 		fmt.Println("Matrix found:", matrixData)
// 	} else {
// 		fmt.Println("Matrix not found")
// 	}
// }
