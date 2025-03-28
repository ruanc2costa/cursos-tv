package models

import "time"

// Aluno representa o modelo equivalente à entidade Java.
type Aluno struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome     string `gorm:"not null" json:"nome"`
	CPF      string `gorm:"not null" json:"cpf"`
	Email    string `gorm:"not null" json:"email"`
	Sexo     string `gorm:"not null" json:"sexo"`
	Telefone string `json:"telefone"`
	// DataNascto pode ser customizada na serialização JSON se for necessário o padrão dd/MM/yyyy.
	DataNascto time.Time `gorm:"not null" json:"dataNascto"`
}
