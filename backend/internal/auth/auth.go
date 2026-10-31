package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	cookieName = "lj_crud_token"
	tokenTTL   = 8 * time.Hour
)

// Manager menangani autentikasi password untuk fitur CRUD.
// Password hanya disimpan di sisi server sehingga tidak pernah ter-bundle
// ke JavaScript frontend ataupun terlihat lewat fitur browser.
type Manager struct {
	password string
	secret   []byte
}

// NewManager membuat Manager baru. Secret berasal dari env SESSION_SECRET;
// bila kosong akan digenerate acak sehingga token lama tidak berlaku setelah restart.
func NewManager(password, secret string) *Manager {
	key := []byte(secret)
	if len(key) == 0 {
		key = make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			panic(err)
		}
	}
	return &Manager{password: password, secret: key}
}

func (m *Manager) sign(payload string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (m *Manager) issue() string {
	expires := strconv.FormatInt(time.Now().Add(tokenTTL).Unix(), 10)
	return expires + "." + m.sign(expires)
}

func (m *Manager) valid(token string) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return false
	}
	expected := m.sign(parts[0])
	if subtle.ConstantTimeCompare([]byte(expected), []byte(parts[1])) != 1 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}
	return time.Now().Unix() < expires
}

func tokenFromRequest(c *fiber.Ctx) string {
	if v := c.Cookies(cookieName); v != "" {
		return v
	}
	header := c.Get(fiber.HeaderAuthorization)
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}

// Authenticated mengecek apakah request memiliki token yang valid.
func (m *Manager) Authenticated(c *fiber.Ctx) bool {
	return m.valid(tokenFromRequest(c))
}

// Middleware menolak request tanpa token valid.
func (m *Manager) Middleware(c *fiber.Ctx) error {
	if !m.Authenticated(c) {
		return fiber.NewError(fiber.StatusUnauthorized, "password diperlukan untuk mengakses fitur CRUD")
	}
	return c.Next()
}

// Login godoc: POST /api/auth/login
func (m *Manager) Login(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body tidak valid")
	}
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(m.password)) != 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "password salah")
	}
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    m.issue(),
		Path:     "/api",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   int(tokenTTL.Seconds()),
	})
	return c.JSON(fiber.Map{"authenticated": true})
}

// Status godoc: GET /api/auth/status
func (m *Manager) Status(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"authenticated": m.Authenticated(c)})
}

// Logout godoc: POST /api/auth/logout
func (m *Manager) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/api",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   -1,
	})
	return c.JSON(fiber.Map{"authenticated": false})
}

// RegisterRoutes mendaftarkan route autentikasi.
func (m *Manager) RegisterRoutes(app fiber.Router) {
	group := app.Group("/auth")
	group.Post("/login", m.Login)
	group.Get("/status", m.Status)
	group.Post("/logout", m.Logout)
}
