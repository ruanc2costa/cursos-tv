package service

import (
	"errors"
	"log"
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

// Implementação concreta do serviço - USANDO NOME DIFERENTE
type inscricaoServiceImpl struct {
	inscricaoRepo repository.InscricaoRepository
}

// Função construtora
func NewInscricaoService(inscricaoRepo repository.InscricaoRepository) InscricaoService {
	return &inscricaoServiceImpl{
		inscricaoRepo: inscricaoRepo,
	}
}

// Implementação dos métodos
func (s *inscricaoServiceImpl) ListarInscricoesDetalhadas() ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindAllWithDetails()
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições detalhadas")
	}

	return inscricoes, nil
}

func (s *inscricaoServiceImpl) ObterInscricaoPorID(id uint) (*models.Inscricao, error) {
	inscricao, err := s.inscricaoRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, errors.New("inscrição não encontrada")
	}

	return inscricao, nil
}

func (s *inscricaoServiceImpl) CriarInscricao(inscricao *models.Inscricao) error {
	if inscricao.AlunoID == 0 || inscricao.CursoID == 0 {
		return errors.New("aluno e curso são obrigatórios para uma inscrição")
	}

	return s.inscricaoRepo.Save(inscricao)
}

func (s *inscricaoServiceImpl) CancelarInscricao(id uint) error {
	// Verificar se a inscrição existe
	_, err := s.inscricaoRepo.FindByID(id)
	if err != nil {
		return errors.New("inscrição não encontrada")
	}

	// Deletar pelo ID, não pelo objeto
	return s.inscricaoRepo.Delete(id)
}

func (s *inscricaoServiceImpl) ListarInscricoesPorAluno(alunoID uint) ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindByAlunoWithDetails(alunoID)
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições do aluno")
	}

	return inscricoes, nil
}

func (s *inscricaoServiceImpl) ListarInscricoesPorCurso(cursoID uint) ([]models.Inscricao, error) {
	inscricoes, err := s.inscricaoRepo.FindByCursoWithDetails(cursoID)
	if err != nil {
		return nil, errors.New("falha ao recuperar inscrições do curso")
	}

	return inscricoes, nil
}

func (s *inscricaoServiceImpl) GerarRelatorio(dados []map[string]interface{}) error {
	// Implementação simplificada - você pode expandir conforme necessário
	if len(dados) == 0 {
		return errors.New("nenhum dado fornecido para gerar relatório")
	}

	// Aqui você poderia implementar lógica para salvar os dados em um formato específico,
	// gerar estatísticas, ou processar os dados conforme necessário para sua aplicação
	log.Printf("Gerando relatório com %d registros", len(dados))

	return nil
}
