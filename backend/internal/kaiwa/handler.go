package kaiwa

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk fitur kaiwa.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc: GET /api/kaiwa?lesson=&buku=
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

// Get godoc: GET /api/kaiwa/:id
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

// Create godoc: POST /api/kaiwa
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

// Update godoc: PUT /api/kaiwa/:id
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

// Delete godoc: DELETE /api/kaiwa/:id
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

// BulkDelete godoc: POST /api/kaiwa/bulk-delete
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
		if len(p.Contoh) == 0 {
			return fmt.Sprintf("pola ke-%d harus memiliki minimal satu contoh", pi+1)
		}
		for ci, c := range p.Contoh {
			if len(c.Pertanyaan) == 0 {
				return fmt.Sprintf("contoh ke-%d pada pola ke-%d harus memiliki pertanyaan", ci+1, pi+1)
			}
			if len(c.Jawaban) == 0 {
				return fmt.Sprintf("contoh ke-%d pada pola ke-%d harus memiliki jawaban", ci+1, pi+1)
			}
			for di, d := range c.Pertanyaan {
				if strings.TrimSpace(d.Kanji) == "" && strings.TrimSpace(d.Kana) == "" {
					return fmt.Sprintf("pertanyaan ke-%d pada pola ke-%d contoh ke-%d wajib diisi", di+1, pi+1, ci+1)
				}
			}
			for di, d := range c.Jawaban {
				if strings.TrimSpace(d.Kanji) == "" && strings.TrimSpace(d.Kana) == "" {
					return fmt.Sprintf("jawaban ke-%d pada pola ke-%d contoh ke-%d wajib diisi", di+1, pi+1, ci+1)
				}
			}
		}
	}
	return ""
}
