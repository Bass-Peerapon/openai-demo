package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

func (r Role) String() string {
	return string(r)
}

type Message struct {
	ID        *uuid.UUID `json:"id"`
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Messages []Message

func (h *Messages) Scan(v interface{}) error {
	switch v := v.(type) {
	case []byte:
		json.Unmarshal(v, &h)
		return nil
	case string:
		json.Unmarshal([]byte(v), &h)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (h Messages) Value() (driver.Value, error) {
	return json.Marshal(h)
}

type Chat struct {
	ID        *uuid.UUID `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id" db:"user_id"`
	Messages  Messages   `json:"messages" db:"messages"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

func NewChat(userID uuid.UUID) *Chat {
	id, _ := uuid.NewV4()

	now := time.Now()
	return &Chat{
		ID:        &id,
		UserID:    &userID,
		Messages:  make(Messages, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Chat) AddMessage(role Role, content string) {
	now := time.Now()
	id, _ := uuid.NewV4()
	c.Messages = append(c.Messages, Message{
		ID:        &id,
		Role:      role.String(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	})
}
