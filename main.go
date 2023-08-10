package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fiber-basic-auth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/jet/v2"
)

func verifyEnvironment() {
	problemFound := false
	for _, str := range []string{
		"SERVER_URL", "SERVER_PORT",
	} {
		val := os.Getenv(str)
		if len(strings.Trim(val, " \t\n\r")) == 0 {
			problemFound = true
			_, _ = fmt.Fprintf(os.Stderr, "Error: Unable to find environmental variable %s\n", str)
		}
	}

	if problemFound {
		panic("Set variables above and restart.")
	}
}

func main() {
	verifyEnvironment()

	// Create a new rendering engine
	engine := jet.New("./views", ".jet")
	engine.Verbose = os.Getenv("DEBUG") != ""

	// Pass the engine to the views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Static("/", "./public")

	routes.SetUpAuthRoutes(app)
	routes.SetUpAppRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
}
