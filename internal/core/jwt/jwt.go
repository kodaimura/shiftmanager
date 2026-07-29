package jwt

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpackage "github.com/golang-jwt/jwt/v4"

	"shiftmanager/config"
)

func SetRefreshTokenToCookie(c *gin.Context, pl Payload) error {
	jwtStr, err := EncodeJwt(pl)
	if err != nil {
		return err
	}
	cf := config.GetConfig()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(COOKIE_KEY_REFRESH_TOKEN, jwtStr, int(REFRESH_TOKEN_EXPIRES), "/", cookieDomain(cf.AppHost), isSecureCookie(c), true)
	return nil
}

func RemoveRefreshTokenFromCookie(c *gin.Context) {
	cf := config.GetConfig()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(COOKIE_KEY_REFRESH_TOKEN, "", -1, "/", cookieDomain(cf.AppHost), isSecureCookie(c), true)
}

func GetPayload(c *gin.Context) Payload {
	pl := c.Keys[CONTEXT_KEY_PAYLOAD]
	if pl == nil {
		return Payload{}
	}
	return pl.(Payload)
}

func EncodeJwt(pl Payload) (string, error) {
	return encodeJwt(pl)
}

func ExpireJwt(pl Payload) Payload {
	pl.IssuedAt = time.Now().Unix()
	pl.ExpiresAt = time.Now().Unix()
	return pl
}

func AuthAccessToken(c *gin.Context) error {
	tokenStr, err := getAccessToken(c)
	if err != nil {
		return err
	}

	pl, err := decodeJwt(tokenStr, TOKEN_TYPE_ACCESS)
	if err != nil {
		return err
	}

	c.Set(CONTEXT_KEY_PAYLOAD, pl)
	return nil
}

func AuthRefreshToken(c *gin.Context) error {
	tokenStr, err := getRefreshToken(c)
	if err != nil {
		return err
	}

	pl, err := decodeJwt(tokenStr, TOKEN_TYPE_REFRESH)
	if err != nil {
		return err
	}

	c.Set(CONTEXT_KEY_PAYLOAD, pl)
	return nil
}

func encodeJwt(pl Payload) (string, error) {
	cf := config.GetConfig()
	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, pl)
	return token.SignedString([]byte(cf.JwtSecretKey))
}

func decodeJwt(encoded string, tokenType string) (Payload, error) {
	cf := config.GetConfig()
	token, err := jwtpackage.Parse(encoded, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(cf.JwtSecretKey), nil
	})
	if err != nil {
		return Payload{}, err
	}
	if !token.Valid {
		return Payload{}, errors.New("Invalid token")
	}

	pl, err := convertToPayload(token)
	if err != nil {
		return Payload{}, err
	}
	if pl.TokenType != tokenType {
		return Payload{}, errors.New("Invalid token type")
	}
	return pl, nil
}

func getAccessToken(c *gin.Context) (string, error) {
	token := c.Request.Header.Get(HEADER_KEY_ACCESS_TOKEN)
	if token != "" {
		return token, nil
	}
	return "", errors.New("Access token not found")
}

func getRefreshToken(c *gin.Context) (string, error) {
	token, err := c.Cookie(COOKIE_KEY_REFRESH_TOKEN)
	if err == nil {
		return token, nil
	}
	return "", errors.New("Refresh token not found")
}

func convertToPayload(token *jwtpackage.Token) (Payload, error) {
	var pl Payload

	jsonString, err := json.Marshal(token.Claims.(jwtpackage.MapClaims))

	if err == nil {
		err = json.Unmarshal(jsonString, &pl)
	}

	return pl, err
}

func cookieDomain(host string) string {
	if host == "" || host == "localhost" || host == "127.0.0.1" {
		return ""
	}
	return host
}

func isSecureCookie(c *gin.Context) bool {
	if c.Request.TLS != nil || strings.EqualFold(c.Request.Header.Get("X-Forwarded-Proto"), "https") {
		return true
	}
	host := config.GetConfig().AppHost
	return host != "" && host != "localhost" && host != "127.0.0.1"
}
