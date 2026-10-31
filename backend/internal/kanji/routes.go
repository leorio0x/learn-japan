package kanji

import "github.com/gofiber/fiber/v2"

// RegisterRoutes mendaftarkan seluruh route fitur kanji.
// protect adalah middleware autentikasi yang dipasang pada endpoint mutasi (CRUD).
func RegisterRoutes(app fiber.Router, repo *Repository, protect fiber.Handler) {
	handler := NewHandler(NewService(repo))

	group := app.Group("/kanji")
	group.Get("/", handler.List)
	group.Get("/:id", handler.Get)
	group.Post("/bulk", protect, handler.BulkImport)
	group.Post("/bulk-delete", protect, handler.BulkDelete)
	group.Post("/", protect, handler.Create)
	group.Put("/:id", protect, handler.Update)
	group.Delete("/:id", protect, handler.Delete)
}
