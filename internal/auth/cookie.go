package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CookieAccess  = "access_token"
	CookieRefresh = "refresh_token"
	CookieState   = "oauthstate"
)

func isSecure() bool { return os.Getenv("APP_ENV") == "production" }

func setAuthCookies(w http.ResponseWriter, access, refresh string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieAccess,
		Value:    access,
		Path:     "/",
		MaxAge:   int(15 * time.Minute / time.Second),
		HttpOnly: true,
		Secure:   isSecure(),
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     CookieRefresh,
		Value:    refresh,
		Path:     "/auth/refresh",
		MaxAge:   int(30 * 24 * time.Hour / time.Second),
		HttpOnly: true,
		Secure:   isSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAuthCookies(w http.ResponseWriter) {
	for name, path := range map[string]string{
		CookieAccess:  "/",
		CookieRefresh: "/auth/refresh",
	} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   isSecure(),
			SameSite: http.SameSiteLaxMode,
		})
	}
}

var stateSecret = []byte(os.Getenv("JWT_SECRET"))

func generateState() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	// Format: nonce.expiry_timestamp.hmac_signature
	payload := hex.EncodeToString(nonce) + "." +
		strconv.FormatInt(time.Now().Add(5*time.Minute).Unix(), 10)

	mac := hmac.New(sha256.New, stateSecret)
	mac.Write([]byte(payload))

	return payload + "." + hex.EncodeToString(mac.Sum(nil)), nil
}

func validateState(state string) bool {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return false
	}

	// Check expiry
	expiry, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return false
	}

	// Verify HMAC signature
	mac := hmac.New(sha256.New, stateSecret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSig), []byte(parts[2]))
}
