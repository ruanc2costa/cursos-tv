package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const DateLayout = "02/01/2006" // dd/MM/yyyy

// CustomDate encapsula time.Time para utilizar um formato customizado
type CustomDate struct {
	time.Time
}

// UnmarshalJSON desserializa a data no formato dd/MM/yyyy
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	// Remove as aspas da string JSON
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("data inválida")
	}
	s = s[1 : len(s)-1]
	// Converte a string para time.Time usando o layout definido
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

// MarshalJSON serializa a data no formato dd/MM/yyyy
func (cd *CustomDate) MarshalJSON() ([]byte, error) {
	formatted := cd.Format(DateLayout)
	return json.Marshal(formatted)
}

// Value implementa o método driver.Valuer para que o CustomDate seja armazenado no banco como time.Time.
func (cd *CustomDate) Value() (driver.Value, error) {
	return cd.Time, nil
}

// Scan implementa o método sql.Scanner para converter dados do banco para CustomDate.
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		cd.Time = v
		return nil
	case []byte:
		t, err := time.Parse(DateLayout, string(v))
		if err != nil {
			return err
		}
		cd.Time = t
		return nil
	case string:
		t, err := time.Parse(DateLayout, v)
		if err != nil {
			return err
		}
		cd.Time = t
		return nil
	default:
		return fmt.Errorf("não é possível converter o tipo %T para CustomDate", value)
	}
}
