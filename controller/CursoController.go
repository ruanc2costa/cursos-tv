package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"   // ajuste para o caminho correto dos seus models
	"tvtec/services" // ajuste para o caminho correto dos seus serviços
)

// CursoController gerencia os endpoints relacionados à entidade Curso.
type CursoController struct {
	cursoService *service.CursoService
}

// NewCursoController cria uma nova instância do controller.
func NewCursoController(cs *service.CursoService) *CursoController {
	return &CursoController{cursoService: cs}
}

// RegisterRoutes registra as rotas do controller para a entidade Curso.
func (ctrl *CursoController) RegisterRoutes(r *gin.Engine) {
	grupo := r.Group("/curso")
	{
		grupo.GET("", ctrl.GetAllCursos)
		grupo.GET("/:id", ctrl.GetCursoByID)
		grupo.POST("", ctrl.AddCurso)
		grupo.PUT("/:id", ctrl.UpdateCurso)
		grupo.DELETE("/:id", ctrl.DeleteCurso)
	}
}

// GetAllCursos retorna todos os cursos.
func (ctrl *CursoController) GetAllCursos(c *gin.Context) {
	cursos, err := ctrl.cursoService.GetAllCursos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cursos)
}

// GetCursoByID retorna um curso específico pelo ID.
func (ctrl *CursoController) GetCursoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	curso, err := ctrl.cursoService.GetCurso(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, curso)
}

// AddCurso cria um novo curso.
func (ctrl *CursoController) AddCurso(c *gin.Context) {
	var curso models.Curso
	if err := c.ShouldBindJSON(&curso); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	novoCurso, err := ctrl.cursoService.AddCurso(&curso)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, novoCurso)
}

// UpdateCurso atualiza os dados de um curso existente.
func (ctrl *CursoController) UpdateCurso(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var curso models.Curso
	if err := c.ShouldBindJSON(&curso); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cursoAtualizado, err := ctrl.cursoService.UpdateCurso(uint(id), &curso)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cursoAtualizado)
}

// DeleteCurso remove um curso com base no ID.
func (ctrl *CursoController) DeleteCurso(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if err := ctrl.cursoService.DeleteCurso(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
