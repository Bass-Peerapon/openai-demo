package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Order struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Orders []Order

func (o *Orders) Scan(v interface{}) error {
	switch v := v.(type) {
	case []byte:
		json.Unmarshal(v, &o)
		return nil
	case string:
		json.Unmarshal([]byte(v), &o)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (o Orders) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type Customer struct {
	ID         *uuid.UUID `json:"id" db:"id"`
	FirstName  string     `json:"first_name" db:"first_name"`
	LastName   string     `json:"last_name" db:"last_name"`
	Age        int        `json:"age" db:"age"`
	Membership string     `json:"membership" db:"membership"`
	Orders     Orders     `json:"orders" db:"orders"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
