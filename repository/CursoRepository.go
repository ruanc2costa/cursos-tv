package repository

import (
	"errors"
	"tvtec/models"

	"gorm.io/gorm"
)

type CursoRepository interface {
	FindAll() ([]models.Curso, error)
	FindByID(id uint) (*models.Curso, error)
	Save(curso *models.Curso) error
	Update(curso *models.Curso) error
	Delete(id uint) error
	IncrementarVagasPreenchidas(cursoID uint) error
	DecrementarVagasPreenchidas(cursoID uint) error
}

type cursoRepository struct {
	db *gorm.DB
}

func NewCursoRepository(db *gorm.DB) CursoRepository {
	return &cursoRepository{db: db}
}

func (r *cursoRepository) FindAll() ([]models.Curso, error) {
	var cursos []models.Curso
	result := r.db.Find(&cursos)
	return cursos, result.Error
}

func (r *cursoRepository) FindByID(id uint) (*models.Curso, error) {
	var curso models.Curso
	result := r.db.First(&curso, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("curso não encontrado")
		}
		return nil, result.Error
	}
	return &curso, nil
}

func (r *cursoRepository) Save(curso *models.Curso) error {
	// Garantir que vagas preenchidas começa com zero
	curso.VagasPreenchidas = 0
	return r.db.Create(curso).Error
}

func (r *cursoRepository) Update(curso *models.Curso) error {
	return r.db.Save(curso).Error
}

func (r *cursoRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Curso{}, id)
	if result.RowsAffected == 0 {
		return errors.New("curso não encontrado")
	}
	return result.Error
}

func (r *cursoRepository) IncrementarVagasPreenchidas(cursoID uint) error {
	// Usar uma transação para evitar condições de corrida
	return r.db.Transaction(func(tx *gorm.DB) error {
		var curso models.Curso
		// Primeiro, obter o curso atual com bloqueio para atualização
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&curso, cursoID).Error; err != nil {
			return err
		}

		// Verificar se há vagas disponíveis
		if curso.VagasPreenchidas >= curso.VagasTotais {
			return errors.New("não há vagas disponíveis no curso")
		}

		// Incrementar e salvar
		curso.VagasPreenchidas++
		return tx.Save(&curso).Error
	})
}

func (r *cursoRepository) DecrementarVagasPreenchidas(cursoID uint) error {
	// Usar uma transação para evitar condições de corrida
	return r.db.Transaction(func(tx *gorm.DB) error {
		var curso models.Curso
		// Primeiro, obter o curso atual com bloqueio para atualização
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&curso, cursoID).Error; err != nil {
			return err
		}

		// Verificar se há vagas preenchidas para decrementar
		if curso.VagasPreenchidas <= 0 {
			return errors.New("não há vagas preenchidas para decrementar")
		}

		// Decrementar e salvar
		curso.VagasPreenchidas--
		return tx.Save(&curso).Error
	})
}
