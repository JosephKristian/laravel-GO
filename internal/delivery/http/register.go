package http

import (
	"log"
	"strings"

	"github.com/JosephKristian/project-migration/internal/models"
	"github.com/JosephKristian/project-migration/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RegisterController struct {
	registerService          *service.RegisterService
	accountActivationService *service.AccountActivationService
	validate                 *validator.Validate
}

func NewRegisterController(
	registerService *service.RegisterService,
	accountActivationService *service.AccountActivationService,
) *RegisterController {
	return &RegisterController{
		registerService:          registerService,
		accountActivationService: accountActivationService,
		validate:                 validator.New(),
	}
}

func (r *RegisterController) Register(c *gin.Context) {
	var userInput models.RegisterInput

	// Ambil Bearer Token dari header Authorization
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"status":    "error",
			"data":      nil,
			"message":   "Authorization token is required",
			"errorCode": 400,
		})
		return
	}

	// Hapus prefix "Bearer " pada token (jika ada)
	token = strings.TrimPrefix(token, "Bearer ")

	// Log token untuk debugging (hati-hati saat log token di production)
	log.Printf("Received Bearer token: %s", token)

	// Log menerima request
	log.Printf("Received registration request: %+v", c.Request)

	// Bind input multipart-form data ke struct
	if err := c.ShouldBind(&userInput); err != nil {
		log.Printf("Error binding request data: %v", err) // Log error jika binding gagal
		c.JSON(422, gin.H{
			"status":    "error",
			"data":      nil,
			"message":   "Data tidak valid",
			"errorCode": 422,
			"errors":    err.Error(),
		})
		return
	}

	// Log validasi input
	log.Printf("Validating input for registration: %+v", userInput)

	// Panggil service untuk registrasi user
	registeredUser, err := r.registerService.Register(&userInput) // Kirim token ke service
	if err != nil {
		log.Printf("Registration failed: %v", err) // Log error jika registrasi gagal
		c.JSON(500, gin.H{
			"status":    "error",
			"data":      nil,
			"message":   "Registration failed",
			"errorCode": 500,
			"errors":    err.Error(),
		})
		return
	}

	// Log sukses registrasi
	log.Printf("User registered successfully: %+v", registeredUser)

	// Respons sukses setelah berhasil registrasi
	c.Header("Location", "/db/v1/auth/register/"+registeredUser.Email)
	c.JSON(201, gin.H{
		"status":  "success",
		"data":    userInput,
		"message": "User registered successfully",
	})
}

// Tambahkan handler untuk Account Activation
func (r *RegisterController) AccountActivation(c *gin.Context) {
	var req struct {
		EmailOrPhone     string `json:"email_or_phone" binding:"required"`
		VerificationCode int    `json:"verification_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	result, err := r.accountActivationService.Activate(req.EmailOrPhone, req.VerificationCode, c.ClientIP())
	if err != nil {
		c.JSON(result.StatusCode, gin.H{"status": "error", "message": result.Message, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "message": result.Message})
}
