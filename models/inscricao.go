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

	// Associações (opcionais, mas úteis para carregamento relacionado)
	Aluno Aluno `gorm:"foreignKey:AlunoID" json:"aluno,omitempty"`
	Curso Curso `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
}
