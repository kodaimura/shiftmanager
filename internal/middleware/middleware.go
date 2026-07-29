package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"shiftmanager/config"
	"shiftmanager/internal/core/jwt"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.AuthRefreshToken(c); err != nil {
			c.Redirect(303, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func JwtAuthApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.AuthAccessToken(c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cf := config.GetConfig()

		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != cf.BasicAuthUser || pass != cf.BasicAuthPass {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func CsrfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isSafeMethod(c.Request.Method) || isSameOrigin(c) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "不正なリクエストです。"})
		c.Abort()
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func isSameOrigin(c *gin.Context) bool {
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		return isAllowedOrigin(c, origin)
	}

	referer := c.Request.Header.Get("Referer")
	if referer != "" {
		return isAllowedOrigin(c, referer)
	}

	return false
}

func isAllowedOrigin(c *gin.Context, rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}

	return strings.EqualFold(parsed.Scheme, requestScheme(c)) &&
		strings.EqualFold(parsed.Host, requestHost(c))
}

func requestScheme(c *gin.Context) string {
	if proto := firstForwardedValue(c.Request.Header.Get("X-Forwarded-Proto")); proto != "" {
		return proto
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}

func requestHost(c *gin.Context) string {
	if host := firstForwardedValue(c.Request.Header.Get("X-Forwarded-Host")); host != "" {
		return host
	}
	return c.Request.Host
}

func firstForwardedValue(value string) string {
	if value == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(value, ",")[0])
}
