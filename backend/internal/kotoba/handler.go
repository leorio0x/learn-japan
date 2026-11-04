package kotoba

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"learn-japan/backend/internal/bulk"
)

// Handler menangani HTTP request untuk fitur kotoba.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc: GET /api/kotoba?lesson=
func (h *Handler) List(c *fiber.Ctx) error {
	filter := Filter{
		Lesson: c.Query("lesson"),
	}
	list, err := h.svc.List(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderCacheControl, "no-cache")
	return c.JSON(list)
}

// Lessons godoc: GET /api/kotoba/lessons
func (h *Handler) Lessons(c *fiber.Ctx) error {
	lessons, err := h.svc.Lessons()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderCacheControl, "no-cache")
	return c.JSON(lessons)
}

// Get godoc: GET /api/kotoba/:id
func (h *Handler) Get(c *fiber.Ctx) error {
	k, err := h.svc.Get(c.Params("id"))
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(k)
}

// Create godoc: POST /api/kotoba
func (h *Handler) Create(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if msg := validate(req); msg != "" {
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	k, err := h.svc.Create(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(k)
}

// Update godoc: PUT /api/kotoba/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if msg := validate(req); msg != "" {
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	k, err := h.svc.Update(c.Params("id"), req)
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(k)
}

// BulkImport godoc: POST /api/kotoba/bulk
func (h *Handler) BulkImport(c *fiber.Ctx) error {
	rows, err := bulk.ParseUpload(c, "file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := h.svc.Bulk(rows)
	if err != nil {
		var ve *bulk.ValidationError
		if errors.As(err, &ve) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"imported": len(created),
		"data":     created,
	})
}

// Delete godoc: DELETE /api/kotoba/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	err := h.svc.Delete(c.Params("id"))
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// BulkDelete godoc: POST /api/kotoba/bulk-delete
func (h *Handler) BulkDelete(c *fiber.Ctx) error {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if len(req.IDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "minimal satu id wajib diisi")
	}
	removed, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"deleted": removed})
}

func validate(req Request) string {
	if len(cleanTerjemahan(req.Terjemahan)) == 0 {
		return "minimal satu terjemahan wajib diisi"
	}
	if strings.TrimSpace(req.Kanji) == "" && strings.TrimSpace(req.Furigana) == "" {
		return "kanji atau furigana wajib diisi"
	}
	return ""
}
