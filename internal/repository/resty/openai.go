package resty

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Bass-Peerapon/openai-demo/domain"
	"github.com/go-resty/resty/v2"
)

type OpenaiService struct {
	client *resty.Client
	Model  string `json:"model"`
}

func NewOpenaiService(host string, secret string, model string) *OpenaiService {
	client := resty.New()
	client.SetBaseURL(host)
	client.SetAuthToken(secret)
	return &OpenaiService{
		client: client,
		Model:  model,
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type chatResponse struct {
	Message message `json:"message"`
}

type getEmbeddingRequest struct {
	Model    string `json:"model"`
	Input    string `json:"input"`
	Truncate bool   `json:"truncate"`
}

type getEmbeddingResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *OpenaiService) Chat(ctx context.Context, meesages []domain.Message) (*domain.Message, error) {
	messages := make([]message, len(meesages))

	for i, m := range meesages {
		messages[i] = message{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	req := chatRequest{
		Model:    o.Model,
		Messages: messages,
		Stream:   false,
	}

	resp, err := o.client.R().
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		Post("/api/chat")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, errors.New(resp.String())
	}

	chatResp := chatResponse{}
	err = json.Unmarshal(resp.Body(), &chatResp)
	if err != nil {
		return nil, err
	}

	ansMessage := domain.Message{
		Role:    domain.RoleAssistant.String(),
		Content: chatResp.Message.Content,
	}

	return &ansMessage, nil
}

func (o *OpenaiService) GetEmbedding(ctx context.Context, text string) ([][]float32, error) {
	req := getEmbeddingRequest{
		Model:    o.Model,
		Input:    text,
		Truncate: true,
	}

	resp, err := o.client.R().
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		Post("/api/embed")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, errors.New(resp.String())
	}

	var embedding getEmbeddingResponse
	err = json.Unmarshal(resp.Body(), &embedding)
	if err != nil {
		return nil, err
	}
	return embedding.Embeddings, nil
}
