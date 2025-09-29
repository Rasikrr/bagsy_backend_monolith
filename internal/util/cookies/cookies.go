package cookies

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
)

const (
	accessTokenName  = "access_token"
	refreshTokenName = "refresh_token"
	defaultPath      = "/"
)

var (
	ttl = time.Hour * 24 * 7
)

func SetAuthTokens(w http.ResponseWriter, tokens *entity.Auth) {
	SetAccessToken(w, tokens.AccessToken)
	SetRefreshToken(w, tokens.RefreshToken)
}

func SetAccessToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenName,
		Value:    token,
		Path:     defaultPath,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		MaxAge:   int(ttl.Seconds()),
	})
}

func SetRefreshToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenName,
		Value:    token,
		Path:     defaultPath,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		MaxAge:   int(ttl.Seconds()),
	})
}

func ClearAuthTokens(w http.ResponseWriter) {
	ClearAccessToken(w)
	ClearRefreshToken(w)
}

func ClearAccessToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenName,
		Value:    "",
		Path:     defaultPath,
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func ClearRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenName,
		Value:    "",
		Path:     defaultPath,
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func GetAccessToken(r *http.Request) string {
	cookie, err := r.Cookie(accessTokenName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GetRefreshToken(r *http.Request) string {
	cookie, err := r.Cookie(refreshTokenName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
