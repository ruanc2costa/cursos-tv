package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tvtec/models"
	"tvtec/service"
)

type InscricaoController struct {
	service service.InscricaoService
}

func NewInscricaoController(service service.InscricaoService) *InscricaoController {
	return &InscricaoController{service: service}
}

func (c *InscricaoController) ListarInscricoes(ctx *gin.Context) {
	inscricoes, err := c.service.ListarInscricoesDetalhadas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar inscrições",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, inscricoes)
}

func (c *InscricaoController) ObterInscricaoPorID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de inscrição inválido",
		})
		return
	}

	inscricao, err := c.service.ObterInscricaoPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Inscrição não encontrada",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, inscricao)
}

func (c *InscricaoController) CriarInscricao(ctx *gin.Context) {
	var inscricao models.Inscricao

	if err := ctx.ShouldBindJSON(&inscricao); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	if err := c.service.CriarInscricao(&inscricao); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao criar inscrição",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, inscricao)
}

func (c *InscricaoController) CancelarInscricao(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de inscrição inválido",
		})
		return
	}

	if err := c.service.CancelarInscricao(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao cancelar inscrição",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Inscrição cancelada com sucesso",
	})
}

func (c *InscricaoController) ListarInscricoesPorAluno(ctx *gin.Context) {
	idStr := ctx.Param("alunoId")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de aluno inválido",
		})
		return
	}

	inscricoes, err := c.service.ListarInscricoesPorAluno(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar inscrições do aluno",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, inscricoes)
}

func (c *InscricaoController) ListarInscricoesPorCurso(ctx *gin.Context) {
	idStr := ctx.Param("cursoId")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	inscricoes, err := c.service.ListarInscricoesPorCurso(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar inscrições do curso",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, inscricoes)
}

func (c *InscricaoController) GerarRelatorio(ctx *gin.Context) {
	var dados []map[string]interface{}
	if err := ctx.ShouldBindJSON(&dados); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao processar dados do relatório",
			"details": err.Error(),
		})
		return
	}

	if err := c.service.GerarRelatorio(dados); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao gerar relatório",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Relatório gerado com sucesso",
		"timestamp": ctx.Request.Context().Value("requestTime"),
	})
}
