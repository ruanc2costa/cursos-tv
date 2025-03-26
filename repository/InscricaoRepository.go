package repository

import (
	"errors"
	"gorm.io/gorm"
	"tvtec/models"
)

// InscricaoRepository gerencia as operações de persistência para a entidade Inscrição.
type InscricaoRepository struct {
	db *gorm.DB
}

// NewInscricaoRepository cria uma nova instância do repositório com o DB injetado.
func NewInscricaoRepository(db *gorm.DB) *InscricaoRepository {
	return &InscricaoRepository{db: db}
}

// Save persiste a inscrição no banco de dados.
func (r *InscricaoRepository) Save(inscricao *models.Inscricao) error {
	// Validações antes de salvar
	if inscricao.AlunoID == 0 || inscricao.CursoID == 0 {
		return errors.New("AlunoID e CursoID são obrigatórios")
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

		// Verificar se aluno já está inscrito
		var existingInscricao models.Inscricao
		if err := tx.Where("aluno_id = ? AND curso_id = ?",
			inscricao.AlunoID, inscricao.CursoID).First(&existingInscricao).Error; err == nil {
			return errors.New("aluno já está inscrito neste curso")
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
}

// FindByID busca uma inscrição pelo ID.
func (r *InscricaoRepository) FindByID(id uint) (*models.Inscricao, error) {
	var inscricao models.Inscricao

	// Carregar inscrição com dados relacionados de Aluno e Curso
	err := r.db.Preload("Aluno").Preload("Curso").First(&inscricao, id).Error
	if err != nil {
		return nil, err
	}

	return &inscricao, nil
}

// FindByAluno busca todas as inscrições de um aluno.
func (r *InscricaoRepository) FindByAluno(alunoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao

	err := r.db.Where("aluno_id = ?", alunoID).
		Preload("Curso").
		Find(&inscricoes).Error

	return inscricoes, err
}

// FindByCurso busca todas as inscrições de um curso.
func (r *InscricaoRepository) FindByCurso(cursoID uint) ([]models.Inscricao, error) {
	var inscricoes []models.Inscricao

	err := r.db.Where("curso_id = ?", cursoID).
		Preload("Aluno").
		Find(&inscricoes).Error

	return inscricoes, err
}

// Delete remove a inscrição do banco de dados.
func (r *InscricaoRepository) Delete(inscricao *models.Inscricao) error {
	// Iniciar transação para garantir consistência
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Recuperar o curso para atualizar vagas
		var curso models.Curso
		if err := tx.First(&curso, inscricao.CursoID).Error; err != nil {
			return err
		}

		// Remover inscrição
		if err := tx.Delete(inscricao).Error; err != nil {
			return err
		}

		// Atualizar vagas preenchidas
		curso.VagasPreenchidas--
		if err := tx.Save(&curso).Error; err != nil {
			return err
		}

		return nil
	})
}

// CountByAluno conta o número de inscrições de um aluno.
func (r *InscricaoRepository) CountByAluno(alunoID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Inscricao{}).
		Where("aluno_id = ?", alunoID).
		Count(&count).Error

	return count, err
}

// CountByCurso conta o número de inscrições em um curso.
func (r *InscricaoRepository) CountByCurso(cursoID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Inscricao{}).
		Where("curso_id = ?", cursoID).
		Count(&count).Error

	return count, err
}
