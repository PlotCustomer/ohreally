// Package klant models the Klant BusinessActor: the customer that uses the
// Webshop and performs the Selecteer producten business process.
package klant

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CookieName is the cookie that carries the Klant identity.
const CookieName = "klant_id"

// Middleware makes sure every request of a Klant carries a stable identity and
// stores it on the Echo context under CookieName.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(CookieName)
			if err != nil || cookie.Value == "" {
				value := newID()
				c.SetCookie(&http.Cookie{
					Name:     CookieName,
					Value:    value,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
				c.Set(CookieName, value)
			} else {
				c.Set(CookieName, cookie.Value)
			}
			return next(c)
		}
	}
}

// ID returns the identity of the Klant that performs the request.
func ID(c echo.Context) string {
	if value, ok := c.Get(CookieName).(string); ok {
		return value
	}
	return ""
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "klant"
	}
	return hex.EncodeToString(b)
}
