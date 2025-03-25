package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"tvtec/models" // ajuste o caminho conforme sua estrutura de pastas
)

// CursoService gerencia as operações para a entidade Curso.
type CursoService struct {
	db *gorm.DB
}

// NewCursoService cria uma nova instância do serviço com o DB injetado.
func NewCursoService(db *gorm.DB) *CursoService {
	return &CursoService{db: db}
}

// AddCurso adiciona um novo curso ao banco de dados.
func (s *CursoService) AddCurso(curso *models.Curso) (*models.Curso, error) {
	if curso == nil {
		return nil, errors.New("curso é nil")
	}
	if err := s.db.Create(curso).Error; err != nil {
		return nil, err
	}
	return curso, nil
}

// GetAllCursos retorna todos os cursos cadastrados.
func (s *CursoService) GetAllCursos() ([]models.Curso, error) {
	var cursos []models.Curso
	if err := s.db.Find(&cursos).Error; err != nil {
		return nil, err
	}
	return cursos, nil
}

// GetCurso busca um curso pelo ID.
func (s *CursoService) GetCurso(id uint) (*models.Curso, error) {
	var curso models.Curso
	if err := s.db.First(&curso, id).Error; err != nil {
		return nil, err
	}
	return &curso, nil
}

// UpdateCurso atualiza os dados de um curso existente.
func (s *CursoService) UpdateCurso(id uint, novoCurso *models.Curso) (*models.Curso, error) {
	var cursoExistente models.Curso
	if err := s.db.First(&cursoExistente, id).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("Curso não encontrado com o id %d", id))
	}

	// Atualiza os campos desejados usando o campo único "Data"
	cursoExistente.Nome = novoCurso.Nome
	cursoExistente.Professor = novoCurso.Professor
	cursoExistente.Data = novoCurso.Data
	cursoExistente.CargaHoraria = novoCurso.CargaHoraria
	cursoExistente.Certificado = novoCurso.Certificado

	if err := s.db.Save(&cursoExistente).Error; err != nil {
		return nil, err
	}
	return &cursoExistente, nil
}

// DeleteCurso remove um curso do banco de dados.
func (s *CursoService) DeleteCurso(id uint) error {
	var curso models.Curso
	if err := s.db.First(&curso, id).Error; err != nil {
		return err
	}
	return s.db.Delete(&curso).Error
}
