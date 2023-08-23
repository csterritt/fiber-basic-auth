package routes

import (
	"log"
	"math/rand"
	"os" // PRODUCTION:REMOVE
	"regexp"
	"strings"
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

func getErrorIfAny(sess *WrappedSession) interface{} {
	errVal := sess.Get(constants.ErrorKey)
	if errVal != nil && errVal != "" {
		sess.Delete(constants.ErrorKey)
	}

	return errVal
}

func getMessageIfAny(sess *WrappedSession) interface{} {
	msgVal := sess.Get(constants.MessageKey)
	if msgVal != nil && msgVal != "" {
		sess.Delete(constants.MessageKey)
	}

	return msgVal
}

func getSignInUpCode() string {
	result := digits[rand.Intn(9)+1]
	for index := 0; index < 5; index++ {
		result = result + digits[rand.Intn(10)]
	}

	return result
}

func redirectIfSignedIn(c *fiber.Ctx, sess *WrappedSession) bool {
	isSignedIn := sess.Get(constants.IsSignedInKey)
	if isSignedIn == constants.IsSignedInValue {
		_ = c.Redirect(constants.IndexPath, fiber.StatusSeeOther)

		return true
	}

	return false
}

type WrappedSession struct {
	dirty         bool
	actualSession *session.Session
}

func (ms *WrappedSession) Get(key string) interface{} {
	return ms.actualSession.Get(key)
}

func (ms *WrappedSession) Set(key string, value interface{}) {
	ms.actualSession.Set(key, value)
	ms.dirty = true
}

func (ms *WrappedSession) Delete(key string) {
	ms.actualSession.Delete(key)
	ms.dirty = true
}

func withSession(c *fiber.Ctx, wrapt func(c *fiber.Ctx, sess *WrappedSession) error) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	mySess := WrappedSession{
		dirty:         false,
		actualSession: sess,
	}
	err = wrapt(c, &mySess)

	if err == nil && mySess.dirty {
		if err = sess.Save(); err != nil {
			log.Printf("========> Error saving Session: %v\n", err)
		}
	}

	return err
}

func SetUpAuthRoutes(app *fiber.App) {
	setupStorage()

	app.Get(constants.AuthSignInPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
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
	})

	app.Post(constants.AuthSubmitSignInPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			if redirectIfSignedIn(c, sess) {
				return nil
			}

			email := c.FormValue("email")
			if isValidEmail(email) {
				sess.Set(constants.EmailKey, email)
				codeVal := getSignInUpCode()
				log.Printf("codeVal is %s\n", codeVal) // PRODUCTION:REMOVE
				sess.Set(constants.ExpectedCodeKey, codeVal)
				sess.Set(constants.SubmitTimeKey, time.Now().Unix())
				sess.Set(constants.CameFromKey, constants.AuthSignInPath)

				return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
			} else {
				sess.Set(constants.ErrorKey, "You must provide an email.")

				return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
			}
		})
	})

	app.Get(constants.AuthSignUpPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
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
	})

	app.Post(constants.AuthSubmitSignUpPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			if redirectIfSignedIn(c, sess) {
				return nil
			}

			email := c.FormValue("email")
			if isValidEmail(email) {
				sess.Set(constants.EmailKey, email)
				codeVal := getSignInUpCode()
				log.Printf("codeVal is %s\n", codeVal) // PRODUCTION:REMOVE
				sess.Set(constants.ExpectedCodeKey, codeVal)
				sess.Set(constants.SubmitTimeKey, time.Now().Unix())
				sess.Set(constants.CameFromKey, constants.AuthSignUpPath)

				return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
			} else {
				sess.Set(constants.ErrorKey, "You must provide an email.")

				return c.Redirect(constants.AuthSignUpPath, fiber.StatusSeeOther)
			}
		})
	})

	app.Get(constants.AuthEnterCodePath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			if redirectIfSignedIn(c, sess) {
				return nil
			}

			email := sess.Get(constants.EmailKey)
			errVal := getErrorIfAny(sess)
			// Render index within layouts/main
			return c.Render(constants.AuthEnterCodePath, fiber.Map{
				"Title": "Enter Code",
				"Email": email,
				"Error": errVal,
			}, constants.LayoutsMainPath)
		})
	})

	app.Post(constants.AuthSubmitCodePath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			if redirectIfSignedIn(c, sess) {
				return nil
			}

			email := sess.Get(constants.EmailKey)
			if email == "" {
				sess.Set(constants.ErrorKey, "You must provide an email address.")

				return c.Redirect(constants.AuthSignInPath, fiber.StatusSeeOther)
			}

			timeSubmitted := sess.Get(constants.SubmitTimeKey)
			isExpired := timeSubmitted == nil || time.Now().Unix()-timeSubmitted.(int64) > constants.CodeExpireTimeInSeconds
			code := c.FormValue("code")

			if code == "7654321" && os.Getenv("DEBUG") == "true" { // PRODUCTION:REMOVE
				isExpired = true // PRODUCTION:REMOVE
			} // PRODUCTION:REMOVE

			if isExpired {
				if timeSubmitted == nil {
					log.Printf("Time submitted is somehow nil?!?\n")
				} else {
					log.Printf("Code expired %v seconds ago.\n", time.Now().Unix()-timeSubmitted.(int64))
				}
				cameFrom := sess.Get(constants.CameFromKey).(string)
				sess.Delete(constants.ExpectedCodeKey)
				sess.Delete(constants.CameFromKey)
				sess.Set(constants.ErrorKey, "The code has expired, please try again.")

				return c.Redirect(cameFrom, fiber.StatusSeeOther)
			}
			expectedCode := sess.Get(constants.ExpectedCodeKey)

			if code == "1234567" && os.Getenv("DEBUG") == "true" { // PRODUCTION:REMOVE
				code = expectedCode.(string) // PRODUCTION:REMOVE
			} // PRODUCTION:REMOVE

			if code == "" {
				sess.Set(constants.ErrorKey, "You must provide the code.") // PRODUCTION: Might want to indicate where to expect the code (e.g., email, spam filter)

				return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
			} else if code != expectedCode {
				countVal := sess.Get(constants.WrongCodeEnteredCount)
				var count = 0
				if countVal != nil {
					count = countVal.(int)
				}
				count += 1
				if count > constants.WrongCodeFailureCount {
					sess.Delete(constants.WrongCodeEnteredCount)
					sess.Delete(constants.EmailKey)
					sess.Delete(constants.ExpectedCodeKey)
					sess.Delete(constants.UrlToReturnToKey)
					sess.Delete(constants.CameFromKey)
					sess.Set(constants.ErrorKey, "The wrong code was given too many times.")
					return c.Redirect(constants.IndexPath, fiber.StatusSeeOther)
				}

				sess.Set(constants.ErrorKey, "That code is incorrect.")
				sess.Set(constants.WrongCodeEnteredCount, count)

				return c.Redirect(constants.AuthEnterCodePath, fiber.StatusSeeOther)
			} else {
				pathToGoTo := sess.Get(constants.UrlToReturnToKey)
				if pathToGoTo == nil || pathToGoTo == "" {
					pathToGoTo = constants.IndexPath
				}
				cameFrom := sess.Get(constants.CameFromKey).(string)

				sess.Delete(constants.ExpectedCodeKey)
				sess.Delete(constants.UrlToReturnToKey)
				sess.Delete(constants.CameFromKey)
				sess.Set(constants.IsSignedInKey, constants.IsSignedInValue)
				if strings.Index(cameFrom, "sign-in") != -1 {
					sess.Set(constants.MessageKey, "You are signed in.")
				} else {
					sess.Set(constants.MessageKey, "You are signed up.")
				}

				return c.Redirect(pathToGoTo.(string), fiber.StatusSeeOther)
			}
		})
	})

	app.Post(constants.AuthCancelPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			if redirectIfSignedIn(c, sess) {
				return nil
			}

			sess.Delete(constants.EmailKey)
			sess.Delete(constants.ExpectedCodeKey)
			sess.Delete(constants.CameFromKey)
			sess.Delete(constants.UrlToReturnToKey)

			return c.Redirect(constants.IndexPath, fiber.StatusSeeOther)
		})
	})

	app.Post(constants.AuthSignOutPath, func(c *fiber.Ctx) error {
		return withSession(c, func(c *fiber.Ctx, sess *WrappedSession) error {
			sess.Delete(constants.IsSignedInKey)
			sess.Delete(constants.EmailKey)
			sess.Delete(constants.ExpectedCodeKey)
			sess.Delete(constants.CameFromKey)
			sess.Delete(constants.UrlToReturnToKey)

			return c.Redirect(constants.IndexPath, fiber.StatusSeeOther)
		})
	})
}
