package repository

import (
	"gorm.io/gorm"
	"tvtec/models" // ajuste o import para o caminho correto dos seus modelos
)

// AlunoRepository gerencia as operações de persistência para a entidade Aluno.
type AlunoRepository struct {
	db *gorm.DB
}

// NewAlunoRepository cria uma nova instância do repositório com o DB injetado.
func NewAlunoRepository(db *gorm.DB) *AlunoRepository {
	return &AlunoRepository{db: db}
}

// Save persiste o aluno no banco de dados.
func (r *AlunoRepository) Save(aluno *models.Aluno) error {
	return r.db.Save(aluno).Error
}

// FindAll retorna todos os alunos registrados.
func (r *AlunoRepository) FindAll() ([]models.Aluno, error) {
	var alunos []models.Aluno
	err := r.db.Find(&alunos).Error
	return alunos, err
}

// FindByID busca um aluno pelo ID.
func (r *AlunoRepository) FindByID(id uint) (*models.Aluno, error) {
	var aluno models.Aluno
	if err := r.db.First(&aluno, id).Error; err != nil {
		return nil, err
	}
	return &aluno, nil
}

// Delete remove o aluno do banco de dados.
func (r *AlunoRepository) Delete(aluno *models.Aluno) error {
	return r.db.Delete(aluno).Error
}

// FindByEmail busca um aluno pelo email.
func (r *AlunoRepository) FindByEmail(email string) (*models.Aluno, error) {
	var aluno models.Aluno
	if err := r.db.Where("email = ?", email).First(&aluno).Error; err != nil {
		return nil, err
	}
	return &aluno, nil
}

// FindByTelefone busca um aluno pelo telefone.
func (r *AlunoRepository) FindByTelefone(telefone string) (*models.Aluno, error) {
	var aluno models.Aluno
	if err := r.db.Where("telefone = ?", telefone).First(&aluno).Error; err != nil {
		return nil, err
	}
	return &aluno, nil
}
