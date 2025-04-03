package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"tvtec/models"
	"tvtec/service"
)

// AlunoController gerencia as rotas e lógica HTTP para alunos
type AlunoController struct {
	service service.AlunoService
}

// NewAlunoController cria uma nova instância do controlador de alunos
func NewAlunoController(service service.AlunoService) *AlunoController {
	return &AlunoController{service: service}
}

// Estrutura para receber os dados do formulário
type InscricaoRequest struct {
	Nome              string `json:"nome"`
	CPF               string `json:"cpf"`
	Email             string `json:"email"`
	Curso             uint   `json:"curso"`
	Sexo              string `json:"sexo"`
	DataNascto        string `json:"dataNascto"`
	Telefone          string `json:"telefone"`
	Escolaridade      string `json:"escolaridade"`
	Trabalhando       string `json:"trabalhando"`
	Bairro            string `json:"bairro"`
	EhCuidador        string `json:"ehCuidador"`
	EhPCD             string `json:"ehPCD"`
	TipoPCD           string `json:"tipoPCD"`
	NecessitaElevador string `json:"necessitaElevador"`
	ComoSoube         string `json:"comoSoube"`
	AutorizaWhatsApp  string `json:"autorizaWhatsApp"`
}

// CadastrarAlunoEInscrever cadastra um novo aluno e o inscreve em um curso
func (c *AlunoController) CadastrarAlunoEInscrever(ctx *gin.Context) {
	var request InscricaoRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados de formulário inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validações básicas
	if request.Nome == "" || request.CPF == "" || request.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Nome, CPF e email são campos obrigatórios",
		})
		return
	}

	if request.Curso == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "É necessário selecionar um curso",
		})
		return
	}

	// Converter a data de nascimento de string para time.Time
	dataNascto, err := time.Parse("02/01/2006", request.DataNascto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Formato de data inválido. Use DD/MM/AAAA",
			"details": err.Error(),
		})
		return
	}

	// Criar o objeto Aluno
	aluno := &models.Aluno{
		Nome:       request.Nome,
		CPF:        request.CPF,
		Email:      request.Email,
		Sexo:       request.Sexo,
		Telefone:   request.Telefone,
		DataNascto: dataNascto,
	}

	// Criar o objeto de inscrição com os novos campos
	inscricao := &models.Inscricao{
		CursoID:           request.Curso,
		DataInscricao:     time.Now(),
		Escolaridade:      request.Escolaridade,
		Trabalhando:       request.Trabalhando,
		Bairro:            request.Bairro,
		EhCuidador:        request.EhCuidador,
		EhPCD:             request.EhPCD,
		TipoPCD:           request.TipoPCD,
		NecessitaElevador: request.NecessitaElevador,
		ComoSoube:         request.ComoSoube,
		AutorizaWhatsApp:  request.AutorizaWhatsApp,
	}

	// Chama o serviço para cadastrar o aluno e inscrevê-lo no curso
	if err := c.service.CadastrarAlunoEInscrever(aluno, inscricao); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Falha ao cadastrar aluno e inscrever no curso",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Aluno cadastrado e inscrito com sucesso",
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

// AdicionarAlunoCurso adiciona um aluno existente a um curso
func (c *AlunoController) AdicionarAlunoCurso(ctx *gin.Context) {
	alunoIDStr := ctx.Param("id")
	cursoIDStr := ctx.Param("cursoId")

	// Converte os IDs para uint
	alunoID, err := strconv.ParseUint(alunoIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de aluno inválido",
		})
		return
	}

	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de curso inválido",
		})
		return
	}

	// Dados adicionais para a inscrição
	var inscricaoData models.Inscricao
	if err := ctx.ShouldBindJSON(&inscricaoData); err == nil {
		// Se foram fornecidos dados adicionais, usa-os
		inscricaoData.AlunoID = uint(alunoID)
		inscricaoData.CursoID = uint(cursoID)
		inscricaoData.DataInscricao = time.Now()

		// Chama o serviço para criar a inscrição com detalhes
		if err := c.service.CriarInscricaoDetalhada(&inscricaoData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Falha ao adicionar aluno ao curso",
				"details": err.Error(),
			})
			return
		}
	} else {
		// Se não foram fornecidos dados adicionais, usa o método básico
		if err := c.service.AdicionarAlunoCurso(uint(alunoID), uint(cursoID)); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Falha ao adicionar aluno ao curso",
				"details": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Aluno adicionado ao curso com sucesso",
	})
}

// ListarInscricoesAluno lista todas as inscrições de um aluno
func (c *AlunoController) ListarInscricoesAluno(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Converte o ID para uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de aluno inválido",
		})
		return
	}

	inscricoes, err := c.service.ListarInscricoesAluno(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Falha ao recuperar inscrições",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, inscricoes)
}
