package postgres

import (
	"context"
	"encoding/json"

	"github.com/Bass-Peerapon/openai-demo/domain"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	Client *sqlx.DB
}

func NewChatRepository(client *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		Client: client,
	}
}

func (r *ChatRepository) GetChatHistory(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	query := `SELECT * FROM chats WHERE id = $1`
	rows, err := r.Client.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*domain.Chat, 0, 1)
	for rows.Next() {
		chat := &domain.Chat{}
		if err := rows.StructScan(chat); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if len(chats) > 0 {
		return chats[0], nil
	}
	return nil, nil
}

func (r *ChatRepository) SaveChatHistory(ctx context.Context, chat *domain.Chat) error {
	query := `INSERT INTO chats (id, user_id, messages, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (id) DO UPDATE SET messages = $3, updated_at = $5
	`

	historyJson, _ := json.Marshal(chat.Messages)
	_, err := r.Client.ExecContext(ctx, query, chat.ID, chat.UserID, historyJson, chat.CreatedAt, chat.UpdatedAt)
	return err
}
