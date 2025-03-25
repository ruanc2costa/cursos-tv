package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tvtec/models" // ajuste o import para o caminho correto do pacote models
)

// AlunoService gerencia as operações com a entidade Aluno.
type AlunoService struct {
	db *gorm.DB
}

// NewAlunoService cria uma nova instância do serviço com o DB injetado.
func NewAlunoService(db *gorm.DB) *AlunoService {
	return &AlunoService{db: db}
}

// AddAluno adiciona um novo aluno. Se o aluno for nulo, retorna um erro.
func (s *AlunoService) AddAluno(aluno *models.Aluno) (*models.Aluno, error) {
	if aluno == nil {
		return nil, errors.New("aluno é nil")
	}

	result := s.db.Create(aluno)
	if result.Error != nil {
		return nil, result.Error
	}
	return aluno, nil
}

// GetAllAlunos retorna todos os alunos registrados.
func (s *AlunoService) GetAllAlunos() ([]models.Aluno, error) {
	var alunos []models.Aluno
	result := s.db.Find(&alunos)
	return alunos, result.Error
}

// GetAluno busca um aluno pelo ID.
func (s *AlunoService) GetAluno(id uint) (*models.Aluno, error) {
	var aluno models.Aluno
	result := s.db.First(&aluno, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New(fmt.Sprintf("Aluno não encontrado com o id %d", id))
	}
	return &aluno, result.Error
}

// UpdateAluno atualiza os dados de um aluno existente.
func (s *AlunoService) UpdateAluno(id uint, novoAluno *models.Aluno) (*models.Aluno, error) {
	var alunoExistente models.Aluno
	result := s.db.First(&alunoExistente, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New(fmt.Sprintf("Aluno não encontrado com o id %d", id))
	}

	// Atualiza os campos desejados
	alunoExistente.Nome = novoAluno.Nome
	alunoExistente.Sobrenome = novoAluno.Sobrenome
	alunoExistente.Sexo = novoAluno.Sexo
	alunoExistente.DataNascto = novoAluno.DataNascto

	result = s.db.Save(&alunoExistente)
	return &alunoExistente, result.Error
}

// DeleteAluno remove um aluno com base no ID.
func (s *AlunoService) DeleteAluno(id uint) error {
	var aluno models.Aluno
	result := s.db.First(&aluno, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("Aluno não encontrado com o id %d", id))
	}
	result = s.db.Delete(&aluno)
	return result.Error
}
