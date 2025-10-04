package cookies

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

const (
	accessTokenName  = "access_token"
	refreshTokenName = "refresh_token"
	defaultPath      = "/"
)

var ttl = 7 * 24 * time.Hour

func SetAuthTokens(w http.ResponseWriter, tokens *entity.Auth) {
	setPartitionedCookie(w, accessTokenName, tokens.AccessToken, int(ttl.Seconds()))
	setPartitionedCookie(w, refreshTokenName, tokens.RefreshToken, int(ttl.Seconds()))
}

func ClearAuthTokens(w http.ResponseWriter) {
	clearCookie(w, accessTokenName)
	clearCookie(w, refreshTokenName)
}

func GetAccessToken(r *http.Request) string {
	c, err := r.Cookie(accessTokenName)
	if err != nil {
		return ""
	}
	return c.Value
}

func GetRefreshToken(r *http.Request) string {
	c, err := r.Cookie(refreshTokenName)
	if err != nil {
		return ""
	}
	return c.Value
}

// --- helpers ---

// CHIPS: SameSite=None; Secure; Partitioned
// stdlib не имеет поля Partitioned → добавим вручную.
func setPartitionedCookie(w http.ResponseWriter, name, value string, maxAge int) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     defaultPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   maxAge,
		// Domain: "stage-backoffice.bagsy.kz", // укажи при необходимости (если кука не host-only)
	}
	// Пишем базовые атрибуты:
	http.SetCookie(w, c)
	// Добавляем ; Partitioned (Chrome/Edge понимают)
	// Важно: нужно ДОБАВИТЬ второй Set-Cookie с тем же значением + "; Partitioned"
	// чтобы не потерять атрибут. Проще — перезаписать заголовок вручную:
	w.Header().Del("Set-Cookie")
	// Соберём строку вручную:
	s := c.String() + "; Partitioned"
	w.Header().Add("Set-Cookie", s)
}

func clearCookie(w http.ResponseWriter, name string) {
	// Для корректного удаления повторяем Path/Domain/SameSite/Secure как при установке
	c := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     defaultPath,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		// Domain: "stage-backoffice.bagsy.kz",
	}
	w.Header().Add("Set-Cookie", c.String()+"; Partitioned")
}
