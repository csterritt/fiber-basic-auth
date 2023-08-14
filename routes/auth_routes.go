package routes

import (
	"log"
	"math/rand"
	"regexp"
	"time"

	"fiber-basic-auth/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var digits = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var storage *sqlite3.Storage // PRODUCTION: You probably want to use something else, especially on a serverless host
var store *session.Store
var emailPattern *regexp.Regexp

func init() {
	emailPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
}

// validate that the given string looks like an email address
func isValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

func setupStorage() {
	if storage == nil {
		storage = sqlite3.New() // From github.com/gofiber/storage/sqlite3
		store = session.New(session.Config{
			Storage:        storage,
			Expiration:     180 * 24 * time.Hour, // six months, approximately
			CookieSameSite: "Lax",                // PRODUCTION: Make sure you understand this before changing!
			//CookieDomain: "your-web-site-here.whatever", // PRODUCTION: THIS MUST BE SET TO YOUR WEB DOMAIN
			//CookieHTTPOnly: true,                        // PRODUCTION: THIS MUST BE SET TO true
			//CookieSecure: true,                          // PRODUCTION: THIS MUST BE SET TO true
		})
	}
}

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

func redirectIfSignedIn(c *fiber.Ctx, sess *session.Session) bool {
	isSignedIn := sess.Get("is-signed-in")
	if isSignedIn == "true" {
		_ = c.Redirect(constants.IndexPath, fiber.StatusSeeOther)

		return true
	}

	return false
}

func SetUpAuthRoutes(app *fiber.App) {
	setupStorage()

	app.Get(constants.AuthSignInPath, func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		errVal := getErrorIfAny(sess)

		// Render index within layouts/main
		return c.Render(constants.AuthSignInPath, fiber.Map{
			"Title": "Protected",
			"Error": errVal,
		}, constants.LayoutsMainPath)
	})

	app.Post(constants.AuthSubmitSignInPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		email := c.FormValue("email")
		if isValidEmail(email) {
			sess.Set("email", email)
			codeVal := getSignInUpCode()
			log.Printf("codeVal is %s\n", codeVal) // PRODUCTION: GET RID OF THIS LINE!!!
			sess.Set("expected-code", codeVal)
			_ = sess.Save()

			return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
		} else {
			sess.Set("error", "You must provide an email.")
			_ = sess.Save()

			return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
		}
	})

	app.Get(constants.AuthSignUpPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		errVal := getErrorIfAny(sess)

		// Render index within layouts/main
		return c.Render(constants.AuthSignUpPath, fiber.Map{
			"Title": "Protected",
			"Error": errVal,
		}, constants.LayoutsMainPath)
	})

	app.Post(constants.AuthSubmitSignUpPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		email := c.FormValue("email")
		if isValidEmail(email) {
			sess.Set("email", email)
			codeVal := getSignInUpCode()
			log.Printf("codeVal is %s\n", codeVal) // PRODUCTION: GET RID OF THIS LINE!!!
			sess.Set("expected-code", codeVal)
			_ = sess.Save()

			return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
		} else {
			sess.Set("error", "You must provide an email.")
			_ = sess.Save()

			return c.Redirect(constants.AuthSignUpPath, fiber.StatusSeeOther)
		}
	})

	app.Get(constants.AuthEnterCodePath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		email := sess.Get("email")
		errVal := getErrorIfAny(sess)
		// Render index within layouts/main
		return c.Render(constants.AuthEnterCodePath, fiber.Map{
			"Title": "Enter Code",
			"Email": email,
			"Error": errVal,
		}, constants.LayoutsMainPath)
	})

	app.Post(constants.AuthSubmitCodePath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		email := sess.Get("email")
		if email == "" {
			sess.Set("error", "You must provide an email address.")
			_ = sess.Save()

			return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
		}

		code := c.FormValue("code")
		expectedCode := sess.Get("expected-code")
		if code == "" {
			sess.Set("error", "You must provide the code.") // PRODUCTION: Might want to indicate where it is
			_ = sess.Save()

			return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
		} else if code != expectedCode {
			sess.Set("error", "That code is incorrect or expired.")
			_ = sess.Save()

			return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
		} else {
			pathToGoTo := sess.Get("url-to-return-to")
			if pathToGoTo == nil || pathToGoTo == "" {
				pathToGoTo = constants.IndexPath
			}
			sess.Set("code", "")
			sess.Set("is-signed-in", "true")
			sess.Set("url-to-return-to", "")
			_ = sess.Save()

			return c.Redirect(pathToGoTo.(string), fiber.StatusSeeOther)
		}
	})

	app.Post(constants.AuthCancelPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if redirectIfSignedIn(c, sess) {
			return nil
		}

		sess.Set("email", "")
		sess.Set("code", "")
		sess.Set("url-to-return-to", "")
		_ = sess.Save()

		return c.Redirect(constants.IndexPath, fiber.StatusSeeOther)
	})

	app.Post(constants.AuthSignOutPath, func(c *fiber.Ctx) error {
		// Get session from storage
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Set("email", "")
		sess.Set("code", "")
		sess.Set("is-signed-in", "")
		sess.Set("url-to-return-to", "")
		_ = sess.Save()

		return c.Redirect(constants.IndexPath, fiber.StatusSeeOther)
	})
}
