package dokkai

import "github.com/gofiber/fiber/v2"

// RegisterRoutes mendaftarkan seluruh route fitur dokkai.
// protect adalah middleware autentikasi untuk endpoint mutasi.
func RegisterRoutes(app fiber.Router, repo *Repository, protect fiber.Handler) {
	handler := NewHandler(NewService(repo))

	group := app.Group("/dokkai")
	group.Get("/:type/:lesson", handler.GetByLesson)
	group.Post("/bulk", protect, handler.BulkImport)
}
