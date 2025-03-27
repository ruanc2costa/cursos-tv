package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"tvtec/models"
	"tvtec/service"

	"github.com/gin-gonic/gin"
)

type AlunoController interface {
	ListarAlunos(c *gin.Context)
	ObterAlunoPorID(c *gin.Context)
	CriarAluno(c *gin.Context)
	AtualizarAluno(c *gin.Context)
	RemoverAluno(c *gin.Context)
	AdicionarAlunoCurso(c *gin.Context)
	CadastrarAlunoEInscrever(c *gin.Context)
	ListarInscricoesAluno(c *gin.Context)
}

type alunoController struct {
	alunoService service.AlunoService
}

func NewAlunoController(alunoService service.AlunoService) AlunoController {
	return &alunoController{alunoService: alunoService}
}

func (ctrl *alunoController) ListarAlunos(c *gin.Context) {
	alunos, err := ctrl.alunoService.ListarAlunos()
	if err != nil {
		log.Printf("Erro ao listar alunos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar alunos"})
		return
	}

	c.JSON(http.StatusOK, alunos)
}

func (ctrl *alunoController) ObterAlunoPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	aluno, err := ctrl.alunoService.ObterAlunoPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aluno)
}

func (ctrl *alunoController) CriarAluno(c *gin.Context) {
	var alunoDTO struct {
		Nome       string `json:"nome" binding:"required"`
		CPF        string `json:"cpf" binding:"required"`
		Email      string `json:"email" binding:"required"`
		Sexo       string `json:"sexo" binding:"required"`
		Telefone   string `json:"telefone"`
		DataNascto string `json:"dataNascto" binding:"required"`
	}

	if err := c.ShouldBindJSON(&alunoDTO); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Converter string de data para time.Time
	dataNascto, err := time.Parse("02/01/2006", alunoDTO.DataNascto)
	if err != nil {
		log.Printf("Erro ao converter data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA"})
		return
	}

	aluno := &models.Aluno{
		Nome:       alunoDTO.Nome,
		CPF:        alunoDTO.CPF,
		Email:      alunoDTO.Email,
		Sexo:       alunoDTO.Sexo,
		Telefone:   alunoDTO.Telefone,
		DataNascto: dataNascto,
	}

	if err := ctrl.alunoService.CriarAluno(aluno); err != nil {
		log.Printf("Erro ao criar aluno: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar aluno"})
		return
	}

	c.JSON(http.StatusCreated, aluno)
}

func (ctrl *alunoController) AtualizarAluno(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o aluno existe
	existingAluno, err := ctrl.alunoService.ObterAlunoPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var alunoDTO struct {
		Nome       string `json:"nome"`
		CPF        string `json:"cpf"`
		Email      string `json:"email"`
		Sexo       string `json:"sexo"`
		Telefone   string `json:"telefone"`
		DataNascto string `json:"dataNascto"`
	}

	if err := c.ShouldBindJSON(&alunoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Atualizar apenas os campos não vazios
	if alunoDTO.Nome != "" {
		existingAluno.Nome = alunoDTO.Nome
	}
	if alunoDTO.CPF != "" {
		existingAluno.CPF = alunoDTO.CPF
	}
	if alunoDTO.Email != "" {
		existingAluno.Email = alunoDTO.Email
	}
	if alunoDTO.Sexo != "" {
		existingAluno.Sexo = alunoDTO.Sexo
	}
	if alunoDTO.Telefone != "" {
		existingAluno.Telefone = alunoDTO.Telefone
	}
	if alunoDTO.DataNascto != "" {
		dataNascto, err := time.Parse("02/01/2006", alunoDTO.DataNascto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA"})
			return
		}
		existingAluno.DataNascto = dataNascto
	}

	if err := ctrl.alunoService.AtualizarAluno(existingAluno); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar aluno"})
		return
	}

	c.JSON(http.StatusOK, existingAluno)
}

func (ctrl *alunoController) RemoverAluno(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := ctrl.alunoService.RemoverAluno(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aluno removido com sucesso"})
}

func (ctrl *alunoController) AdicionarAlunoCurso(c *gin.Context) {
	// Capturar ID do aluno da URL
	alunoIDStr := c.Param("id")
	alunoID, err := strconv.ParseUint(alunoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do aluno inválido"})
		return
	}

	// Capturar ID do curso da URL
	cursoIDStr := c.Param("cursoId")
	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do curso inválido"})
		return
	}

	err = ctrl.alunoService.AdicionarAlunoCurso(uint(alunoID), uint(cursoID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aluno inscrito no curso com sucesso"})
}

func (ctrl *alunoController) CadastrarAlunoEInscrever(c *gin.Context) {
	var cadastroDTO struct {
		Aluno struct {
			Nome       string `json:"nome" binding:"required"`
			CPF        string `json:"cpf" binding:"required"`
			Email      string `json:"email" binding:"required"`
			Sexo       string `json:"sexo" binding:"required"`
			Telefone   string `json:"telefone"`
			DataNascto string `json:"dataNascto" binding:"required"`
		} `json:"aluno" binding:"required"`
		CursoID uint `json:"cursoId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&cadastroDTO); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Converter string de data para time.Time
	dataNascto, err := time.Parse("02/01/2006", cadastroDTO.Aluno.DataNascto)
	if err != nil {
		log.Printf("Erro ao converter data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA"})
		return
	}

	aluno := &models.Aluno{
		Nome:       cadastroDTO.Aluno.Nome,
		CPF:        cadastroDTO.Aluno.CPF,
		Email:      cadastroDTO.Aluno.Email,
		Sexo:       cadastroDTO.Aluno.Sexo,
		Telefone:   cadastroDTO.Aluno.Telefone,
		DataNascto: dataNascto,
	}

	if err := ctrl.alunoService.CadastrarAlunoEInscrever(aluno, cadastroDTO.CursoID); err != nil {
		log.Printf("Erro ao cadastrar aluno e inscrever no curso: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"aluno":   aluno,
		"message": "Aluno cadastrado e inscrito no curso com sucesso",
	})
}

func (ctrl *alunoController) ListarInscricoesAluno(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	inscricoes, err := ctrl.alunoService.ListarInscricoesAluno(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inscricoes)
}
