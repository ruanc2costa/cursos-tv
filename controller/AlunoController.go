package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"
	"tvtec/service"
)

// AlunoController gerencia as rotas e lógica HTTP para alunos
type AlunoController struct {
	service *service.AlunoService
}

// NewAlunoController cria uma nova instância do controlador de alunos
func NewAlunoController(service *service.AlunoService) *AlunoController {
	return &AlunoController{service: service}
}

// CriarAluno lida com a criação de um novo aluno
func (c *AlunoController) CriarAluno(ctx *gin.Context) {
	var aluno models.Aluno

	// Vincula os dados JSON da requisição ao modelo Aluno
	if err := ctx.ShouldBindJSON(&aluno); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Chama o serviço para criar o aluno
	if err := c.service.CriarAluno(&aluno); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao criar aluno",
			"details": err.Error(),
		})
		return
	}

	// Responde com sucesso
	ctx.JSON(http.StatusCreated, aluno)
}

// AdicionarAlunoCurso lida com a inclusão de um aluno em um curso
func (c *AlunoController) AdicionarAlunoCurso(ctx *gin.Context) {
	var aluno models.Aluno
	cursoIDStr := ctx.Param("cursoId")

	// Converte o ID do curso para uint
	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	// Vincula os dados JSON da requisição ao modelo Aluno
	if err := ctx.ShouldBindJSON(&aluno); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Chama o serviço para adicionar aluno ao curso
	if err := c.service.AdicionarAluno(&aluno, uint(cursoID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao adicionar aluno ao curso",
			"details": err.Error(),
		})
		return
	}

	// Responde com sucesso
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Aluno adicionado ao curso com sucesso",
		"aluno":   aluno,
	})
}

// ListarAlunos recupera todos os alunos
func (c *AlunoController) ListarAlunos(ctx *gin.Context) {
	alunos, err := c.service.ListarAlunos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar alunos",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alunos)
}

// ObterAlunoPorID busca um aluno específico
func (c *AlunoController) ObterAlunoPorID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de aluno inválido",
		})
		return
	}

	aluno, err := c.service.ObterAlunoPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Aluno não encontrado",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, aluno)
}

// AtualizarAluno atualiza informações de um aluno
func (c *AlunoController) AtualizarAluno(ctx *gin.Context) {
	var aluno models.Aluno

	// Vincula os dados JSON da requisição ao modelo Aluno
	if err := ctx.ShouldBindJSON(&aluno); err != nil {
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
			"error": "ID de aluno inválido",
		})
		return
	}

	// Define o ID do aluno
	aluno.ID = uint(id)

	// Chama o serviço para atualizar o aluno
	if err := c.service.AtualizarAluno(&aluno); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao atualizar aluno",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, aluno)
}

// RemoverAluno exclui um aluno
func (c *AlunoController) RemoverAluno(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de aluno inválido",
		})
		return
	}

	// Chama o serviço para remover o aluno
	if err := c.service.RemoverAluno(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao remover aluno",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Aluno removido com sucesso",
	})
}
