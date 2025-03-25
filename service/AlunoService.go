package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"tvtec/models" // ajuste o caminho conforme sua estrutura de pastas
)

// AlunoService gerencia as operações com a entidade Aluno.
type AlunoService struct {
	db *gorm.DB
}

// NewAlunoService cria uma nova instância do serviço.
func NewAlunoService(db *gorm.DB) *AlunoService {
	return &AlunoService{db: db}
}

// AddAluno cria um aluno e realiza a inscrição em um curso.
// Se o aluno já existir (verificado pelo email) e não estiver inscrito no curso informado,
// cria uma nova inscrição. Se o aluno não existir, cria o aluno e realiza a inscrição.
// Antes de tudo, verifica se o CPF já existe e, em caso afirmativo, retorna erro de conflito (HTTP 409).
// Além disso, verifica se o curso possui vagas disponíveis e, se sim, decrementa o número de vagas.
func (s *AlunoService) AddAluno(aluno *models.Aluno) (*models.Aluno, error) {
	// Valida se o objeto aluno não é nulo
	if aluno == nil {
		return nil, errors.New("aluno é nil")
	}

	// Verifica se o CPF foi informado
	if aluno.CPF == "" {
		return nil, errors.New("CPF deve ser informado")
	}

	// Valida se pelo menos um curso foi informado (usaremos o primeiro curso do array)
	if len(aluno.Cursos) == 0 {
		return nil, errors.New("deve ser informado pelo menos um curso")
	}

	// Verifica se já existe um aluno com o mesmo CPF
	var alunoCPF models.Aluno
	err := s.db.Where("cpf = ?", aluno.CPF).First(&alunoCPF).Error
	if err == nil {
		// CPF duplicado: retorne um erro (a camada de controller pode mapear para HTTP 409)
		return nil, fmt.Errorf("aluno com CPF duplicado: %s", aluno.CPF)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Seleciona o curso informado (usamos o primeiro curso do array)
	selectedCurso := aluno.Cursos[0]

	// Carrega os dados atuais do curso para verificar as vagas
	var curso models.Curso
	if err := s.db.First(&curso, selectedCurso.ID).Error; err != nil {
		return nil, err
	}
	if curso.Vagas <= 0 {
		return nil, fmt.Errorf("não há vagas disponíveis para o curso %s", curso.Nome)
	}

	// Tenta encontrar um aluno existente pelo email e carrega os cursos associados
	var existingAluno models.Aluno
	err = s.db.Where("email = ?", aluno.Email).Preload("Cursos").First(&existingAluno).Error
	if err == nil {
		// Aluno já existe; verifica se já está inscrito no curso
		for _, c := range existingAluno.Cursos {
			if c.ID == selectedCurso.ID {
				// Já está inscrito, retorna o aluno existente
				return &existingAluno, nil
			}
		}
		// Se não estiver inscrito, cria uma nova inscrição para o curso
		newInscricao := models.Inscricao{
			AlunoID:       existingAluno.ID,
			CursoID:       selectedCurso.ID,
			DataInscricao: time.Now(),
		}
		if err := s.db.Create(&newInscricao).Error; err != nil {
			return nil, err
		}
		// Decrementa o número de vagas do curso e salva a alteração
		curso.Vagas--
		if err := s.db.Save(&curso).Error; err != nil {
			return nil, err
		}
		// Opcional: recarrega o aluno com os cursos atualizados
		if err := s.db.Preload("Cursos").First(&existingAluno, existingAluno.ID).Error; err != nil {
			return nil, err
		}
		return &existingAluno, nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Se o aluno não existir, cria o novo aluno
	if err := s.db.Create(aluno).Error; err != nil {
		return nil, err
	}

	// Após criar o aluno, realiza a inscrição no curso selecionado
	newInscricao := models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       selectedCurso.ID,
		DataInscricao: time.Now(),
	}
	if err := s.db.Create(&newInscricao).Error; err != nil {
		return nil, err
	}

	// Decrementa o número de vagas do curso e salva a alteração
	curso.Vagas--
	if err := s.db.Save(&curso).Error; err != nil {
		return nil, err
	}

	// Opcional: recarrega o aluno com os cursos atualizados
	if err := s.db.Preload("Cursos").First(aluno, aluno.ID).Error; err != nil {
		return nil, err
	}
	return aluno, nil
}

// GetAllAlunos retorna todos os alunos com seus cursos associados.
func (s *AlunoService) GetAllAlunos() ([]models.Aluno, error) {
	var alunos []models.Aluno
	if err := s.db.Preload("Cursos").Find(&alunos).Error; err != nil {
		return nil, err
	}
	return alunos, nil
}

// GetAluno busca um aluno pelo ID, carregando os cursos associados.
func (s *AlunoService) GetAluno(id uint) (*models.Aluno, error) {
	var aluno models.Aluno
	if err := s.db.Preload("Cursos").First(&aluno, id).Error; err != nil {
		return nil, err
	}
	return &aluno, nil
}

// DeleteAluno remove um aluno do banco de dados.
func (s *AlunoService) DeleteAluno(id uint) error {
	var aluno models.Aluno
	if err := s.db.First(&aluno, id).Error; err != nil {
		return err
	}
	return s.db.Delete(&aluno).Error
}
