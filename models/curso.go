package models

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove quotes
	s := strings.Trim(string(b), "\"")
	// Parse using the custom layout
	t, err := time.Parse("2006/01/02", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type Curso struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome         string     `gorm:"not null" json:"nome"`
	Professor    string     `gorm:"not null" json:"professor"`
	Data         CustomTime `gorm:"not null" json:"data"`
	CargaHoraria int32      `gorm:"not null" json:"cargaHoraria"`
	Certificado  string     `gorm:"not null" json:"certificado"`
	AlunoID      *uint      `gorm:"default:null" json:"alunoId,omitempty"`
	Aluno        *Aluno     `gorm:"foreignKey:AlunoID" json:"aluno,omitempty"`
}
