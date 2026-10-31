package kanji

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"learn-japan/backend/internal/bulk"
)

// Handler menangani HTTP request untuk fitur kanji.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc: GET /api/kanji?lesson=&buku=
func (h *Handler) List(c *fiber.Ctx) error {
	filter := Filter{
		Lesson: c.Query("lesson"),
		Buku:   c.Query("buku"),
	}
	list, err := h.svc.List(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

// Get godoc: GET /api/kanji/:id
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

// Create godoc: POST /api/kanji
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

// Update godoc: PUT /api/kanji/:id
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

// BulkImport godoc: POST /api/kanji/bulk
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

// Delete godoc: DELETE /api/kanji/:id
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

// BulkDelete godoc: POST /api/kanji/bulk-delete
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
	if strings.TrimSpace(req.Kanji) == "" {
		return "kanji wajib diisi"
	}
	if strings.TrimSpace(req.Furigana) == "" {
		return "furigana wajib diisi"
	}
	if len(cleanTerjemahan(req.Terjemahan)) == 0 {
		return "minimal satu terjemahan wajib diisi"
	}
	return ""
}
