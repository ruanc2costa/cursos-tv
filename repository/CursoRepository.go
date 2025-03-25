package repository

import (
	"gorm.io/gorm"
	"tvtec/models" // ajuste o import para o caminho correto dos seus modelos
)

// CursoRepository gerencia as operações de persistência para a entidade Curso.
type CursoRepository struct {
	db *gorm.DB
}

// NewCursoRepository cria uma nova instância do repositório com o DB injetado.
func NewCursoRepository(db *gorm.DB) *CursoRepository {
	return &CursoRepository{db: db}
}

// Save persiste o curso no banco de dados.
func (r *CursoRepository) Save(curso *models.Curso) error {
	return r.db.Save(curso).Error
}

// FindAll retorna todos os cursos registrados.
func (r *CursoRepository) FindAll() ([]models.Curso, error) {
	var cursos []models.Curso
	err := r.db.Find(&cursos).Error
	return cursos, err
}

// FindByID busca um curso pelo ID.
func (r *CursoRepository) FindByID(id uint) (*models.Curso, error) {
	var curso models.Curso
	if err := r.db.First(&curso, id).Error; err != nil {
		return nil, err
	}
	return &curso, nil
}

// Delete remove o curso do banco de dados.
func (r *CursoRepository) Delete(curso *models.Curso) error {
	return r.db.Delete(curso).Error
}
