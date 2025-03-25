package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"
	"tvtec/repository"
	"tvtec/service"
)

// AlunoController gerencia os endpoints relacionados à entidade Aluno.
type AlunoController struct {
	alunoService    *service.AlunoService
	alunoRepository *repository.AlunoRepository
}

func NewAlunoController(as *service.AlunoService, ar *repository.AlunoRepository) *AlunoController {
	return &AlunoController{
		alunoService:    as,
		alunoRepository: ar,
	}
}

func (ctrl *AlunoController) RegisterRoutes(r *gin.Engine) {
	grupo := r.Group("/aluno")
	{
		grupo.GET("", ctrl.GetAllAlunos)
		grupo.GET("/:id", ctrl.GetAlunoByID)
		grupo.POST("", ctrl.PostAluno)
		grupo.DELETE("/:id", ctrl.DeleteAluno)
	}
}

func (ctrl *AlunoController) GetAllAlunos(c *gin.Context) {
	alunos, err := ctrl.alunoService.GetAllAlunos()
	if err != nil {
		// Adiciona o erro ao contexto e retorna (o middleware de erro cuidará da resposta)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, alunos)
}

func (ctrl *AlunoController) GetAlunoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	aluno, err := ctrl.alunoService.GetAluno(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, aluno)
}

func (ctrl *AlunoController) PostAluno(c *gin.Context) {
	var aluno models.Aluno
	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.Error(err)
		return
	}
	novoAluno, err := ctrl.alunoService.AddAluno(&aluno)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, novoAluno)
}

func (ctrl *AlunoController) DeleteAluno(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	if err := ctrl.alunoService.DeleteAluno(uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
