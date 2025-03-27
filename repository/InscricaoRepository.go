// repository/inscricao_repository.go
package repository

import (
	"errors"
	"time"
	"tvtec/models"

	"gorm.io/gorm"
)

type InscricaoRepository interface {
	FindAll() ([]models.Inscricao, error)
	FindByID(id uint) (*models.Inscricao, error)
	FindByCursoID(cursoID uint) ([]models.Inscricao, error)
	FindByAlunoID(alunoID uint) ([]models.Inscricao, error)
	FindByAlunoECurso(alunoID, cursoID uint) (*models.Inscricao, error)
	Save(inscricao *models.Inscricao) error
	Delete(id uint) error
}

type inscricaoRepository struct {
	db *gorm.DB
}

func NewInscricaoRepository(db *gorm.DB) InscricaoRepository {
	return &inscricaoRepository{db: db}
}

func (r *inscricaoRepository) FindAll() ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByID(id uint) (*models.Inscricao, error) {
	var inscricao models.Inscricao
	result := r.db.First(&inscricao, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("inscrição não encontrada")
		}
		return nil, result.Error
	}
	return &inscricao, nil
}

func (r *inscricaoRepository) FindByCursoID(cursoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("curso_id = ?", cursoID).Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByAlunoID(alunoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("aluno_id = ?", alunoID).Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByAlunoECurso(alunoID, cursoID uint) (*models.Inscricao, error) {
	var inscricao models.Inscricao
	result := r.db.Where("aluno_id = ? AND curso_id = ?", alunoID, cursoID).First(&inscricao)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("inscrição não encontrada")
		}
		return nil, result.Error
	}
	return &inscricao, nil
}

func (r *inscricaoRepository) Save(inscricao *models.Inscricao) error {
	// Garantir que a data de inscrição seja definida
	if inscricao.DataInscricao.IsZero() {
		inscricao.DataInscricao = time.Now()
	}
	return r.db.Create(inscricao).Error
}

func (r *inscricaoRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Inscricao{}, id)
	if result.RowsAffected == 0 {
		return errors.New("inscrição não encontrada")
	}
	return result.Error
}
