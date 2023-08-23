package wrapped_session

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var storage *sqlite3.Storage // PRODUCTION: You probably want to use something else, especially on a serverless host
var store *session.Store

func SetupStorage() {
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

func WithSession(c *fiber.Ctx, wrapt func(c *fiber.Ctx, sess *WrappedSession) error) error {
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
