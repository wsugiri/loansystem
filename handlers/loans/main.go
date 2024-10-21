package loans

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func getNestedValue(data map[string]interface{}, keys ...string) interface{} {
	var current interface{} = data

	for _, key := range keys {
		// Check if current is a map and contains the key
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				current = val
			} else {
				return nil // Key not found
			}
		} else {
			return nil // Not a map
		}
	}

	return current
}

func CreateLoan(c *fiber.Ctx) error {
	var body map[string]interface{}

	json.Unmarshal(c.Body(), &body)

	return c.JSON(body)
}

func ApproveLoan(c *fiber.Ctx) error {
	var body map[string]interface{}

	json.Unmarshal(c.Body(), &body)

	fmt.Println(body)
	fmt.Println(body["approval_date"])
	fmt.Println(getNestedValue(body, "properties", "id"))
	fmt.Println(getNestedValue(body, "properties", "data", "num"))
	fmt.Println(getNestedValue(body, "properties", "data", "id"))

	resp := fiber.Map{
		"params": fiber.Map{
			"id": c.Params("id"),
		},
		"body": body,
	}

	return c.JSON(resp)
}

func RejectLoan(c *fiber.Ctx) error {
	var body map[string]interface{}

	json.Unmarshal(c.Body(), &body)

	resp := fiber.Map{
		"params": fiber.Map{
			"id": c.Params("id"),
		},
		"body": body,
	}

	return c.JSON(resp)
}

func ListLoan(c *fiber.Ctx) error {
	queries := c.Queries()

	fmt.Println(queries)
	fmt.Println(queries["brand"])

	return c.JSON(queries)
}
