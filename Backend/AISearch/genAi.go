package genAi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

// Define the structure for the matrix equation steps
type MatrixResponse struct {
	MatrixId string            `json:"matrixId" bson:"matrixId"`
	Rows     int               `json:"rows" bson:"rows"`
	Columns  int               `json:"columns" bson:"columns"`
	Cells    map[string]string `json:"cells" bson:"cells"`
}

func cleanJsonResponse(response string) string {
	// First, check if the response is wrapped in JSON code blocks
	if strings.Contains(response, "```json") && strings.Contains(response, "```") {
		startIndex := strings.Index(response, "```json") + 7
		endIndex := strings.LastIndex(response, "```")
		if startIndex > 7 && endIndex > startIndex {
			return strings.TrimSpace(response[startIndex:endIndex])
		}
	}

	// Check if the response is wrapped in generic code blocks
	if strings.Contains(response, "```") {
		startIndex := strings.Index(response, "```") + 3
		endIndex := strings.LastIndex(response, "```")
		if startIndex > 3 && endIndex > startIndex {
			return strings.TrimSpace(response[startIndex:endIndex])
		}
	}

	// If we get here, try to find the first '{' and last '}' to extract JSON
	startIndex := strings.Index(response, "{")
	endIndex := strings.LastIndex(response, "}")

	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		return strings.TrimSpace(response[startIndex : endIndex+1])
	}

	// If nothing worked, return the original response
	return response
}

func SolveEquation(equation string) (*MatrixResponse, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// For testing purposes only, not recommended for production
		apiKey = "asdlowdmasnd-asdmsakdmska-_0lasdk" // Replace with your actual API key if needed
	}

	// Initialize the Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Use the Gemini Pro model
	model := client.GenerativeModel("gemini-2.0-flash-thinking-exp-01-21")

	// Set up the system prompt to get structured JSON output for equation solving
	systemPrompt := `
You are an expert algebraic equation solver and grid formatter. Your task is to receive input and generate a structured JSON output representing the step-by-step solution of algebraic equations in a grid format. If the input is not a valid algebraic equation, you will return an empty matrix. You will also recheck your work to ensure accuracy.

**System Instructions:**

1.  **Input Validation:**
    * Receive input.
    * Determine if the input is a valid algebraic equation. A valid algebraic equation contains:
        * Variables (e.g., x, y, z).
        * Numbers (integers or decimals).
        * Mathematical operators (+, -, \*, /).
        * An equals sign (=).
        * Optional parentheses ().
    * If the input is not a valid algebraic equation, return the following JSON object:
        json
        {
          "matrixId": "invalid_input",
          "rows": 0,
          "columns": 0,
          "cells": {}
        }
        
    * If the input is a valid algebraic equation, proceed to step 2.

2.  **Equation Solving:**
    * Solve the given equation step-by-step, showing all intermediate steps.
    * Adhere to the correct order of operations (PEMDAS/BODMAS).
    * Handle fractions and decimals with precision, showing all steps of simplification and conversion.
    * Show all steps required to move variables to one side of the equals sign.
    * Show all steps required to isolate the variable.

3.  **Grid Formatting:**
    * Represent each step of the solution in a grid format within the "cells" object of the JSON output.
    * Each cell should contain a single number, variable, operator, or parenthesis.
    * Keys in the "cells" object must be in the format "rowxcolumn" (e.g., "1x1", "2x3").
    * Use "" to represent empty cells.
    * Calculate the exact number of "rows" and "columns" required to display all steps completely and accurately.
    * Ensure consistent spacing around operators and terms.
    * Maintain consistent parenthesis placement.
    * Show all steps of fraction simplification, including finding common denominators.
    * Show all steps of decimal conversion if needed.

4.  **JSON Output:**
    * Generate a JSON object with the following structure:
        json
        {
          "matrixId": "original_equation",
          "rows": number_of_rows,
          "columns": number_of_columns,
          "cells": {
            "1x1": "value",
            "1x2": "value",
            ...
            "NxM": "value"
          }
        }
        
    * Replace "original\_equation" with the input equation.
    * Provide ONLY the JSON object as output.

5.  **Recheck and Verification:**
    * After generating the JSON output, recheck the solution and grid formatting for accuracy.
    * Verify that:
        * The equation is solved correctly.
        * All steps are logically ordered and mathematically sound.
        * The grid dimensions are correct.
        * Each cell contains the appropriate value.
        * Spacing and parenthesis placement are consistent.
        * Fraction and decimal calculations are accurate.
    * If any errors are found, correct them and regenerate the JSON output.

6.  **Accuracy and Consistency:**
    * Prioritize accuracy in solving the equation.
    * Ensure all steps are presented in a logical and sequential order.
    * Maintain consistency in formatting, spacing, and parenthesis placement.
    * Break down each step into its most granular components.

7.  **Examples:**
    * For the equation "2x + 3 = 7", the output should be a JSON object with the solution steps in a grid format.
    * For the equation "(1/2)x + 3 = (2/3)x - 1", show all steps required to find common denominators and isolate x.
    * For the equation "0.25(4x + 8) - 1.5x = 3.5 - 0.75x", correctly handle all decimal calculations.
    * For the input "hello world", return the empty matrix json.

By following these instructions, you will provide accurate, consistent, and well-formatted solutions to algebraic equations in a grid-like JSON structure, handle invalid inputs appropriately, and verify the correctness of your work.
IMPORTANT: Return ONLY the raw JSON object with no markdown formatting, no code blocks, and no explanations. Do not wrap the JSON in backticks or add any additional formatting. The response should begin with "{" and end with "}" and contain only valid JSON.
`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemPrompt),
		},
	}

	// Generate response
	resp, err := model.GenerateContent(ctx, genai.Text(equation))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// Extract the response text
	responseText := string(resp.Candidates[0].Content.Parts[0].(genai.Text))

	// Clean up the response - remove markdown code blocks if present
	responseText = cleanJsonResponse(responseText)

	// Parse the JSON response
	var matrixResp MatrixResponse
	err = json.Unmarshal([]byte(responseText), &matrixResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return &matrixResp, nil
}

func main() {
	equation := "2x + 3 = 7"

	// Solve the equation
	result, err := SolveEquation(equation)
	if err != nil {
		log.Fatalf("Error solving equation: %v", err)
	}

	// Output as JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}
	fmt.Println("JSON Output:")
	fmt.Println(string(jsonData))

	// Output as BSON
	bsonData, err := bson.Marshal(result)
	if err != nil {
		log.Fatalf("Error marshaling to BSON: %v", err)
	}

	// For demonstration purposes, convert BSON back to a document to show its structure
	var bsonDoc bson.D
	err = bson.Unmarshal(bsonData, &bsonDoc)
	if err != nil {
		log.Fatalf("Error unmarshaling BSON: %v", err)
	}

	fmt.Println("\nBSON Structure (converted to document for display):")
	fmt.Printf("%+v\n", bsonDoc)

	fmt.Println("\nBSON Raw Byte Length:", len(bsonData), "bytes")
}
