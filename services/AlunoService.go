package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"your_project/models" // ajuste o caminho conforme sua estrutura de pastas
)

// AlunoService gerencia as operações com a entidade Aluno.
type AlunoService struct {
	db *gorm.DB
}

// NewAlunoService cria uma nova instância do serviço com o DB injetado.
func NewAlunoService(db *gorm.DB) *AlunoService {
	return &AlunoService{db: db}
}

// AddAluno cria um aluno e realiza a inscrição em um curso.
// Se o aluno já existir (verificado pelo email):
//   - Se ele já estiver inscrito no curso, nada é feito.
//   - Caso contrário, cria uma nova inscrição para o curso.
//
// Se o aluno não existir, cria o aluno e a inscrição.
func (s *AlunoService) AddAluno(aluno *models.Aluno) (*models.Aluno, error) {
	if aluno == nil {
		return nil, errors.New("aluno é nil")
	}

	// Verifica se pelo menos um curso foi informado.
	if len(aluno.Cursos) == 0 {
		return nil, errors.New("deve ser informado pelo menos um curso")
	}
	selectedCurso := aluno.Cursos[0]

	// Tenta encontrar um aluno existente pelo email e carrega os cursos já associados.
	var existingAluno models.Aluno
	err := s.db.Where("email = ?", aluno.Email).Preload("Cursos").First(&existingAluno).Error
	if err == nil {
		// Aluno já existe. Verifica se já está inscrito no curso.
		var alreadyEnrolled bool
		for _, c := range existingAluno.Cursos {
			if c.ID == selectedCurso.ID {
				alreadyEnrolled = true
				break
			}
		}

		if !alreadyEnrolled {
			// Cria uma nova inscrição (registro na tabela Inscricao).
			newInscricao := models.Inscricao{
				AlunoID:       existingAluno.ID,
				CursoID:       selectedCurso.ID,
				DataInscricao: time.Now(),
			}
			if err := s.db.Create(&newInscricao).Error; err != nil {
				return nil, err
			}
		}
		return &existingAluno, nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Erro diferente de "registro não encontrado"
		return nil, err
	}

	// Se o aluno não existe, cria o novo aluno.
	if err := s.db.Create(aluno).Error; err != nil {
		return nil, err
	}

	// Após criar o aluno, cria a inscrição para o curso selecionado.
	newInscricao := models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       selectedCurso.ID,
		DataInscricao: time.Now(),
	}
	if err := s.db.Create(&newInscricao).Error; err != nil {
		return nil, err
	}

	return aluno, nil
}
