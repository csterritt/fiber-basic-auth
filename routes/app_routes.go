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
		return withSession(c, func(c *fiber.Ctx, sess *session.Session) error {
			isSignedIn := sess.Get(constants.IsSignedInKey) == constants.IsSignedInValue

			// any error
			errVal := getErrorIfAny(sess)
			msgVal := getMessageIfAny(sess)

			// Render index within layouts/main
			return c.Render("index", fiber.Map{
				"Title":      "Welcome!",
				"IsSignedIn": isSignedIn,
				"Error":      errVal,
				"Message":    msgVal,
			}, constants.LayoutsMainPath)
		})
	})

	app.Get(constants.ProtectedPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *session.Session) error {
			isSignedIn := sess.Get(constants.IsSignedInKey)
			if isSignedIn != constants.IsSignedInValue {
				sess.Set(constants.ErrorKey, "You must be signed in to visit that page.")
				sess.Set(constants.UrlToReturnToKey, c.Path())

				return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
			}

			// Render index within layouts/main
			return c.Render(constants.ProtectedPath, fiber.Map{
				"Title": "Protected",
			}, constants.LayoutsMainPath)
		})
	})
}
