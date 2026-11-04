package main

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"learn-japan/backend/internal/auth"
	"learn-japan/backend/internal/dokkai"
	"learn-japan/backend/internal/kaiwa"
	"learn-japan/backend/internal/kanji"
	"learn-japan/backend/internal/kotoba"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	dataPath := env("DATA_PATH", "../data/kotoba.json")
	kanjiPath := env("KANJI_DATA_PATH", "../data/kanji.json")
	kaiwaPath := env("KAIWA_DATA_PATH", "../data/kaiwa.json")
	dokkaiPath := env("DOKKAI_DATA_PATH", "../data/dokkai")
	port := env("PORT", "3000")
	crudPassword := env("CRUD_PASSWORD", "armanganteng123#")

	kotobaRepo := kotoba.NewRepository(dataPath)
	kanjiRepo := kanji.NewRepository(kanjiPath)
	kaiwaRepo := kaiwa.NewRepository(kaiwaPath)
	dokkaiRepo := dokkai.NewRepository(dokkaiPath)

	authManager := auth.NewManager(crudPassword, os.Getenv("SESSION_SECRET"))

	app := fiber.New(fiber.Config{
		AppName:   "learn-japan",
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())
	app.Use(etag.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	authManager.RegisterRoutes(api)
	kotoba.RegisterRoutes(api, kotobaRepo, authManager.Middleware)
	kanji.RegisterRoutes(api, kanjiRepo, authManager.Middleware)
	kaiwa.RegisterRoutes(api, kaiwaRepo, authManager.Middleware)
	dokkai.RegisterRoutes(api, dokkaiRepo, authManager.Middleware)

	log.Printf("server berjalan di http://localhost:%s", port)
	log.Printf("data kotoba: %s | data kanji: %s | data kaiwa: %s | data dokkai: %s", dataPath, kanjiPath, kaiwaPath, dokkaiPath)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
