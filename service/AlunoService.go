package service

import (
	"errors"
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
// Se o aluno já existir (verificado pelo email):
//   - Se ele já estiver inscrito no curso informado, nada é feito.
//   - Caso contrário, é criada uma nova inscrição para o curso.
//
// Se o aluno não existir, ele é criado e, em seguida, é realizada a inscrição.
func (s *AlunoService) AddAluno(aluno *models.Aluno) (*models.Aluno, error) {
	// Validação básica: deve haver ao menos um curso informado.
	if aluno == nil {
		return nil, errors.New("aluno é nil")
	}
	if len(aluno.Cursos) == 0 {
		return nil, errors.New("deve ser informado pelo menos um curso")
	}

	// Seleciona o curso informado (no exemplo, usamos o primeiro curso do array).
	selectedCurso := aluno.Cursos[0]

	// Tenta encontrar um aluno existente pelo email e carrega os cursos associados.
	var existingAluno models.Aluno
	err := s.db.Where("email = ?", aluno.Email).Preload("Cursos").First(&existingAluno).Error
	if err == nil {
		// Aluno já existe; verifica se já está inscrito no curso.
		for _, c := range existingAluno.Cursos {
			if c.ID == selectedCurso.ID {
				// Já está inscrito; não precisa criar nova inscrição.
				return &existingAluno, nil
			}
		}

		// Se não estiver inscrito, cria uma nova inscrição.
		newInscricao := models.Inscricao{
			AlunoID:       existingAluno.ID,
			CursoID:       selectedCurso.ID,
			DataInscricao: time.Now(),
		}
		if err := s.db.Create(&newInscricao).Error; err != nil {
			return nil, err
		}
		// Opcional: recarregar o aluno com os cursos atualizados.
		err = s.db.Preload("Cursos").First(&existingAluno, existingAluno.ID).Error
		return &existingAluno, err
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Erro ao buscar aluno que não seja "não encontrado"
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

	// Opcional: recarrega o aluno com os cursos inscritos.
	err = s.db.Preload("Cursos").First(&aluno, aluno.ID).Error
	return aluno, err
}

// GetAllAlunos retorna todos os alunos com seus cursos.
func (s *AlunoService) GetAllAlunos() ([]models.Aluno, error) {
	var alunos []models.Aluno
	err := s.db.Preload("Cursos").Find(&alunos).Error
	return alunos, err
}

// GetAluno busca um aluno pelo ID, carregando os cursos associados.
func (s *AlunoService) GetAluno(id uint) (*models.Aluno, error) {
	var aluno models.Aluno
	err := s.db.Preload("Cursos").First(&aluno, id).Error
	if err != nil {
		return nil, err
	}
	return &aluno, nil
}

// DeleteAluno remove um aluno com base no ID.
func (s *AlunoService) DeleteAluno(id uint) error {
	var aluno models.Aluno
	if err := s.db.First(&aluno, id).Error; err != nil {
		return err
	}
	return s.db.Delete(&aluno).Error
}
