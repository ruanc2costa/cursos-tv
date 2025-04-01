package service

import (
	"errors"
	"fmt"
	"time"
	"tvtec/models"
	"tvtec/repository"
)

type AlunoService interface {
	// ListarAlunos retorna todos os alunos cadastrados no sistema.
	ListarAlunos() ([]models.Aluno, error)
	// ObterAlunoPorID busca um aluno pelo seu ID.
	ObterAlunoPorID(id uint) (*models.Aluno, error)
	// CriarAluno cadastra um novo aluno no sistema.
	CriarAluno(aluno *models.Aluno) error
	// AtualizarAluno atualiza os dados de um aluno existente.
	AtualizarAluno(aluno *models.Aluno) error
	// RemoverAluno remove um aluno e suas inscrições do sistema.
	RemoverAluno(id uint) error
	// CadastrarAlunoEInscrever cadastra um novo aluno e o inscreve em um curso.
	CadastrarAlunoEInscrever(aluno *models.Aluno, cursoID uint) error
	// AdicionarAlunoCurso inscreve um aluno existente em um curso.
	AdicionarAlunoCurso(alunoID, cursoID uint) error
	// ListarInscricoesAluno retorna todas as inscrições de um aluno.
	ListarInscricoesAluno(alunoID uint) ([]models.Inscricao, error)
}

type alunoService struct {
	alunoRepo     repository.AlunoRepository
	cursoRepo     repository.CursoRepository
	inscricaoRepo repository.InscricaoRepository
}

// NewAlunoService cria uma nova instância do serviço de alunos.
func NewAlunoService(
	alunoRepo repository.AlunoRepository,
	cursoRepo repository.CursoRepository,
	inscricaoRepo repository.InscricaoRepository,
) AlunoService {
	return &alunoService{
		alunoRepo:     alunoRepo,
		cursoRepo:     cursoRepo,
		inscricaoRepo: inscricaoRepo,
	}
}

// ListarAlunos retorna todos os alunos cadastrados no sistema.
func (s *alunoService) ListarAlunos() ([]models.Aluno, error) {
	return s.alunoRepo.FindAll()
}

// ObterAlunoPorID busca um aluno pelo seu ID.
func (s *alunoService) ObterAlunoPorID(id uint) (*models.Aluno, error) {
	return s.alunoRepo.FindByID(id)
}

// CriarAluno cadastra um novo aluno no sistema.
func (s *alunoService) CriarAluno(aluno *models.Aluno) error {
	return s.alunoRepo.Save(aluno)
}

// AtualizarAluno atualiza os dados de um aluno existente.
func (s *alunoService) AtualizarAluno(aluno *models.Aluno) error {
	return s.alunoRepo.Update(aluno)
}

// RemoverAluno remove um aluno e suas inscrições do sistema.
func (s *alunoService) RemoverAluno(id uint) error {
	// Verifica se o aluno existe
	aluno, err := s.alunoRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Busca inscrições do aluno (utilizando o método FindByAluno)
	inscricoes, err := s.inscricaoRepo.FindByAluno(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar inscrições do aluno: %w", err)
	}

	// Para cada inscrição, decrementa as vagas preenchidas no curso e remove a inscrição
	for _, inscricao := range inscricoes {
		if err := s.cursoRepo.DecrementarVagasPreenchidas(inscricao.CursoID); err != nil {
			return fmt.Errorf("erro ao decrementar vagas no curso %d: %w", inscricao.CursoID, err)
		}
		if err := s.inscricaoRepo.Delete(inscricao.ID); err != nil {
			return fmt.Errorf("erro ao remover inscrição %d: %w", inscricao.ID, err)
		}
	}

	// Remove o aluno
	return s.alunoRepo.Delete(aluno.ID)
}

// CadastrarAlunoEInscrever cadastra um novo aluno e o inscreve em um curso.
func (s *alunoService) CadastrarAlunoEInscrever(aluno *models.Aluno, cursoID uint) error {
	// Verifica se o curso existe e se há vagas
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar curso: %w", err)
	}
	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis no curso")
	}

	// Salva o aluno
	if err := s.alunoRepo.Save(aluno); err != nil {
		return fmt.Errorf("erro ao salvar aluno: %w", err)
	}

	// Cria a inscrição
	inscricao := &models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       cursoID,
		DataInscricao: time.Now(),
	}
	if err := s.inscricaoRepo.Save(inscricao); err != nil {
		_ = s.alunoRepo.Delete(aluno.ID)
		return fmt.Errorf("erro ao salvar inscrição: %w", err)
	}

	// Incrementa as vagas preenchidas no curso
	if err := s.cursoRepo.IncrementarVagasPreenchidas(cursoID); err != nil {
		_ = s.inscricaoRepo.Delete(inscricao.ID)
		_ = s.alunoRepo.Delete(aluno.ID)
		return fmt.Errorf("erro ao incrementar vagas no curso: %w", err)
	}

	return nil
}

// AdicionarAlunoCurso inscreve um aluno existente em um curso.
func (s *alunoService) AdicionarAlunoCurso(alunoID, cursoID uint) error {
	// Verifica se o aluno existe
	aluno, err := s.alunoRepo.FindByID(alunoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar aluno: %w", err)
	}

	// Verifica se o curso existe e se há vagas
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar curso: %w", err)
	}
	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis no curso")
	}

	// Verifica se o aluno já está inscrito no curso
	_, err = s.inscricaoRepo.FindByAlunoECurso(alunoID, cursoID)
	if err == nil {
		return errors.New("aluno já está inscrito neste curso")
	}

	// Cria a inscrição
	inscricao := &models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       curso.ID,
		DataInscricao: time.Now(),
	}
	if err := s.inscricaoRepo.Save(inscricao); err != nil {
		return fmt.Errorf("erro ao salvar inscrição: %w", err)
	}

	// Incrementa as vagas preenchidas no curso
	if err := s.cursoRepo.IncrementarVagasPreenchidas(cursoID); err != nil {
		_ = s.inscricaoRepo.Delete(inscricao.ID)
		return fmt.Errorf("erro ao incrementar vagas no curso: %w", err)
	}

	return nil
}

// ListarInscricoesAluno retorna todas as inscrições de um aluno.
func (s *alunoService) ListarInscricoesAluno(alunoID uint) ([]models.Inscricao, error) {
	// Verifica se o aluno existe
	_, err := s.alunoRepo.FindByID(alunoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aluno: %w", err)
	}

	// Retorna as inscrições do aluno (utilizando o método FindByAluno)
	inscricoes, err := s.inscricaoRepo.FindByAluno(alunoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar inscrições: %w", err)
	}

	return inscricoes, nil
}
