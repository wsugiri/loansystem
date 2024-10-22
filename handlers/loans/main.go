package loans

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ListLoan(c *fiber.Ctx) error {
	queries := c.Queries()

	fmt.Println(queries)
	fmt.Println(queries["brand"])

	return c.JSON(queries)
}
