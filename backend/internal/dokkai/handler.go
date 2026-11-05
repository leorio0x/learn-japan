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

const maxUploadSize = 5 << 20

// Handler menangani HTTP request untuk fitur dokkai.
type Handler struct {
	svc *Service
}

// NewHandler membuat handler baru.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ListLessons godoc: GET /api/dokkai/:type/lessons
func (h *Handler) ListLessons(c *fiber.Ctx) error {
	lessons, err := h.svc.Lessons(c.Params("type"))
	if errors.Is(err, ErrInvalidKind) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(lessons)
}

// GetByLesson godoc: GET /api/dokkai/:type/:lesson
func (h *Handler) GetByLesson(c *fiber.Ctx) error {
	list, err := h.svc.ByLesson(c.Params("type"), c.Params("lesson"))
	if errors.Is(err, ErrInvalidKind) || errors.Is(err, ErrInvalidLesson) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

// BulkImport godoc: POST /api/dokkai/bulk
func (h *Handler) BulkImport(c *fiber.Ctx) error {
	kind := strings.TrimSpace(c.FormValue("type"))
	if !ValidKind(kind) {
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidKind.Error())
	}

	items, err := parseUpload(c, "file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	imported, err := h.svc.Bulk(kind, items)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"imported": imported})
}

func parseUpload(c *fiber.Ctx, field string) ([]Dokkai, error) {
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

	var items []Dokkai
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("gagal membaca json: %w", err)
	}
	return items, nil
}

func fileExtension(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i:]
	}
	return ""
}
