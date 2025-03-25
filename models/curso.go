package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// CustomTime define um tipo customizado para datas.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON permite tratar o JSON no formato "dd/MM/yyyy".
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove as aspas do JSON.
	s := strings.Trim(string(b), "\"")
	// Faz o parse usando o layout customizado.
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Value implementa a interface driver.Valuer para salvar no banco.
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan implementa a interface sql.Scanner para ler do banco.
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	default:
		return fmt.Errorf("não foi possível converter %T para CustomTime", value)
	}
}

// Curso representa um curso com data, carga horária, certificado e controle de vagas.
type Curso struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome             string     `gorm:"not null" json:"nome"`
	Professor        string     `gorm:"not null" json:"professor"`
	Data             CustomTime `gorm:"not null" json:"data"`
	CargaHoraria     int32      `gorm:"not null" json:"cargaHoraria"`
	Certificado      string     `gorm:"not null" json:"certificado"`
	VagasTotais      int32      `gorm:"not null" json:"vagasTotais"`
	VagasPreenchidas int32      `gorm:"not null" json:"vagasPreenchidas"`

	AlunoID *uint  `gorm:"default:null" json:"alunoId,omitempty"`
	Aluno   *Aluno `gorm:"foreignKey:AlunoID" json:"aluno,omitempty"`
}
