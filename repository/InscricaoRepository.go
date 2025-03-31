package repository

import (
	"errors"
	"tvtec/models"

	"gorm.io/gorm"
)

type InscricaoRepository interface {
	FindAll() ([]models.Inscricao, error)
	FindByID(id uint) (*models.Inscricao, error)
	FindAllWithDetails() ([]models.Inscricao, error)
	FindByIDWithDetails(id uint) (*models.Inscricao, error)
	FindByAluno(alunoID uint) ([]models.Inscricao, error)
	FindByCurso(cursoID uint) ([]models.Inscricao, error)
	FindByAlunoWithDetails(alunoID uint) ([]models.Inscricao, error)
	FindByCursoWithDetails(cursoID uint) ([]models.Inscricao, error)
	Save(inscricao *models.Inscricao) error
	Delete(id uint) error
	CountByAluno(alunoID uint) (int64, error)
	CountByCurso(cursoID uint) (int64, error)
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

func (r *inscricaoRepository) FindAllWithDetails() ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Preload("Aluno").Preload("Curso").Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByIDWithDetails(id uint) (*models.Inscricao, error) {
	var inscricao models.Inscricao
	result := r.db.Preload("Aluno").Preload("Curso").First(&inscricao, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("inscrição não encontrada")
		}
		return nil, result.Error
	}
	return &inscricao, nil
}

func (r *inscricaoRepository) FindByAluno(alunoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("aluno_id = ?", alunoID).Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByCurso(cursoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("curso_id = ?", cursoID).Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByAlunoWithDetails(alunoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("aluno_id = ?", alunoID).Preload("Aluno").Preload("Curso").Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) FindByCursoWithDetails(cursoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao
	result := r.db.Where("curso_id = ?", cursoID).Preload("Aluno").Preload("Curso").Find(&inscricoes)
	return inscricoes, result.Error
}

func (r *inscricaoRepository) Save(inscricao *models.Inscricao) error {
	if inscricao.ID == 0 {
		// Verifica se já existe uma inscrição deste aluno neste curso
		var count int64
		r.db.Model(&models.Inscricao{}).
			Where("aluno_id = ? AND curso_id = ?", inscricao.AlunoID, inscricao.CursoID).
			Count(&count)

		if count > 0 {
			return errors.New("aluno já está inscrito neste curso")
		}

		// Iniciar transação para garantir consistência
		return r.db.Transaction(func(tx *gorm.DB) error {
			// Verificar disponibilidade de vagas no curso
			var curso models.Curso
			if err := tx.First(&curso, inscricao.CursoID).Error; err != nil {
				return err
			}

			// Verificar se há vagas disponíveis
			if curso.VagasPreenchidas >= curso.VagasTotais {
				return errors.New("não há vagas disponíveis para este curso")
			}

			// Salvar inscrição
			if err := tx.Create(inscricao).Error; err != nil {
				return err
			}

			// Atualizar vagas preenchidas
			curso.VagasPreenchidas++
			if err := tx.Save(&curso).Error; err != nil {
				return err
			}

			return nil
		})
	} else {
		// Atualizar inscrição existente
		return r.db.Save(inscricao).Error
	}
}

func (r *inscricaoRepository) Delete(id uint) error {
	// Buscar a inscrição para obter o curso_id
	var inscricao models.Inscricao
	if err := r.db.First(&inscricao, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inscrição não encontrada")
		}
		return err
	}

	// Iniciar transação para garantir consistência
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Recuperar o curso para atualizar vagas
		var curso models.Curso
		if err := tx.First(&curso, inscricao.CursoID).Error; err != nil {
			return err
		}

		// Remover inscrição
		if err := tx.Delete(&inscricao).Error; err != nil {
			return err
		}

		// Atualizar vagas preenchidas (evitar vagas negativas)
		if curso.VagasPreenchidas > 0 {
			curso.VagasPreenchidas--
			if err := tx.Save(&curso).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *inscricaoRepository) CountByAluno(alunoID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Inscricao{}).Where("aluno_id = ?", alunoID).Count(&count)
	return count, result.Error
}

func (r *inscricaoRepository) CountByCurso(cursoID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Inscricao{}).Where("curso_id = ?", cursoID).Count(&count)
	return count, result.Error
}
