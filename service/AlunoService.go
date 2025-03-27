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
	// Primeiro verificamos se o aluno existe
	aluno, err := s.alunoRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Remover inscrições deste aluno
	inscricoes, err := s.inscricaoRepo.FindByAlunoID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar inscrições do aluno: %w", err)
	}

	// Para cada inscrição, decrementar o contador de vagas no curso
	for _, inscricao := range inscricoes {
		// Decrementar vagas preenchidas no curso
		if err := s.cursoRepo.DecrementarVagasPreenchidas(inscricao.CursoID); err != nil {
			return fmt.Errorf("erro ao decrementar vagas no curso %d: %w", inscricao.CursoID, err)
		}

		// Remover a inscrição
		if err := s.inscricaoRepo.Delete(inscricao.ID); err != nil {
			return fmt.Errorf("erro ao remover inscrição %d: %w", inscricao.ID, err)
		}
	}

	// Finalmente remover o aluno
	return s.alunoRepo.Delete(aluno.ID)
}

// CadastrarAlunoEInscrever cadastra um novo aluno e o inscreve em um curso.
func (s *alunoService) CadastrarAlunoEInscrever(aluno *models.Aluno, cursoID uint) error {
	// Verificar se o curso existe e tem vagas
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar curso: %w", err)
	}

	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis no curso")
	}

	// Salvar o aluno
	if err := s.alunoRepo.Save(aluno); err != nil {
		return fmt.Errorf("erro ao salvar aluno: %w", err)
	}

	// Criar uma inscrição
	inscricao := &models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       cursoID,
		DataInscricao: time.Now(),
	}

	if err := s.inscricaoRepo.Save(inscricao); err != nil {
		deleteErr := s.alunoRepo.Delete(aluno.ID)
		if deleteErr != nil {
			return fmt.Errorf("erro ao salvar inscrição: %w e erro ao remover aluno: %v", err, deleteErr)
		}
		return fmt.Errorf("erro ao salvar inscrição: %w", err)
	}

	// Incrementar vagas preenchidas no curso
	if err := s.cursoRepo.IncrementarVagasPreenchidas(cursoID); err != nil {
		// Se falhar ao incrementar vagas, devemos remover a inscrição e o aluno
		delInscErr := s.inscricaoRepo.Delete(inscricao.ID)
		delAlunoErr := s.alunoRepo.Delete(aluno.ID)

		if delInscErr != nil || delAlunoErr != nil {
			return fmt.Errorf("erro ao incrementar vagas: %w, erro ao remover inscrição: %v, erro ao remover aluno: %v",
				err, delInscErr, delAlunoErr)
		}

		return fmt.Errorf("erro ao incrementar vagas no curso: %w", err)
	}

	return nil
}

// AdicionarAlunoCurso inscreve um aluno existente em um curso.
func (s *alunoService) AdicionarAlunoCurso(alunoID, cursoID uint) error {
	// Verificar se o aluno existe
	aluno, err := s.alunoRepo.FindByID(alunoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar aluno: %w", err)
	}

	// Verificar se o curso existe e tem vagas
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar curso: %w", err)
	}

	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis no curso")
	}

	// Verificar se o aluno já está inscrito no curso
	_, err = s.inscricaoRepo.FindByAlunoECurso(alunoID, cursoID)
	if err == nil {
		return errors.New("aluno já está inscrito neste curso")
	}

	// Criar uma inscrição
	inscricao := &models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       curso.ID,
		DataInscricao: time.Now(),
	}

	if err := s.inscricaoRepo.Save(inscricao); err != nil {
		return fmt.Errorf("erro ao salvar inscrição: %w", err)
	}

	// Incrementar vagas preenchidas no curso
	if err := s.cursoRepo.IncrementarVagasPreenchidas(cursoID); err != nil {
		delErr := s.inscricaoRepo.Delete(inscricao.ID)
		if delErr != nil {
			return fmt.Errorf("erro ao incrementar vagas: %w e erro ao remover inscrição: %v", err, delErr)
		}
		return fmt.Errorf("erro ao incrementar vagas no curso: %w", err)
	}

	return nil
}

// ListarInscricoesAluno retorna todas as inscrições de um aluno.
func (s *alunoService) ListarInscricoesAluno(alunoID uint) ([]models.Inscricao, error) {
	// Verificar se o aluno existe
	_, err := s.alunoRepo.FindByID(alunoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aluno: %w", err)
	}

	// Retornar inscrições do aluno
	inscricoes, err := s.inscricaoRepo.FindByAlunoID(alunoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar inscrições: %w", err)
	}

	return inscricoes, nil
}
