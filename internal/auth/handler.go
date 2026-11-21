package auth

import (
	"api/internal/respond"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type handler struct {
	service AuthService
}

var googleOauthConfig *oauth2.Config

func InitOAuthConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	if googleOauthConfig.ClientID == "" ||
		googleOauthConfig.ClientSecret == "" ||
		googleOauthConfig.RedirectURL == "" {
		log.Fatal("missing GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / GOOGLE_REDIRECT_URL env vars")
	}
	log.Printf("OAuth configured with redirect: %s", googleOauthConfig.RedirectURL)
}

func NewAuthHandler(service AuthService) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		respond.Error(w, http.StatusBadRequest, "All fields are required")
		return
	}

	res, err := h.service.Register(r.Context(), &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, res)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, res)
}

func (h *handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to build state")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{
		"redirect_url": googleOauthConfig.AuthCodeURL(state),
	})
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		respond.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	claims, err := ValidateTokenClaims(req.RefreshToken)
	if err != nil || claims.Type != TokenTypeRefresh {
		respond.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	access, _ := GenerateToken(claims.UserID)
	refresh, _ := GenerateRefreshToken(claims.UserID)
	respond.JSON(w, http.StatusOK, &AuthResponse{
		Token:        access,
		RefreshToken: refresh,
	})
}

func (h *handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if !validateState(r.FormValue("state")) {
		log.Printf("[OAuth Callback] Invalid or expired state")
		http.Error(w, "Invalid or expired state", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		googleError := r.FormValue("error")
		log.Printf("[OAuth Callback] Missing authorization code. Google error: %s", googleError)
		http.Error(w, "Authorization denied or failed", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("[OAuth Callback] Token exchange failed: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[OAuth Callback] Network error fetching userinfo: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[OAuth Callback] Google API returned status %d: %s", resp.StatusCode, string(bodyBytes))
		http.Error(w, "Failed to get user info from Google", http.StatusInternalServerError)
		return
	}

	var info UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Printf("[OAuth Callback] JSON parse error: %v", err)
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	res, err := h.service.LoginWithGoogle(r.Context(), &info)
	if err != nil {
		log.Printf("[OAuth Callback] Service/DB error during login: %v", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	setAuthCookies(w, res.Token, res.RefreshToken)

	redirect := os.Getenv("FRONTEND_URL")
	if redirect == "" {
		redirect = "/"
	}

	log.Printf("[OAuth Callback] Successfully logged in user: %s", info.Email)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	log.Printf("[GetMe] Token contains User ID: %q", userID)

	user, err := h.service.GetMe(r.Context(), userID)

	log.Printf("[GetMe] Database Result -> User: %v | Error: %v", user, err)

	if err != nil || user == nil {
		respond.Error(w, http.StatusNotFound, "User not found")
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookies(w)
	respond.JSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}
