package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kaka/kodeakademia/be/internal/domain"
	"github.com/kaka/kodeakademia/be/internal/domain/repository"
	repoimpl "github.com/kaka/kodeakademia/be/internal/repository"
	"github.com/kaka/kodeakademia/be/internal/usecase"
)

func main() {
	_ = godotenv.Load()

	// --- Database ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := repoimpl.NewPostgresUserRepository(db)

	// --- Echo ---
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/login", loginPageHandler)

	e.GET("/auth/google/login", googleLoginHandler)
	e.POST("/auth/google/login", googleLoginHandler)

	e.GET("/auth/google/callback", func(c echo.Context) error {
		return googleCallbackHandlerWithRepo(c, userRepo)
	})

	e.GET("/me", meHandler, jwtMiddleware)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("🚀 server running on %s", addr)
	log.Fatal(e.Start(addr))
}

// ----------------------------------------------------------------------
// Handlers
// ----------------------------------------------------------------------

func loginPageHandler(c echo.Context) error {
	html := `<!doctype html>
<html>
<head><meta charset="utf-8"><title>Login</title></head>
<body>
<p>Redirecting to Google login…</p>
<script>window.location.href='/auth/google/login';</script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

// Step 1: Redirect user to Google
func googleLoginHandler(c echo.Context) error {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")
	scopes := os.Getenv("OAUTH_SCOPES")

	if clientID == "" || redirectURI == "" || scopes == "" {
		return c.String(http.StatusInternalServerError, "oauth config missing")
	}

	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("redirect_uri", redirectURI)
	v.Set("response_type", "code")
	v.Set("scope", scopes)
	v.Set("access_type", "offline")
	v.Set("include_granted_scopes", "true")

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Step 2: Google callback
func googleCallbackHandlerWithRepo(
	c echo.Context,
	userRepo repository.UserRepository,
) error {

	ctx := c.Request().Context()
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing code",
		})
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")

	tokenEndpoint := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenEndpoint,
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "token request failed"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "token exchange failed"})
	}

	var tokResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	_ = json.Unmarshal(body, &tokResp)

	// --- Fetch user info ---
	req2, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/oauth2/v3/userinfo",
		nil,
	)
	req2.Header.Set("Authorization", "Bearer "+tokResp.AccessToken)

	resp2, err := client.Do(req2)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "userinfo failed"})
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)

	var profile map[string]interface{}
	_ = json.Unmarshal(body2, &profile)

	profileStr := map[string]string{}
	for k, v := range profile {
		if s, ok := v.(string); ok {
			profileStr[k] = s
		}
	}

	user, _, err := usecase.LoginOrProvisionUser(ctx, userRepo, profileStr)
	if err != nil {
         log.Printf("user provision error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user provision failed"})
	}

	jwtToken, err := generateJWTForUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "jwt failed"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": jwtToken,
		"user":  user,
	})
}

// ----------------------------------------------------------------------
// JWT
// ----------------------------------------------------------------------

func generateJWTForUser(user *domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		c.Set("user", token.Claims)
		return next(c)
	}
}

func meHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("user"))
}
