package bunpou

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk fitur bunpou.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc: GET /api/bunpou?lesson=&buku=
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

// Get godoc: GET /api/bunpou/:id
func (h *Handler) Get(c *fiber.Ctx) error {
	b, err := h.svc.Get(c.Params("id"))
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(b)
}

// Create godoc: POST /api/bunpou
func (h *Handler) Create(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if msg := validate(req); msg != "" {
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	b, err := h.svc.Create(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(b)
}

// Update godoc: PUT /api/bunpou/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if msg := validate(req); msg != "" {
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	b, err := h.svc.Update(c.Params("id"), req)
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(b)
}

// Delete godoc: DELETE /api/bunpou/:id
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

// BulkDelete godoc: POST /api/bunpou/bulk-delete
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
	if len(req.Pola) == 0 {
		return "minimal satu pola wajib diisi"
	}
	for pi, p := range req.Pola {
		if strings.TrimSpace(p.NamaPola) == "" && strings.TrimSpace(p.Rumus) == "" {
			return fmt.Sprintf("pola ke-%d wajib memiliki nama pola atau rumus", pi+1)
		}
		if len(p.Contoh) == 0 {
			return fmt.Sprintf("pola ke-%d harus memiliki minimal satu contoh", pi+1)
		}
		for ci, c := range p.Contoh {
			if c.Jawaban == nil || strings.TrimSpace(c.Jawaban.Kalimat) == "" {
				return fmt.Sprintf("contoh ke-%d pada pola ke-%d harus memiliki jawaban", ci+1, pi+1)
			}
			if c.Pertanyaan != nil && strings.TrimSpace(c.Pertanyaan.Kalimat) == "" {
				return fmt.Sprintf("pertanyaan pada contoh ke-%d pola ke-%d tidak boleh kosong", ci+1, pi+1)
			}
		}
	}
	return ""
}
