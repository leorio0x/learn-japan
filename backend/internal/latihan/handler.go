package latihan

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk latihan bunpou.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ListLessons godoc: GET /api/latihan/bunpou/lessons
func (h *Handler) ListLessons(c *fiber.Ctx) error {
	list, err := h.svc.ListLessons()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

// GetByLesson godoc: GET /api/latihan/bunpou/:lesson
func (h *Handler) GetByLesson(c *fiber.Ctx) error {
	list, err := h.svc.GetByLesson(c.Params("lesson"))
	if errors.Is(err, ErrInvalidLesson) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

// Create godoc: POST /api/latihan/bunpou
func (h *Handler) Create(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	soal, err := h.svc.Create(req)
	if errors.Is(err, ErrInvalidLesson) || errors.Is(err, ErrInvalidRequest) || errors.Is(err, ErrInvalidKategori) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(soal)
}

// Update godoc: PUT /api/latihan/bunpou/:lesson/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	soal, err := h.svc.Update(c.Params("lesson"), c.Params("id"), req)
	if errors.Is(err, ErrInvalidLesson) || errors.Is(err, ErrInvalidRequest) || errors.Is(err, ErrInvalidKategori) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(soal)
}

// Delete godoc: DELETE /api/latihan/bunpou/:lesson/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	err := h.svc.Delete(c.Params("lesson"), c.Params("id"))
	if errors.Is(err, ErrInvalidLesson) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
