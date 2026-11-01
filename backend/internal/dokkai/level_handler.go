package dokkai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ListLevels godoc: GET /api/dokkai/level
func (h *Handler) ListLevels(c *fiber.Ctx) error {
	levels, err := h.svc.Levels()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(levels)
}

// GetByLevel godoc: GET /api/dokkai/level/:level
func (h *Handler) GetByLevel(c *fiber.Ctx) error {
	list, err := h.svc.ByLevel(c.Params("level"))
	if errors.Is(err, ErrInvalidLevel) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

// BulkLevel godoc: POST /api/dokkai/level/bulk
func (h *Handler) BulkLevel(c *fiber.Ctx) error {
	defaultLevel := c.FormValue("level")
	items, err := parseLevelUpload(c, "file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	imported, err := h.svc.BulkLevel(defaultLevel, items)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"imported": imported})
}

// UpdateLevel godoc: PUT /api/dokkai/level/:level/:id
func (h *Handler) UpdateLevel(c *fiber.Ctx) error {
	var item LevelDokkai
	if err := json.Unmarshal(c.Body(), &item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("gagal membaca json: %v", err))
	}

	story, err := h.svc.UpdateLevel(c.Params("level"), c.Params("id"), item)
	if errors.Is(err, ErrInvalidLevel) || errors.Is(err, ErrInvalidStoryID) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.Is(err, ErrStoryNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(story)
}

// DeleteLevel godoc: DELETE /api/dokkai/level/:level/:id
func (h *Handler) DeleteLevel(c *fiber.Ctx) error {
	err := h.svc.DeleteLevel(c.Params("level"), c.Params("id"))
	if errors.Is(err, ErrInvalidLevel) || errors.Is(err, ErrInvalidStoryID) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.Is(err, ErrStoryNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseLevelUpload(c *fiber.Ctx, field string) ([]LevelDokkai, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return nil, fmt.Errorf(`file wajib diunggah pada field %q`, field)
	}
	if strings.ToLower(fileExtension(file.Filename)) != ".json" {
		return nil, errors.New("format file dokkai harus .json")
	}
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxUploadSize+1))
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file: %w", err)
	}
	if len(data) > maxUploadSize {
		return nil, fmt.Errorf("ukuran file melebihi %d MB", maxUploadSize>>20)
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	var items []LevelDokkai
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("gagal membaca json: %w", err)
	}
	return items, nil
}
