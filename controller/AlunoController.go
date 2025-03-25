package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"     // ajuste para o caminho correto dos seus models
	"tvtec/repository" // ajuste para o caminho correto dos seus repositórios
	"tvtec/services"   // ajuste para o caminho correto dos seus serviços
)

// AlunoController gerencia os endpoints relacionados à entidade Aluno.
type AlunoController struct {
	alunoService    *service.AlunoService
	alunoRepository *repository.AlunoRepository
}

// NewAlunoController cria uma nova instância do controller.
func NewAlunoController(as *service.AlunoService, ar *repository.AlunoRepository) *AlunoController {
	return &AlunoController{
		alunoService:    as,
		alunoRepository: ar,
	}
}

// RegisterRoutes registra as rotas do controller.
func (ctrl *AlunoController) RegisterRoutes(r *gin.Engine) {
	grupo := r.Group("/aluno")
	{
		grupo.GET("", ctrl.GetAllAlunos)
		grupo.GET("/:id", ctrl.GetAlunoByID)
		grupo.POST("/alunos", ctrl.PostAlunos)
		grupo.DELETE("/:id", ctrl.DeleteAluno)
	}
}

// GetAllAlunos retorna todos os alunos.
func (ctrl *AlunoController) GetAllAlunos(c *gin.Context) {
	alunos, err := ctrl.alunoService.GetAllAlunos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alunos)
}

// GetAlunoByID retorna um aluno específico pelo ID.
func (ctrl *AlunoController) GetAlunoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	aluno, err := ctrl.alunoService.GetAluno(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aluno)
}

// PostAlunos cria ou atualiza um aluno baseado no email e adiciona um curso.
func (ctrl *AlunoController) PostAlunos(c *gin.Context) {
	var alunoRequest models.Aluno
	if err := c.ShouldBindJSON(&alunoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se há pelo menos um curso informado.
	if len(alunoRequest.Cursos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário informar pelo menos um curso"})
		return
	}
	cursoSelecionado := alunoRequest.Cursos[0]

	// Busca aluno pelo email
	alunoExistente, err := ctrl.alunoRepository.FindByEmail(alunoRequest.Email)
	var aluno models.Aluno
	if err == nil && alunoExistente != nil {
		// Aluno já existe, adiciona o curso se ainda não estiver associado
		aluno = *alunoExistente
		exists := false
		for _, curso := range aluno.Cursos {
			if curso.ID == cursoSelecionado.ID {
				exists = true
				break
			}
		}
		if !exists {
			aluno.Cursos = append(aluno.Cursos, cursoSelecionado)
			if err = ctrl.alunoRepository.Save(&aluno); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		// Cria novo aluno e associa o curso informado.
		aluno = models.Aluno{
			Nome:   alunoRequest.Nome,
			Email:  alunoRequest.Email,
			Cursos: []models.Curso{cursoSelecionado},
		}
		if err = ctrl.alunoRepository.Save(&aluno); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, aluno)
}

// DeleteAluno remove um aluno com base no ID.
func (ctrl *AlunoController) DeleteAluno(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := ctrl.alunoService.DeleteAluno(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
