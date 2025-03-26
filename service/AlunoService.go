package service

import (
	"errors"
	"time"

	"tvtec/models"
	"tvtec/repository"
)

// AlunoService contém a lógica de negócio para manipulação de alunos
type AlunoService struct {
	alunoRepo     *repository.AlunoRepository
	cursoRepo     *repository.CursoRepository
	inscricaoRepo *repository.InscricaoRepository
}

// NewAlunoService cria uma nova instância do serviço de alunos
func NewAlunoService(
	alunoRepo *repository.AlunoRepository,
	cursoRepo *repository.CursoRepository,
	inscricaoRepo *repository.InscricaoRepository,
) *AlunoService {
	return &AlunoService{
		alunoRepo:     alunoRepo,
		cursoRepo:     cursoRepo,
		inscricaoRepo: inscricaoRepo,
	}
}

// AdicionarAluno gerencia a inclusão de um aluno em um curso
func (s *AlunoService) AdicionarAluno(aluno *models.Aluno, cursoID uint) error {
	// Iniciar transação para garantir atomicidade
	return s.alunoRepo.db.Transaction(func(tx *gorm.DB) error {
		// Verificar se o curso existe
		curso, err := s.cursoRepo.FindByID(cursoID)
		if err != nil {
			return errors.New("curso não encontrado")
		}

		// Verificar disponibilidade de vagas
		if curso.VagasPreenchidas >= curso.VagasTotais {
			return errors.New("não há vagas disponíveis para este curso")
		}

		// Verificar se o aluno já existe pelo email
		existingAluno, err := s.alunoRepo.FindByEmail(aluno.Email)

		if err != nil {
			// Aluno não existe, criar novo
			if err := s.alunoRepo.Save(aluno); err != nil {
				return err
			}
		} else {
			// Aluno já existe, usar o aluno existente
			aluno = existingAluno
		}

		// Verificar se aluno já está inscrito neste curso
		var inscricaoExistente models.Inscricao
		err = tx.Where("aluno_id = ? AND curso_id = ?", aluno.ID, cursoID).First(&inscricaoExistente).Error
		if err == nil {
			return errors.New("aluno já está inscrito neste curso")
		}

		// Criar nova inscrição
		inscricao := &models.Inscricao{
			AlunoID:       aluno.ID,
			CursoID:       cursoID,
			DataInscricao: time.Now(),
		}

		// Salvar inscrição
		if err := tx.Create(inscricao).Error; err != nil {
			return err
		}

		// Atualizar vagas do curso
		curso.VagasPreenchidas++
		if err := tx.Save(curso).Error; err != nil {
			return err
		}

		return nil
	})
}

// ObterAlunoPorID busca um aluno específico
func (s *AlunoService) ObterAlunoPorID(id uint) (*models.Aluno, error) {
	aluno, err := s.alunoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("aluno não encontrado")
	}
	return aluno, nil
}

// ListarAlunos recupera todos os alunos
func (s *AlunoService) ListarAlunos() ([]models.Aluno, error) {
	return s.alunoRepo.FindAll()
}

// CriarAluno realiza a criação de um novo aluno
func (s *AlunoService) CriarAluno(aluno *models.Aluno) error {
	// Validações básicas
	if aluno.Nome == "" {
		return errors.New("nome do aluno é obrigatório")
	}

	if aluno.Email == "" {
		return errors.New("email do aluno é obrigatório")
	}

	// Verificar se já existe aluno com este email
	existente, _ := s.alunoRepo.FindByEmail(aluno.Email)
	if existente != nil {
		return errors.New("já existe um aluno cadastrado com este email")
	}

	return s.alunoRepo.Save(aluno)
}

// AtualizarAluno atualiza as informações de um aluno existente
func (s *AlunoService) AtualizarAluno(aluno *models.Aluno) error {
	// Verificar se o aluno existe
	existente, err := s.alunoRepo.FindByID(aluno.ID)
	if err != nil {
		return errors.New("aluno não encontrado")
	}

	// Manter campos importantes do registro original
	if aluno.Nome == "" {
		aluno.Nome = existente.Nome
	}

	if aluno.Email == "" {
		aluno.Email = existente.Email
	}

	return s.alunoRepo.Save(aluno)
}

// RemoverAluno exclui um aluno
func (s *AlunoService) RemoverAluno(id uint) error {
	// Buscar aluno
	aluno, err := s.alunoRepo.FindByID(id)
	if err != nil {
		return errors.New("aluno não encontrado")
	}

	// Verificar se aluno tem inscrições ativas
	inscricoes, err := s.inscricaoRepo.FindByAluno(id)
	if err != nil {
		return err
	}

	if len(inscricoes) > 0 {
		return errors.New("não é possível remover aluno com inscrições ativas")
	}

	return s.alunoRepo.Delete(aluno)
}
