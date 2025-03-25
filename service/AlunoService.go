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
// Além disso, verifica se há vagas disponíveis no curso e, se houver, decrementa o número de vagas.
func (s *AlunoService) AddAluno(aluno *models.Aluno) (*models.Aluno, error) {
	// Validação básica
	if aluno == nil {
		return nil, errors.New("aluno é nil")
	}
	if aluno.CPF == "" {
		return nil, errors.New("CPF deve ser informado")
	}
	if len(aluno.Cursos) == 0 {
		return nil, errors.New("deve ser informado pelo menos um curso")
	}

	// Verifica se já existe um aluno com o mesmo CPF
	var alunoCPF models.Aluno
	err := s.db.Where("cpf = ?", aluno.CPF).First(&alunoCPF).Error
	if err == nil {
		return nil, fmt.Errorf("Aluno com CPF duplicado: %s", aluno.CPF)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Seleciona o curso informado (usamos o primeiro curso do array)
	selectedCurso := aluno.Cursos[0]

	// Carrega os dados atuais do curso para verificar vagas
	var curso models.Curso
	if err := s.db.First(&curso, selectedCurso.ID).Error; err != nil {
		return nil, err
	}
	if curso.Vagas <= 0 {
		return nil, fmt.Errorf("não há vagas disponíveis para o curso %s", curso.Nome)
	}

	// Tenta encontrar um aluno existente pelo email (com preload de cursos)
	var existingAluno models.Aluno
	err = s.db.Where("email = ?", aluno.Email).Preload("Cursos").First(&existingAluno).Error
	if err == nil {
		// Aluno já existe; verifica se já está inscrito no curso
		for _, c := range existingAluno.Cursos {
			if c.ID == selectedCurso.ID {
				// Já está inscrito; não precisa criar nova inscrição.
				return &existingAluno, nil
			}
		}
		// Se não estiver inscrito, cria uma nova inscrição
		newInscricao := models.Inscricao{
			AlunoID:       existingAluno.ID,
			CursoID:       selectedCurso.ID,
			DataInscricao: time.Now(),
		}
		if err := s.db.Create(&newInscricao).Error; err != nil {
			return nil, err
		}
		// Decrementa as vagas do curso
		curso.Vagas--
		if err := s.db.Save(&curso).Error; err != nil {
			return nil, err
		}
		// Opcional: recarrega o aluno com os cursos atualizados.
		err = s.db.Preload("Cursos").First(&existingAluno, existingAluno.ID).Error
		return &existingAluno, err
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Se o aluno não existe, cria o novo aluno.
	if err := s.db.Create(aluno).Error; err != nil {
		return nil, err
	}

	// Após criar o aluno, realiza a inscrição no curso selecionado.
	newInscricao := models.Inscricao{
		AlunoID:       aluno.ID,
		CursoID:       selectedCurso.ID,
		DataInscricao: time.Now(),
	}
	if err := s.db.Create(&newInscricao).Error; err != nil {
		return nil, err
	}

	// Decrementa as vagas do curso
	curso.Vagas--
	if err := s.db.Save(&curso).Error; err != nil {
		return nil, err
	}

	// Opcional: recarrega o aluno com os cursos inscritos.
	err = s.db.Preload("Cursos").First(aluno, aluno.ID).Error
	return aluno, err
}
