package middleware

import (
	"fmt"
	"strings"

	conf "thlAttractionService/pkg/config"

	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// FQDN format for kubernetes service
	authServiceURL = "http://service-auth-thl.thl-services.svc.cluster.local:8080"
	verifyEndpoint = "/api/auth/verify" // แก้ไขเป็น endpoint ที่ถูกต้อง
)

type AuthResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UserID    uint   `json:"user_id"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		SessionID string `json:"session_id"`
	} `json:"data"`
}

type AuthErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"massage"`
	Error   string `json:"error"`
}

type AuthMiddleware struct {
	config *conf.Config
}

func NewAuthMiddleware(config *conf.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (am *AuthMiddleware) Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    401,
				"message": "กรุณาระบุ token",
				"error":   "missing authorization header",
			})
		}

		// ตรวจสอบรูปแบบ token
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    401,
				"message": "รูปแบบ token ไม่ถูกต้อง",
				"error":   "invalid token format",
			})
		}

		// แยก token ออกมา
		tokenString := bearerToken[1]

		// ตรวจสอบ token ด้วย secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("รูปแบบการเข้ารหัสไม่ถูกต้อง: %v", token.Header["alg"])
			}

			// Decode base64 secret key
			decodedKey, err := base64.StdEncoding.DecodeString(am.config.App.JWT.SecretKey)

			if err != nil {
				return nil, fmt.Errorf("ไม่สามารถ decode secret key: %v", err)
			}

			// Debug log
			fmt.Printf("--------- Debug Info ---------\n")
			fmt.Printf("Algorithm: %v\n", token.Header["alg"])
			fmt.Printf("Secret Key Length: %d\n", len(decodedKey))
			fmt.Printf("Token: %s\n", tokenString[:10])
			fmt.Printf("---------------------------\n")

			parts := strings.Split(tokenString, ".")
			fmt.Printf("\n----- ข้อมูลการตรวจสอบ Token -----\n")
			fmt.Printf("1. ข้อมูล Token:\n")
			fmt.Printf("   - Algorithm: %v\n", token.Header["alg"])
			fmt.Printf("   - Token Type: %v\n", token.Header["typ"])
			fmt.Printf("\n2. ข้อมูล Secret Key:\n")
			fmt.Printf("   - ความยาว: %d ตัวอักษร\n", len(decodedKey))
			fmt.Printf("   - 10 ตัวแรก: %x\n", decodedKey[:10])
			fmt.Printf("\n3. ส่วนประกอบ Token:\n")
			fmt.Printf("   - Header: %s\n", parts[0])
			fmt.Printf("   - Payload: %s\n", parts[1])
			fmt.Printf("   - ความยาว Signature: %d\n", len(parts[2]))
			fmt.Printf("--------------------------------\n\n")

			return decodedKey, nil
		})

		if err != nil || token == nil || !token.Valid {
			decodedKey, _ := base64.StdEncoding.DecodeString(am.config.App.JWT.SecretKey)
			parts := strings.Split(tokenString, ".")

			// Default error
			errorMsg := "ไม่สามารถตรวจสอบ token ได้"
			if err != nil {
				errorMsg = fmt.Sprintf("สาเหตุ: %v", err)
			} else if token != nil && !token.Valid {
				errorMsg = "Token ไม่ valid (invalid signature หรือ expired)"
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    401,
				"message": "การตรวจสอบ Token ล้มเหลว",
				"error":   errorMsg,
				"debug": fiber.Map{
					"algorithm":      token.Header["alg"],
					"key_length":     len(am.config.App.JWT.SecretKey),
					"decoded_length": len(decodedKey),
					"decoded_key":    fmt.Sprintf("%x", decodedKey),
					"token_parts":    len(parts),
					"signing_method": "HS256",
					"header":         parts[0],
					"payload":        parts[1],
					"signature_len":  len(parts[2]),
				},
			})
		}

		// ตรวจสอบข้อมูลใน token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user_id", uint(claims["id"].(float64)))
			c.Locals("email", claims["email"].(string))
			c.Locals("full_name", claims["full_name"].(string))

			if sessionID, ok := claims["session_id"].(string); ok {
				c.Locals("session_id", sessionID)
			}

			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    401,
			"message": "Token ไม่ถูกต้อง",
			"error":   "ข้อมูลใน token ไม่ถูกต้อง",
		})
	}
}
