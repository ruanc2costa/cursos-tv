package models

import "time"

// Aluno representa o modelo equivalente à entidade Java.
// Aluno representa o modelo equivalente à entidade Java.
// Versão corrigida
type Aluno struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome       string    `gorm:"not null" json:"nome"`
	CPF        string    `gorm:"not null;uniqueIndex" json:"cpf"`   // Adicionar uniqueIndex
	Email      string    `gorm:"not null;uniqueIndex" json:"email"` // Adicionar uniqueIndex
	Sexo       string    `gorm:"not null" json:"sexo"`
	Telefone   string    `json:"telefone"`
	DataNascto time.Time `gorm:"not null" json:"dataNascto"`

	// Remover relação direta com Curso
	// Em vez disso, podemos adicionar relação com Inscrições (se necessário)
	Inscricoes []Inscricao `gorm:"foreignKey:AlunoID" json:"inscricoes,omitempty"`
}
