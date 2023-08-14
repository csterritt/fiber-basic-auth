package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fiber-basic-auth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(helmet.New())

	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		//AllowOrigins: "https://the-final-website.com", // PRODUCTION: Comment in, set this to your site(s)
		AllowOrigins: "http://localhost:3000", // PRODUCTION: REMOVE THIS LINE or comment it out
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_KEY"),
	}))

	app.Use(etag.New(etag.Config{
		Weak: true, // speed is important to me, but caching may be more important to you. See the docs.
	}))

	// Not sure how to use this. I can sign in via a POST, but if I wait the Expiration time, then I can't
	// sign out because that's a POST!
	//app.Use(csrf.New(csrf.Config{
	//	KeyLookup:  "cookie:csrf_",
	//	CookieName: "csrf_",
	//	//CookieDomain: "your-web-site-here.whatever", // PRODUCTION: THIS MUST BE SET TO YOUR WEB DOMAIN
	//	CookieDomain:   "localhost",
	//	CookieSameSite: "Lax", // PRODUCTION: Make sure you understand this before changing!
	//	//CookieHTTPOnly: true,                        // PRODUCTION: THIS MUST BE SET TO true
	//	//CookieSecure: true,                          // PRODUCTION: THIS MUST BE SET TO true
	//	Expiration: 1 * time.Hour,                     // Hmmmmmm... see above
	//}))

	routes.SetUpAuthRoutes(app)
	routes.SetUpAppRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
}
