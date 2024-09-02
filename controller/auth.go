package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"login-service/config"
	"login-service/helper/at"
	"login-service/helper/atdb"
	"login-service/helper/auth"
	"login-service/helper/watoken"
	"login-service/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func AuthUser(c *fiber.Ctx) error {
	var resp model.Response
	var request struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&request); err != nil {
		resp.Response = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	creds, err := atdb.GetOneDoc[auth.GoogleCredential](config.Mongoconn, "credentials", bson.M{})
	if err != nil {
		resp.Response = "Database Connection Problem: Unable to fetch credentials"
		return c.Status(http.StatusBadGateway).JSON(resp)
	}

	payload, err := auth.VerifyIDToken(request.Token, creds.ClientID)
	if err != nil {
		resp.Response = "Invalid token: Token verification failed"
		return c.Status(http.StatusUnauthorized).JSON(resp)
	}
	userInfo := model.User{
		Nama:  payload.Claims["name"].(string),
		Email: payload.Claims["email"].(string),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.Mongoconn.Collection("user")
	filter := bson.M{"email": userInfo.Email}

	var existingUser model.User
	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil || existingUser.PhoneNumber == "" {
		response := map[string]interface{}{
			"message": "Please scan the QR code to provide your phone number",
			"user":    userInfo,
			"token":   "",
		}
		resp.Response = "Phone number not found"
		return c.Status(http.StatusUnauthorized).JSON(response)
	} else if existingUser.PhoneNumber != "" {
		token, err := watoken.EncodeforHours(existingUser.PhoneNumber, existingUser.Nama, config.PrivateKey, 18)
		if err != nil {
			resp.Response = "Token generation failed"
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}
		response := map[string]interface{}{
			"message": "Authenticated successfully",
			"user":    userInfo,
			"token":   token,
		}
		return c.Status(http.StatusOK).JSON(response)
	}

	update := bson.M{
		"$set": userInfo,
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		resp.Response = "Failed to save user info: Database update failed"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	response := map[string]interface{}{
		"user": userInfo,
	}
	return c.Status(http.StatusOK).JSON(response)
}

func GeneratePasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
		Captcha     string `json:"captcha"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	// Validate CAPTCHA
	captchaResponse, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {"0x4AAAAAAAfj2NjfaHRBhkd2VjcfmRe5gvI"},
		"response": {request.Captcha},
	})
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify captcha"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusServiceUnavailable, respn)
		return
	}
	defer captchaResponse.Body.Close()

	var captchaResult struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(captchaResponse.Body).Decode(&captchaResult); err != nil {
		var respn model.Response
		respn.Status = "Failed to decode captcha response"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}
	if !captchaResult.Success {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Invalid captcha"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Validate phone number
	re := regexp.MustCompile(`^62\d{9,15}$`)
	if !re.MatchString(request.PhoneNumber) {
		var respn model.Response
		respn.Status = "Bad Request"
		respn.Response = "Invalid phone number format"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Check if phone number exists in the 'user' collection
	userFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", userFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Phone number not registered"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Generate random password
	randomPassword, err := auth.GenerateRandomPassword(12)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to generate password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(randomPassword)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to hash password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Update or insert the user in the database
	stpFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, err = atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", stpFilter)
	var responseMessage string

	if err == mongo.ErrNoDocuments {
		// Document not found, insert new one
		newUser := model.Stp{
			PhoneNumber:  request.PhoneNumber,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
		}
		_, err = atdb.InsertOneDoc(config.Mongoconn, "stp", newUser)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to insert new user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotModified, respn)
			return
		}
		responseMessage = "New user created and password generated successfully"
	} else {
		// Document found, update the existing one
		stpUpdate := bson.M{
			"phonenumber": request.PhoneNumber,
			"password":    hashedPassword,
			"createdAt":   time.Now(),
		}
		_, err = atdb.UpdateOneDoc(config.Mongoconn, "stp", stpFilter, stpUpdate)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to update user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}
		responseMessage = "User info updated and password generated successfully"
	}

	// Respond with success and the generated password
	response := map[string]interface{}{
		"message":     responseMessage,
		"phonenumber": request.PhoneNumber,
	}
	at.WriteJSON(respw, http.StatusOK, response)

	// Send the random password via WhatsApp
	auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
}

var (
	rl = auth.NewRateLimiter(1, 5) // 1 request per second, burst of 5
)

func VerifyPasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Implementasi rate limiting
	limiter := rl.GetLimiter(request.PhoneNumber)
	if !limiter.Allow() {
		var respn model.Response
		respn.Status = "Too Many Requests"
		respn.Response = "Please try again later."
		at.WriteJSON(respw, http.StatusTooManyRequests, respn)
		return
	}

	// Find user in the database
	userFilter := bson.M{"phonenumber": request.PhoneNumber}
	user, err := atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", userFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Verify password and expiry
	if time.Now().After(user.CreatedAt.Add(4 * time.Minute)) {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Password Expired"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Find user in the 'user' collection
	myiUserFilter := bson.M{"phonenumber": request.PhoneNumber}
	existingUser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", myiUserFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Phone number not registered"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	token, err := watoken.EncodeforHours(existingUser.PhoneNumber, existingUser.Name, config.PrivateKey, 18)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to give the token"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"message": "Authenticated successfully",
		"token":   token,
		"name":    existingUser.Name,
	}

	// Respond with success
	at.WriteJSON(respw, http.StatusOK, response)
}

func ResendPasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Generate random password
	randomPassword, err := auth.GenerateRandomPassword(12)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to generate password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(randomPassword)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to hash password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Check if phone number exists in the 'stp' collection
	stpFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, stpErr := atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", stpFilter)

	if stpErr == mongo.ErrNoDocuments {
		// Document not found, insert new one
		newUser := model.Stp{
			PhoneNumber:  request.PhoneNumber,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
		}
		_, err = atdb.InsertOneDoc(config.Mongoconn, "stp", newUser)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to insert new user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}
		responseMessage := "New user created and password generated successfully"

		// Respond with success and the generated password
		response := map[string]interface{}{
			"message":     responseMessage,
			"phonenumber": request.PhoneNumber,
		}
		at.WriteJSON(respw, http.StatusOK, response)

		// Send the random password via WhatsApp
		auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
		return
	} else if stpErr != nil {
		var respn model.Response
		respn.Status = "Failed to fetch user info"
		respn.Response = stpErr.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Document found, update the existing one
	stpUpdate := bson.M{
		"phonenumber": request.PhoneNumber,
		"password":    hashedPassword,
		"createdAt":   time.Now(),
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "stp", stpFilter, stpUpdate)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to update user"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}
	responseMessage := "User info updated and password generated successfully"

	// Respond with success and the generated password
	response := map[string]interface{}{
		"message":     responseMessage,
		"phonenumber": request.PhoneNumber,
	}
	at.WriteJSON(respw, http.StatusOK, response)

	// Send the random password via WhatsApp
	auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
}
