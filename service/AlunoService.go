package service

import (
	"errors"
	"time"

	"tvtec/models"
	"tvtec/repository"
)

// Interface para o serviço de alunos
type AlunoService interface {
	ListarAlunos() ([]models.Aluno, error)
	ObterAlunoPorID(id uint) (*models.Aluno, error)
	CriarAluno(aluno *models.Aluno) error
	AtualizarAluno(aluno *models.Aluno) error
	RemoverAluno(id uint) error
	CadastrarAlunoEInscrever(aluno *models.Aluno, inscricao *models.Inscricao) error
	AdicionarAlunoCurso(alunoID, cursoID uint) error
	CriarInscricaoDetalhada(inscricao *models.Inscricao) error
	ListarInscricoesAluno(alunoID uint) ([]models.Inscricao, error)
}

// Implementação do serviço de alunos
type alunoServiceImpl struct {
	alunoRepo     repository.AlunoRepository
	cursoRepo     repository.CursoRepository
	inscricaoRepo repository.InscricaoRepository
}

// Função construtora para o serviço de alunos
func NewAlunoService(alunoRepo repository.AlunoRepository, cursoRepo repository.CursoRepository, inscricaoRepo repository.InscricaoRepository) AlunoService {
	return &alunoServiceImpl{
		alunoRepo:     alunoRepo,
		cursoRepo:     cursoRepo,
		inscricaoRepo: inscricaoRepo,
	}
}

// Métodos de implementação

func (s *alunoServiceImpl) ListarAlunos() ([]models.Aluno, error) {
	return s.alunoRepo.FindAll()
}

func (s *alunoServiceImpl) ObterAlunoPorID(id uint) (*models.Aluno, error) {
	aluno, err := s.alunoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("aluno não encontrado")
	}
	return aluno, nil
}

func (s *alunoServiceImpl) CriarAluno(aluno *models.Aluno) error {
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

func (s *alunoServiceImpl) AtualizarAluno(aluno *models.Aluno) error {
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

func (s *alunoServiceImpl) RemoverAluno(id uint) error {
	// Verificar se o aluno existe (sem armazenar o resultado)
	if _, err := s.alunoRepo.FindByID(id); err != nil {
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

	return s.alunoRepo.Delete(id)
}

func (s *alunoServiceImpl) CadastrarAlunoEInscrever(aluno *models.Aluno, inscricao *models.Inscricao) error {
	// Tenta encontrar o aluno pelo email
	alunoExistente, _ := s.alunoRepo.FindByEmail(aluno.Email)
	if alunoExistente != nil {
		// O aluno já existe: usa o registro existente
		aluno = alunoExistente
	} else {
		// Se não encontrou pelo email, pode verificar pelo CPF, se necessário
		existenteCPF, _ := s.alunoRepo.FindByCPF(aluno.CPF)
		if existenteCPF != nil {
			aluno = existenteCPF
		} else {
			// Se o aluno não existe, salva o novo aluno
			if err := s.alunoRepo.Save(aluno); err != nil {
				return err
			}
		}
	}

	// Verifica se o curso existe
	curso, err := s.cursoRepo.FindByID(inscricao.CursoID)
	if err != nil {
		return errors.New("curso não encontrado")
	}

	// Verifica disponibilidade de vagas
	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis para este curso")
	}

	// Opcional: Verifica se já existe inscrição para este aluno e curso
	inscricaoExistente, _ := s.inscricaoRepo.FindByAlunoECurso(aluno.ID, curso.ID)
	if inscricaoExistente != nil {
		return errors.New("já existe uma inscrição para este curso")
	}

	// Associa o aluno à inscrição
	inscricao.AlunoID = aluno.ID

	// Define valor padrão para campos opcionais
	if inscricao.EhPCD == "" {
		inscricao.EhPCD = "não"
	}

	// Salva a inscrição
	if err := s.inscricaoRepo.Save(inscricao); err != nil {
		return err
	}

	return nil
}

func (s *alunoServiceImpl) AdicionarAlunoCurso(alunoID, cursoID uint) error {
	// Verificar se o aluno existe
	if _, err := s.alunoRepo.FindByID(alunoID); err != nil {
		return errors.New("aluno não encontrado")
	}

	// Verificar se o curso existe e tem vagas disponíveis
	curso, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return errors.New("curso não encontrado")
	}

	// Verificar disponibilidade de vagas - aqui usamos a variável curso
	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis para este curso")
	}

	// Verificar se o aluno já está inscrito neste curso
	inscricoes, err := s.inscricaoRepo.FindByAluno(alunoID)
	if err != nil {
		return err
	}

	for _, inscricao := range inscricoes {
		if inscricao.CursoID == cursoID {
			return errors.New("aluno já está inscrito neste curso")
		}
	}

	// Criar uma nova inscrição
	inscricao := &models.Inscricao{
		AlunoID:       alunoID,
		CursoID:       cursoID,
		DataInscricao: time.Now(),
		EhPCD:         "não", // Valor padrão
	}

	// Salvar a inscrição
	return s.inscricaoRepo.Save(inscricao)
}

func (s *alunoServiceImpl) CriarInscricaoDetalhada(inscricao *models.Inscricao) error {
	// Verificar se o aluno existe
	_, err := s.alunoRepo.FindByID(inscricao.AlunoID)
	if err != nil {
		return errors.New("aluno não encontrado")
	}

	// Verificar se o curso existe
	curso, err := s.cursoRepo.FindByID(inscricao.CursoID)
	if err != nil {
		return errors.New("curso não encontrado")
	}

	// Verificar disponibilidade de vagas
	if curso.VagasPreenchidas >= curso.VagasTotais {
		return errors.New("não há vagas disponíveis para este curso")
	}

	// Verificar se o aluno já está inscrito neste curso
	inscricoes, err := s.inscricaoRepo.FindByAluno(inscricao.AlunoID)
	if err != nil {
		return err
	}

	for _, existente := range inscricoes {
		if existente.CursoID == inscricao.CursoID {
			return errors.New("aluno já está inscrito neste curso")
		}
	}

	// Definir valores padrão para campos opcionais se não forem informados
	if inscricao.EhPCD == "" {
		inscricao.EhPCD = "não"
	}

	// Definir a data de inscrição como a data atual se não for informada
	if inscricao.DataInscricao.IsZero() {
		inscricao.DataInscricao = time.Now()
	}

	// Salvar a inscrição
	return s.inscricaoRepo.Save(inscricao)
}

func (s *alunoServiceImpl) ListarInscricoesAluno(alunoID uint) ([]models.Inscricao, error) {
	// Verificar se o aluno existe
	_, err := s.alunoRepo.FindByID(alunoID)
	if err != nil {
		return nil, errors.New("aluno não encontrado")
	}

	// Buscar inscrições do aluno com detalhes de cursos
	return s.inscricaoRepo.FindByAlunoWithDetails(alunoID)
}
