package models

// Curso representa a entidade Curso.
type Curso struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome         string     `gorm:"not null" json:"nome"`
	Professor    string     `gorm:"not null" json:"professor"`
	DataInicio   CustomDate `gorm:"not null" json:"data"`
	CargaHoraria int32      `gorm:"not null" json:"cargaHoraria"`
	Certificado  string     `gorm:"not null" json:"certificado"`

	// Permite que a coluna aluno_id seja nula.
	AlunoID *uint  `gorm:"default:null" json:"alunoId,omitempty"`
	Aluno   *Aluno `gorm:"foreignKey:AlunoID" json:"aluno,omitempty"`
}
