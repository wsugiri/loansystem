package master

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/utils"
)

type User struct {
	ID    int
	Name  string
	Email string
	Role  string
}

func ListAllUser(c *fiber.Ctx) error {
	role := c.Params("role")

	var (
		rows *sql.Rows
		err  error
	)

	if role != "" {
		// Query with a condition if the role parameter is provided
		rows, err = utils.DB.Query("SELECT id, name, email, role FROM users WHERE role = ?", role)
	} else {
		// Query without conditions if the role parameter is not provided
		rows, err = utils.DB.Query("SELECT id, name, email, role FROM users")
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	resp := fiber.Map{
		"status": "success",
		"total":  len(users),
		"items":  users,
	}

	return c.JSON(resp)
}
