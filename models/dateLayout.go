package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const DateLayout = "02/01/2006" // dd/MM/yyyy

// CustomDate encapsula time.Time para utilizar um formato customizado
type CustomDate struct {
	time.Time
}

// UnmarshalJSON implementa a desserialização do JSON para o CustomDate.
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	// Remove as aspas
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("data inválida")
	}
	s = s[1 : len(s)-1]
	// Tenta converter a string para time.Time com o layout definido
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

// MarshalJSON implementa a serialização do CustomDate para JSON.
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	formatted := cd.Format(DateLayout)
	return json.Marshal(formatted)
}
