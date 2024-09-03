package at

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// func TestGetPresensiThisMonth(t *testing.T) {
// 	uri := "" //SRVLookup("mongodb+srv://xx:xxx@cxxx.xxx.mongodb.net/")
// 	print(uri)

// }

func TestGetLoginFromHeader(t *testing.T) {
	app := fiber.New()

	// Define a route to test
	app.Get("/test", func(ctx *fiber.Ctx) error {
		loginSecret, err := GetLoginFromHeader(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println(loginSecret)
		return ctx.JSON(fiber.Map{
			"login": loginSecret,
		})
	})

	// Test case 1: Missing login header
	// req := httptest.NewRequest("GET", "/test", nil)
	// resp, _ := app.Test(req)
	// assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	// Test case 2: Valid login header
	// req = httptest.NewRequest("GET", "/test", nil)
	// req.Header.Set("login", "test-secret")
	// resp, _ = app.Test(req)
	// assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Test case 3: Valid Login header (case insensitive)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Login", "test-secret")
	resp, _ := app.Test(req)
	fmt.Println(req, "||", resp)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}