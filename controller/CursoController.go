package controller

import (
	"log"
	"net/http"
	"strconv"
	"tvtec/models"
	"tvtec/service"

	"github.com/gin-gonic/gin"
)

type CursoController interface {
	ListarCursos(c *gin.Context)
	ObterCursoPorID(c *gin.Context)
	CriarCurso(c *gin.Context)
	AtualizarCurso(c *gin.Context)
	RemoverCurso(c *gin.Context)
	VerificarDisponibilidadeVagas(c *gin.Context)
	ListarInscricoesCurso(c *gin.Context)
}

type cursoController struct {
	cursoService service.CursoService
}

func NewCursoController(cursoService service.CursoService) CursoController {
	return &cursoController{cursoService: cursoService}
}

func (ctrl *cursoController) ListarCursos(c *gin.Context) {
	cursos, err := ctrl.cursoService.ListarCursos()
	if err != nil {
		log.Printf("Erro ao listar cursos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar cursos"})
		return
	}

	c.JSON(http.StatusOK, cursos)
}

func (ctrl *cursoController) ObterCursoPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	curso, err := ctrl.cursoService.ObterCursoPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curso)
}

func (ctrl *cursoController) CriarCurso(c *gin.Context) {
	var cursoDTO struct {
		Nome         string `json:"nome" binding:"required"`
		Professor    string `json:"professor" binding:"required"`
		Data         string `json:"data" binding:"required"`
		CargaHoraria int32  `json:"cargaHoraria" binding:"required"`
		Certificado  string `json:"certificado" binding:"required"`
		VagasTotais  int32  `json:"vagasTotais" binding:"required"`
	}

	if err := c.ShouldBindJSON(&cursoDTO); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Criar CustomTime a partir da string de data
	var data models.CustomTime
	if err := data.UnmarshalJSON([]byte(`"` + cursoDTO.Data + `"`)); err != nil {
		log.Printf("Erro ao converter data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA"})
		return
	}

	curso := &models.Curso{
		Nome:         cursoDTO.Nome,
		Professor:    cursoDTO.Professor,
		Data:         data,
		CargaHoraria: cursoDTO.CargaHoraria,
		Certificado:  cursoDTO.Certificado,
		VagasTotais:  cursoDTO.VagasTotais,
	}

	if err := ctrl.cursoService.CriarCurso(curso); err != nil {
		log.Printf("Erro ao criar curso: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar curso"})
		return
	}

	c.JSON(http.StatusCreated, curso)
}

func (ctrl *cursoController) AtualizarCurso(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o curso existe
	existingCurso, err := ctrl.cursoService.ObterCursoPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var cursoDTO struct {
		Nome         string `json:"nome"`
		Professor    string `json:"professor"`
		Data         string `json:"data"`
		CargaHoraria *int32 `json:"cargaHoraria"`
		Certificado  string `json:"certificado"`
		VagasTotais  *int32 `json:"vagasTotais"`
	}

	if err := c.ShouldBindJSON(&cursoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Atualizar apenas os campos não vazios
	if cursoDTO.Nome != "" {
		existingCurso.Nome = cursoDTO.Nome
	}
	if cursoDTO.Professor != "" {
		existingCurso.Professor = cursoDTO.Professor
	}
	if cursoDTO.Data != "" {
		var data models.CustomTime
		if err := data.UnmarshalJSON([]byte(`"` + cursoDTO.Data + `"`)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA"})
			return
		}
		existingCurso.Data = data
	}
	if cursoDTO.CargaHoraria != nil {
		existingCurso.CargaHoraria = *cursoDTO.CargaHoraria
	}
	if cursoDTO.Certificado != "" {
		existingCurso.Certificado = cursoDTO.Certificado
	}
	if cursoDTO.VagasTotais != nil {
		// Verificamos se o novo número de vagas totais é pelo menos o número de vagas já preenchidas
		if *cursoDTO.VagasTotais < existingCurso.VagasPreenchidas {
			c.JSON(http.StatusBadRequest, gin.H{"error": "O número de vagas totais não pode ser menor que o número de vagas já preenchidas"})
			return
		}
		existingCurso.VagasTotais = *cursoDTO.VagasTotais
	}

	if err := ctrl.cursoService.AtualizarCurso(existingCurso); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar curso"})
		return
	}

	c.JSON(http.StatusOK, existingCurso)
}

func (ctrl *cursoController) RemoverCurso(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := ctrl.cursoService.RemoverCurso(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curso removido com sucesso"})
}

func (ctrl *cursoController) VerificarDisponibilidadeVagas(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	vagasDisponiveis, err := ctrl.cursoService.VerificarDisponibilidadeVagas(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vagasDisponiveis": vagasDisponiveis,
	})
}

func (ctrl *cursoController) ListarInscricoesCurso(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	inscricoes, err := ctrl.cursoService.ListarInscricoesCurso(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inscricoes)
}
