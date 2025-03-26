package service

import (
	"errors"
	"time"

	"tvtec/models"
	"tvtec/repository"
)

// CursoService contém a lógica de negócio para manipulação de cursos
type CursoService struct {
	repo *repository.CursoRepository
}

// NewCursoService cria uma nova instância do serviço de cursos
func NewCursoService(repo *repository.CursoRepository) *CursoService {
	return &CursoService{repo: repo}
}

// CriarCurso realiza a criação de um novo curso com validações de negócio
func (s *CursoService) CriarCurso(curso *models.Curso) error {
	// Validações de negócio
	if curso.Nome == "" {
		return errors.New("nome do curso é obrigatório")
	}

	if curso.Professor == "" {
		return errors.New("professor do curso é obrigatório")
	}

	if curso.CargaHoraria <= 0 {
		return errors.New("carga horária deve ser maior que zero")
	}

	if curso.VagasTotais <= 0 {
		return errors.New("número de vagas deve ser maior que zero")
	}

	// Inicializa vagas preenchidas como zero
	curso.VagasPreenchidas = 0

	// Define data de criação como momento atual se não informada
	if curso.Data.IsZero() {
		curso.Data = models.CustomTime{Time: time.Now()}
	}

	// Persiste o curso
	return s.repo.Save(curso)
}

// ObterCursoPorID busca um curso específico
func (s *CursoService) ObterCursoPorID(id uint) (*models.Curso, error) {
	curso, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("curso não encontrado")
	}
	return curso, nil
}

// ListarCursos recupera todos os cursos
func (s *CursoService) ListarCursos() ([]models.Curso, error) {
	return s.repo.FindAll()
}

// AtualizarCurso atualiza as informações de um curso existente
func (s *CursoService) AtualizarCurso(curso *models.Curso) error {
	// Verifica se o curso existe
	existente, err := s.repo.FindByID(curso.ID)
	if err != nil {
		return errors.New("curso não encontrado para atualização")
	}

	// Validações de negócio
	if curso.Nome == "" {
		curso.Nome = existente.Nome
	}

	if curso.Professor == "" {
		curso.Professor = existente.Professor
	}

	// Mantém o controle de vagas original
	curso.VagasPreenchidas = existente.VagasPreenchidas

	// Persiste a atualização
	return s.repo.Save(curso)
}

// RemoverCurso exclui um curso
func (s *CursoService) RemoverCurso(id uint) error {
	// Busca o curso primeiro para garantir que existe
	curso, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("curso não encontrado")
	}

	// Verifica se existem vagas preenchidas
	if curso.VagasPreenchidas > 0 {
		return errors.New("não é possível remover curso com inscrições ativas")
	}

	// Remove o curso
	return s.repo.Delete(curso)
}

// VerificarDisponibilidadeVagas verifica se ainda há vagas disponíveis
func (s *CursoService) VerificarDisponibilidadeVagas(cursoID uint) (bool, error) {
	curso, err := s.repo.FindByID(cursoID)
	if err != nil {
		return false, errors.New("curso não encontrado")
	}

	return curso.VagasPreenchidas < curso.VagasTotais, nil
}
