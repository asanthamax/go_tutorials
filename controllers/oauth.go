package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"tododemo.com/m/models"
)

func randomState() string {
	// generate a random string for oauth2 state parameter to protect against CSRF attacks
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

var (
	googleOauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = randomState()
)

func HandleGoogleLogin(w http.ResponseWriter, c *http.Request) {
	url := googleOauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, c, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Error",
			"error": "Invalid OAuth state",
		})
		return
	}

	code := c.Query("code")
	token, err := googleOauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": "Error",
			"error": "Code exchange failed",
		})
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": "Error",
			"error": "Failed getting user info",
		})
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": "Error",
			"error": "Failed reading response body",
		})
		return
	}

	var userInfo models.UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": "Error",
			"error": "Failed parsing user info",
		})
		return
	}

	// Render template with user info
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "Todo App",
		"userInfo": userInfo,
	})
}
