package routes

import (
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var digits = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var storage *sqlite3.Storage
var store *session.Store

func getErrorIfAny(sess *session.Session) interface{} {
	errVal := sess.Get("error")
	if errVal != nil && errVal != "" {
		sess.Set("error", nil)
		_ = sess.Save()
	}

	return errVal
}

func getSignInUpCode() string {
	result := digits[rand.Intn(9)+1]
	for index := 0; index < 5; index++ {
		result = result + digits[rand.Intn(10)]
	}

	return result
}

func SetUpAuthRoutes(app *fiber.App) {
	// set up local sqlite3 storage for session information
	if storage == nil {
		storage = sqlite3.New() // From github.com/gofiber/storage/sqlite3
		store = session.New(session.Config{
			Storage: storage,
		})
	}

	app.Get("/auth/sign-in", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		errVal := getErrorIfAny(sess)

		// Render index within layouts/main
		return c.Render("auth/sign-in", fiber.Map{
			"Title": "Protected",
			"Error": errVal,
		}, "layouts/main")
	})

	app.Post("/auth/submit-sign-in", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		email := c.FormValue("email")
		log.Printf("Entered submit-sign-in, email is '%s'\n", email)
		if email == "" {
			sess.Set("error", "You must provide an email.")
			_ = sess.Save()

			return c.Redirect("/auth/sign-in", fiber.StatusSeeOther)
		} else {
			sess.Set("email", email)
			codeVal := getSignInUpCode()
			log.Printf("codeVal is %s\n", codeVal)
			sess.Set("expected-code", codeVal)
			_ = sess.Save()

			return c.Redirect("/auth/enter-code", fiber.StatusSeeOther)
		}
	})

	app.Get("/auth/sign-up", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		errVal := getErrorIfAny(sess)

		// Render index within layouts/main
		return c.Render("auth/sign-up", fiber.Map{
			"Title": "Protected",
			"Error": errVal,
		}, "layouts/main")
	})

	app.Post("/auth/submit-sign-up", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		email := c.FormValue("email")
		log.Printf("Entered submit-sign-up, email is '%s'\n", email)
		if email == "" {
			sess.Set("error", "You must provide an email.")
			_ = sess.Save()

			return c.Redirect("/auth/sign-up", fiber.StatusSeeOther)
		} else {
			sess.Set("email", email)
			codeVal := getSignInUpCode()
			log.Printf("codeVal is %s\n", codeVal)
			sess.Set("expected-code", codeVal)
			_ = sess.Save()

			return c.Redirect("/auth/enter-code", fiber.StatusSeeOther)
		}
	})

	app.Get("/auth/enter-code", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		email := sess.Get("email")
		errVal := getErrorIfAny(sess)
		log.Printf("Entered enter-code, email is '%s'\n", email)
		// Render index within layouts/main
		return c.Render("auth/enter-code", fiber.Map{
			"Title": "Enter Code",
			"Email": email,
			"Error": errVal,
		}, "layouts/main")
	})

	app.Post("/auth/submit-code", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		email := sess.Get("email")
		if email == "" {
			sess.Set("error", "You must provide an email address.")
			_ = sess.Save()

			return c.Redirect("/auth/sign-in", fiber.StatusSeeOther)
		}

		code := c.FormValue("code")
		expectedCode := sess.Get("expected-code")
		log.Printf("Entered submit-code, code is '%s', expected is '%s'\n", code, expectedCode)
		if code == "" {
			sess.Set("error", "You must provide an code.")
			_ = sess.Save()

			return c.Redirect("/auth/enter-code", fiber.StatusSeeOther)
		} else if code != expectedCode {
			sess.Set("error", "That code is incorrect or expired.")
			_ = sess.Save()

			return c.Redirect("/auth/enter-code", fiber.StatusSeeOther)
		} else {
			sess.Set("code", "")
			sess.Set("is-signed-in", "true")
			_ = sess.Save()

			return c.Redirect("/", fiber.StatusSeeOther)
		}
	})

	app.Post("/auth/cancel-sign-in", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Set("email", "")
		sess.Set("code", "")
		_ = sess.Save()

		return c.Redirect("/", fiber.StatusSeeOther)
	})

	app.Post("/auth/sign-out", func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Set("email", "")
		sess.Set("code", "")
		sess.Set("is-signed-in", "")
		_ = sess.Save()

		return c.Redirect("/", fiber.StatusSeeOther)
	})
}
