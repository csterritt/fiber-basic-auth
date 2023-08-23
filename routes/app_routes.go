package routes

import (
	"fiber-basic-auth/constants"
	"fiber-basic-auth/wrapped_session"
	"github.com/gofiber/fiber/v2"
)

func SetUpAppRoutes(app *fiber.App) {
	wrapped_session.SetupStorage()

	app.Get(constants.IndexPath, func(c *fiber.Ctx) error {
		return wrapped_session.WithSession(c, func(c *fiber.Ctx, sess *wrapped_session.WrappedSession) error {
			isSignedIn := sess.Get(constants.IsSignedInKey) == constants.IsSignedInValue

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
		return wrapped_session.WithSession(c, func(c *fiber.Ctx, sess *wrapped_session.WrappedSession) error {
			isSignedIn := sess.Get(constants.IsSignedInKey)
			if isSignedIn != constants.IsSignedInValue {
				sess.Set(constants.ErrorKey, "You must be signed in to visit that page.")
				sess.Set(constants.UrlToReturnToKey, c.Path())

				return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
			}

			errVal := getErrorIfAny(sess)
			msgVal := getMessageIfAny(sess)

			// Render index within layouts/main
			return c.Render(constants.ProtectedPath, fiber.Map{
				"Title":   "Protected",
				"Error":   errVal,
				"Message": msgVal,
			}, constants.LayoutsMainPath)
		})
	})
}
