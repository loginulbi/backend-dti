package at

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetSecretFromHeader(r *http.Request) (secret string) {
	if r.Header.Get("secret") != "" {
		secret = r.Header.Get("secret")
	} else if r.Header.Get("Secret") != "" {
		secret = r.Header.Get("Secret")
	}
	return
}

func GetLoginFromHeader(ctx *fiber.Ctx) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is nil")
	}
	login := ctx.Get("login")
	if login != "" {
		return login, nil
	}
	login = ctx.Get("Login")
	if login != "" {
		return login, nil
	}
	return "", fmt.Errorf("login header is missing")
}

func Jsonstr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}

func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

func WriteString(respw http.ResponseWriter, statusCode int, content string) {
	respw.WriteHeader(statusCode)
	respw.Write([]byte(content))
}

func WriteJSONWithHeader(w http.ResponseWriter, r *http.Request, statusCode int, content interface{}, Origins []string) {
	origin := r.Header.Get("Origin")

	if IsAllowedOrigin(origin, Origins) {
		// Set CORS headers for the preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
		}
		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, secret")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(Jsonstr(content)))
}

func IsAllowedOrigin(origin string, Origins []string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}
