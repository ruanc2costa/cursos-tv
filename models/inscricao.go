package models

import (
	"time"
)

// Inscricao representa o registro de inscrição de um aluno em um curso.
type Inscricao struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AlunoID       uint      `gorm:"not null" json:"alunoId"`
	CursoID       uint      `gorm:"not null" json:"cursoId"`
	DataInscricao time.Time `gorm:"not null" json:"dataInscricao"`

	// Novos campos adicionados
	Escolaridade      string `json:"escolaridade"`
	Trabalhando       string `json:"trabalhando"`
	Bairro            string `json:"bairro"`
	EhCuidador        string `json:"ehCuidador"`
	EhPCD             string `json:"ehPCD"`
	TipoPCD           string `json:"tipoPCD"`
	NecessitaElevador string `json:"necessitaElevador"`
	ComoSoube         string `json:"comoSoube"`
	AutorizaWhatsApp  string `json:"autorizaWhatsApp"`
	LevaNotebook      string `json:"levaNotebook"`

	// Associações
	Aluno Aluno `gorm:"foreignKey:AlunoID" json:"aluno,omitempty"`
	Curso Curso `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
}
