package repository

import (
	"errors"
	"tvtec/models"

	"gorm.io/gorm"
)

type AlunoRepository interface {
	FindAll() ([]models.Aluno, error)
	FindByID(id uint) (*models.Aluno, error)
	Save(aluno *models.Aluno) error
	Update(aluno *models.Aluno) error
	Delete(id uint) error
	FindByCPF(cpf string) (*models.Aluno, error)
	FindByEmail(email string) (*models.Aluno, error)
}

type alunoRepository struct {
	db *gorm.DB
}

func NewAlunoRepository(db *gorm.DB) AlunoRepository {
	return &alunoRepository{db: db}
}

func (r *alunoRepository) FindAll() ([]models.Aluno, error) {
	var alunos []models.Aluno
	result := r.db.Find(&alunos)
	return alunos, result.Error
}

func (r *alunoRepository) FindByID(id uint) (*models.Aluno, error) {
	var aluno models.Aluno
	result := r.db.First(&aluno, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("aluno não encontrado")
		}
		return nil, result.Error
	}
	return &aluno, nil
}

func (r *alunoRepository) Save(aluno *models.Aluno) error {
	return r.db.Create(aluno).Error
}

func (r *alunoRepository) Update(aluno *models.Aluno) error {
	return r.db.Save(aluno).Error
}

func (r *alunoRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Aluno{}, id)
	if result.RowsAffected == 0 {
		return errors.New("aluno não encontrado")
	}
	return result.Error
}

func (r *alunoRepository) FindByCPF(cpf string) (*models.Aluno, error) {
	var aluno models.Aluno
	result := r.db.Where("cpf = ?", cpf).First(&aluno)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("aluno não encontrado")
		}
		return nil, result.Error
	}
	return &aluno, nil
}
func (r *alunoRepository) FindByEmail(email string) (*models.Aluno, error) {
	var aluno models.Aluno
	result := r.db.Where("email = ?", email).First(&aluno)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("aluno não encontrado")
		}
		return nil, result.Error
	}
	return &aluno, nil
}
