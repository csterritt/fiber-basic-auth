package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

func SetUpAppRoutes(app *fiber.App) {
	// set up local sqlite3 storage for session information
	if storage == nil {
		storage = sqlite3.New() // From github.com/gofiber/storage/sqlite3
		store = session.New(session.Config{
			Storage: storage,
		})
	}

	app.Get("/", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		name := sess.Get("name")
		isSignedIn := sess.Get("is-signed-in")

		// any error
		errVal := getErrorIfAny(sess)

		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title":      "Welcome!",
			"Name":       name,
			"IsSignedIn": isSignedIn,
			"Error":      errVal,
		}, "layouts/main")
	})

	app.Get("/protected", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		isSignedIn := sess.Get("is-signed-in")
		log.Printf("/protected route found isSignedIn '%v'\n", isSignedIn)
		if isSignedIn == nil || isSignedIn != "true" {
			sess.Set("error", "You must be signed in to visit that page.")
			_ = sess.Save()

			return c.Redirect("/auth/sign-in", fiber.StatusSeeOther)
		} else {
			// Render index within layouts/main
			return c.Render("protected", fiber.Map{
				"Title": "Protected",
			}, "layouts/main")
		}
	})
}
