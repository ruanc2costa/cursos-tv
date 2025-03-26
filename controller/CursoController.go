package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"
	"tvtec/service"
)

// CursoController gerencia as rotas e lógica HTTP para cursos
type CursoController struct {
	service *service.CursoService
}

// NewCursoController cria uma nova instância do controlador de cursos
func NewCursoController(service *service.CursoService) *CursoController {
	return &CursoController{service: service}
}

// CriarCurso lida com a criação de um novo curso
func (c *CursoController) CriarCurso(ctx *gin.Context) {
	var curso models.Curso

	// Vincula os dados JSON da requisição ao modelo Curso
	if err := ctx.ShouldBindJSON(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Chama o serviço para criar o curso
	if err := c.service.CriarCurso(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao criar curso",
			"details": err.Error(),
		})
		return
	}

	// Responde com sucesso e retorna o curso criado com ID gerado
	ctx.JSON(http.StatusCreated, curso)
}

// ListarCursos recupera todos os cursos
func (c *CursoController) ListarCursos(ctx *gin.Context) {
	cursos, err := c.service.ListarCursos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar cursos",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cursos)
}

// ObterCursoPorID busca um curso específico
func (c *CursoController) ObterCursoPorID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	curso, err := c.service.ObterCursoPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Curso não encontrado",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, curso)
}

// AtualizarCurso atualiza informações de um curso
func (c *CursoController) AtualizarCurso(ctx *gin.Context) {
	var curso models.Curso

	// Vincula os dados JSON da requisição ao modelo Curso
	if err := ctx.ShouldBindJSON(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Obtém o ID da URL
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	// Define o ID do curso
	curso.ID = uint(id)

	// Chama o serviço para atualizar o curso
	if err := c.service.AtualizarCurso(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao atualizar curso",
			"details": err.Error(),
		})
		return
	}

	// Recupera o curso atualizado para retornar ao cliente
	cursoAtualizado, err := c.service.ObterCursoPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Curso atualizado, mas falha ao recuperar dados atualizados",
		})
		return
	}

	ctx.JSON(http.StatusOK, cursoAtualizado)
}

// RemoverCurso exclui um curso
func (c *CursoController) RemoverCurso(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	// Chama o serviço para remover o curso
	if err := c.service.RemoverCurso(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao remover curso",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Curso removido com sucesso",
	})
}

// VerificarDisponibilidadeVagas verifica se um curso tem vagas disponíveis
func (c *CursoController) VerificarDisponibilidadeVagas(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	// Verifica disponibilidade de vagas
	disponivel, err := c.service.VerificarDisponibilidadeVagas(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Falha ao verificar disponibilidade",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"disponivel": disponivel,
	})
}
