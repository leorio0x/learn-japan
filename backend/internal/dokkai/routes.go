package dokkai

import "github.com/gofiber/fiber/v2"

// RegisterRoutes mendaftarkan seluruh route fitur dokkai.
// protect adalah middleware autentikasi untuk endpoint mutasi.
func RegisterRoutes(app fiber.Router, repo *Repository, protect fiber.Handler) {
	handler := NewHandler(NewService(repo))

	group := app.Group("/dokkai")

	// Dokkai berlevel harus didaftarkan sebelum route generik :type/:lesson.
	group.Get("/level", handler.ListLevels)
	group.Get("/level/:level", handler.GetByLevel)
	group.Post("/level/bulk", protect, handler.BulkLevel)
	group.Put("/level/:level/:id", protect, handler.UpdateLevel)
	group.Delete("/level/:level/:id", protect, handler.DeleteLevel)

	group.Get("/:type/lessons", handler.ListLessons)
	group.Get("/:type/:lesson", handler.GetByLesson)
	group.Post("/bulk", protect, handler.BulkImport)
}
