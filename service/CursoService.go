package service

import (
	"tvtec/models"
	"tvtec/repository"
)

type CursoService interface {
	ListarCursos() ([]models.Curso, error)
	ObterCursoPorID(id uint) (*models.Curso, error)
	CriarCurso(curso *models.Curso) error
	AtualizarCurso(curso *models.Curso) error
	RemoverCurso(id uint) error
	VerificarDisponibilidadeVagas(id uint) (int32, error)
	ListarInscricoesCurso(cursoID uint) ([]models.Inscricao, error)
}

type cursoService struct {
	cursoRepo     repository.CursoRepository
	inscricaoRepo repository.InscricaoRepository
}

func NewCursoService(
	cursoRepo repository.CursoRepository,
	inscricaoRepo repository.InscricaoRepository,
) CursoService {
	return &cursoService{
		cursoRepo:     cursoRepo,
		inscricaoRepo: inscricaoRepo,
	}
}

func (s *cursoService) ListarCursos() ([]models.Curso, error) {
	return s.cursoRepo.FindAll()
}

func (s *cursoService) ObterCursoPorID(id uint) (*models.Curso, error) {
	return s.cursoRepo.FindByID(id)
}

func (s *cursoService) CriarCurso(curso *models.Curso) error {
	return s.cursoRepo.Save(curso)
}

func (s *cursoService) AtualizarCurso(curso *models.Curso) error {
	return s.cursoRepo.Update(curso)
}

func (s *cursoService) RemoverCurso(id uint) error {
	// Primeiro verificamos se o curso existe
	curso, err := s.cursoRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verificamos se existem inscrições para este curso usando o método FindByCurso
	inscricoes, err := s.inscricaoRepo.FindByCurso(id)
	if err != nil {
		return err
	}

	// Se existem inscrições, precisamos removê-las primeiro
	for _, inscricao := range inscricoes {
		if err := s.inscricaoRepo.Delete(inscricao.ID); err != nil {
			return err
		}
	}

	// Finalmente remover o curso
	return s.cursoRepo.Delete(curso.ID)
}

func (s *cursoService) VerificarDisponibilidadeVagas(id uint) (int32, error) {
	curso, err := s.cursoRepo.FindByID(id)
	if err != nil {
		return 0, err
	}

	// Retorna o número de vagas disponíveis
	return curso.VagasTotais - curso.VagasPreenchidas, nil
}

func (s *cursoService) ListarInscricoesCurso(cursoID uint) ([]models.Inscricao, error) {
	// Verificar se o curso existe
	_, err := s.cursoRepo.FindByID(cursoID)
	if err != nil {
		return nil, err
	}

	// Retornar inscrições do curso usando o método FindByCurso
	return s.inscricaoRepo.FindByCurso(cursoID)
}
