package controller

import (
	"net/http"
	"tvtec/middleware"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type authController struct {
}

func NewAuthController() AuthController {
	return &authController{}
}

// Estrutura para receber os dados de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Estrutura para resposta de login
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Login autentica um usuário e retorna um token JWT
func (ctrl *authController) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de login inválidos"})
		return
	}

	// Verificar credenciais com variáveis do middleware
	if loginRequest.Username == middleware.AdminUsername && loginRequest.Password == middleware.AdminPassword {
		// Gerar token
		token, err := middleware.GenerateToken(loginRequest.Username, "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
			return
		}

		// Retornar resposta
		c.JSON(http.StatusOK, LoginResponse{
			Token:    token,
			Username: loginRequest.Username,
			Role:     "admin",
		})
		return
	}

	// Credenciais inválidas
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
}

// ValidateToken verifica se um token é válido
func (ctrl *authController) ValidateToken(c *gin.Context) {
	// O middleware de autenticação já verificou o token
	// Apenas retorne os dados do usuário
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"role":     role,
	})
}
