package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Variáveis para armazenar credenciais e chave JWT
var (
	AdminUsername string
	AdminPassword string
	SecretKey     string
)

// InitAuthConfig inicializa as configurações de autenticação a partir de variáveis de ambiente
func InitAuthConfig() {
	// Carregar variáveis de ambiente
	AdminUsername = getEnvWithDefault("ADMIN_USERNAME", "admin")
	AdminPassword = getEnvWithDefault("ADMIN_PASSWORD", "admin123") // Valor padrão somente para fallback
	SecretKey = getEnvWithDefault("JWT_SECRET_KEY", "chave-secreta-padrao-mudar-em-producao")

	// Avisar se valores padrão estão sendo usados em produção
	if gin.Mode() == gin.ReleaseMode {
		if AdminPassword == "admin123" {
			log.Println("AVISO: Senha de administrador padrão sendo usada em ambiente de produção!")
		}
		if SecretKey == "chave-secreta-padrao-mudar-em-producao" {
			log.Println("AVISO: Chave JWT padrão sendo usada em ambiente de produção!")
		}
	}

	log.Println("Configurações de autenticação inicializadas com sucesso")
}

// Helper para obter variável de ambiente com valor padrão
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// UserClaims define os dados armazenados no token JWT
type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware é o middleware para verificar autenticação
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autorização necessária"})
			c.Abort()
			return
		}

		// Formato esperado: "Bearer TOKEN"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de autorização inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verificar e validar o token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verificar algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verificar se o token é válido
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// Armazenar informações do usuário no contexto
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
	}
}

// AdminAuthMiddleware verifica se o usuário é um administrador
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primeiro verifica se está autenticado
		AuthMiddleware()(c)

		// Se a requisição foi abortada pelo middleware anterior, retorne
		if c.IsAborted() {
			return
		}

		// Verifica se o usuário tem o papel de admin
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado: requer privilégios de administrador"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken gera um token JWT para o usuário
func GenerateToken(username, role string) (string, error) {
	// Define os claims do token
	claims := UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token válido por 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Cria o token com os claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
