package service

import (
	"errors"
	"log"
	"time"
	"tvtec/models"
	"tvtec/repository"
)

// Interface para o serviço de inscrições
type InscricaoService interface {
	ListarInscricoesDetalhadas() ([]models.Inscricao, error)
	ObterInscricaoPorID(id uint) (*models.Inscricao, error)
	CriarInscricao(inscricao *models.Inscricao) error
	CancelarInscricao(id uint) error
	ListarInscricoesPorAluno(alunoID uint) ([]models.Inscricao, error)
	ListarInscricoesPorCurso(cursoID uint) ([]models.Inscricao, error)
	GerarRelatorio(dados []map[string]interface{}) error
}

// Implementação concreta do serviço
type inscricaoServiceImpl struct {
	inscricaoRepo repository.InscricaoRepository
}

// Função construtora
func NewInscricaoService(inscricaoRepo repository.InscricaoRepository) InscricaoService {
	return &inscricaoServiceImpl{
		inscricaoRepo: inscricaoRepo,
	}
}

// ListarInscricoesDetalhadas recupera todas as inscrições com informações de alunos e cursos
func (s *inscricaoServiceImpl) ListarInscricoesDetalhadas() ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindAllWithDetails()
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições detalhadas")
	}

	return inscricoes, nil
}

// ObterInscricaoPorID busca uma inscrição específica com detalhes
func (s *inscricaoServiceImpl) ObterInscricaoPorID(id uint) (*models.Inscricao, error) {
	inscricao, err := s.inscricaoRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, errors.New("inscrição não encontrada")
	}

	return inscricao, nil
}

// CriarInscricao registra uma nova inscrição no sistema
func (s *inscricaoServiceImpl) CriarInscricao(inscricao *models.Inscricao) error {
	if inscricao.AlunoID == 0 || inscricao.CursoID == 0 {
		return errors.New("aluno e curso são obrigatórios para uma inscrição")
	}

	// Definir valores padrão para campos opcionais se não forem informados
	if inscricao.EhPCD == "" {
		inscricao.EhPCD = "N"
	}

	// Definir a data de inscrição como a data atual se não for informada
	if inscricao.DataInscricao.IsZero() {
		inscricao.DataInscricao = time.Now()
	}

	// Validações específicas para os novos campos
	if inscricao.EhPCD == "S" && inscricao.TipoPCD == "" {
		return errors.New("quando marcado como PCD, o tipo de deficiência deve ser informado")
	}

	if inscricao.LevaNotebook == "" {
		inscricao.LevaNotebook = "N"
	}

	return s.inscricaoRepo.Save(inscricao)
}

// CancelarInscricao remove uma inscrição e atualiza vagas do curso
func (s *inscricaoServiceImpl) CancelarInscricao(id uint) error {
	// Verificar se a inscrição existe
	_, err := s.inscricaoRepo.FindByID(id)
	if err != nil {
		return errors.New("inscrição não encontrada")
	}

	// Deletar pelo ID, não pelo objeto
	return s.inscricaoRepo.Delete(id)
}

// ListarInscricoesPorAluno retorna todas as inscrições de um aluno específico
func (s *inscricaoServiceImpl) ListarInscricoesPorAluno(alunoID uint) ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindByAlunoWithDetails(alunoID)
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições do aluno")
	}

	return inscricoes, nil
}

// ListarInscricoesPorCurso retorna todas as inscrições de um curso específico
func (s *inscricaoServiceImpl) ListarInscricoesPorCurso(cursoID uint) ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindByCursoWithDetails(cursoID)
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições do curso")
	}

	return inscricoes, nil
}

// GerarRelatorio processa os dados de inscrições para análise ou exportação
func (s *inscricaoServiceImpl) GerarRelatorio(dados []map[string]interface{}) error {
	// Implementação simplificada - você pode expandir conforme necessário
	if len(dados) == 0 {
		return errors.New("nenhum dado fornecido para gerar relatório")
	}

	// Aqui seria possível implementar um relatório mais elaborado, com estatísticas
	// sobre os novos campos como quantidade de PCDs, bairros mais frequentes,
	// distribuição por escolaridade, etc.

	log.Printf("Gerando relatório com %d registros", len(dados))

	// Registrar informações sobre os dados
	var comPCD int
	var semPCD int
	var cuidadores int
	var necessitamElevador int
	var autorizaramWhatsApp int

	for _, dado := range dados {
		if ehPCD, ok := dado["ehPCD"].(string); ok && ehPCD == "sim" {
			comPCD++
		} else {
			semPCD++
		}

		if ehCuidador, ok := dado["ehCuidador"].(string); ok && ehCuidador == "sim" {
			cuidadores++
		}

		if necessitaElevador, ok := dado["necessitaElevador"].(string); ok && necessitaElevador == "sim" {
			necessitamElevador++
		}

		if autorizaWhatsApp, ok := dado["autorizaWhatsApp"].(string); ok && autorizaWhatsApp == "sim" {
			autorizaramWhatsApp++
		}
	}

	log.Printf("Resumo do relatório: %d PCDs, %d cuidadores, %d necessitam de elevador, %d autorizaram WhatsApp",
		comPCD, cuidadores, necessitamElevador, autorizaramWhatsApp)

	return nil
}
