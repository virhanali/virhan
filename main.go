package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID    uuid.UUID `db:"id" json:"id"`
	Name  string    `db:"name" json:"name"`
	Email string    `db:"email" json:"email"`
}

var db *sqlx.DB

func initDB() {
	var err error
	connStr := "postgres://postgres:Hanhan123@db:5432/virhan_db?sslmode=disable"
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
}

func getUsers(c *fiber.Ctx) error {
	users := []User{}
	err := db.Select(&users, "SELECT id, name, email FROM users")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(users)
}

func createUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	user.ID = uuid.New()
	_, err := db.NamedExec(`INSERT INTO users (id, name, email) VALUES (:id, :name, :email)`, user)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(user)
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello, World!",
		})
	})

	users := api.Group("/users")
	users.Get("/", getUsers)
	users.Post("/", createUser)
}

func main() {
	initDB()
	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
