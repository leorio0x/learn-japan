package latihan

import "github.com/gofiber/fiber/v2"

// RegisterRoutes mendaftarkan route latihan bunpou.
func RegisterRoutes(app fiber.Router, repo *Repository, protect fiber.Handler) {
	handler := NewHandler(NewService(repo))

	group := app.Group("/latihan/bunpou")
	group.Get("/lessons", handler.ListLessons)
	group.Get("/:lesson", handler.GetByLesson)
	group.Post("/", protect, handler.Create)
	group.Put("/:lesson/:id", protect, handler.Update)
	group.Delete("/:lesson/:id", protect, handler.Delete)
}
