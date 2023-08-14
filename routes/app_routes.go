package routes

import (
	"fiber-basic-auth/constants"
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

	app.Get(constants.IndexPath, func(c *fiber.Ctx) error {
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
		}, constants.LayoutsMainPath)
	})

	app.Get(constants.ProtectedPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		isSignedIn := sess.Get("is-signed-in")
		if isSignedIn == nil || isSignedIn != "true" {
			sess.Set("error", "You must be signed in to visit that page.")
			sess.Set("url-to-return-to", c.Path())
			_ = sess.Save()

			return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
		}

		// Render index within layouts/main
		return c.Render(constants.ProtectedPath, fiber.Map{
			"Title": "Protected",
		}, constants.LayoutsMainPath)
	})
}
