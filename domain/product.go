package domain

import (
	"github.com/gofrs/uuid"
	"github.com/pgvector/pgvector-go"
)

type Product struct {
	ID        *uuid.UUID      `json:"id" db:"id"`
	Title     string          `json:"title" db:"title"`
	Content   string          `json:"content" db:"content"`
	Embedding pgvector.Vector `json:"embedding" db:"embedding"`
}
