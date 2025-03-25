package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse define o formato padrão de resposta para erros.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorMiddleware captura os erros adicionados ao contexto e retorna uma resposta padrão.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Se houver erros registrados no contexto...
		if len(c.Errors) > 0 {
			// Podemos escolher o último erro registrado
			err := c.Errors.Last().Err
			// Aqui, você pode implementar lógica para escolher o status code
			statusCode := http.StatusInternalServerError

			// Loga o erro para depuração
			log.Printf("Erro capturado: %v", err)

			// Retorna a resposta padronizada
			c.JSON(statusCode, ErrorResponse{Error: err.Error()})
			// Opcionalmente, interrompe a chain
			c.Abort()
		}
	}
}
