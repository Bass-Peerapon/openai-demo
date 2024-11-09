package openai

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/Bass-Peerapon/openai-demo/domain"
	myTemplate "github.com/Bass-Peerapon/openai-demo/template"
	"github.com/gofrs/uuid"
)

var ErrCustomerNotFound = errors.New("customer not found")

type OpenaiService interface {
	Chat(ctx context.Context, meesage []domain.Message) (*domain.Message, error)
}

type CustomerService interface {
	GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
}

type ProductService interface {
	SearchProduct(ctx context.Context, query string) ([]domain.Product, error)
}

type ChatHistoryService interface {
	SaveChatHistory(ctx context.Context, chat *domain.Chat) error
	GetChatHistory(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
}

type Service struct {
	openaiRepo OpenaiService
	customer   CustomerService
	product    ProductService
	chat       ChatHistoryService
}

func NewService(openaiRepo OpenaiService, customer CustomerService, product ProductService, chat ChatHistoryService) *Service {
	return &Service{
		openaiRepo: openaiRepo,
		customer:   customer,
		product:    product,
		chat:       chat,
	}
}

func (s *Service) NewChat(ctx context.Context, userID uuid.UUID, message string) (*domain.Chat, error) {
	customer, err := s.customer.GetCustomer(ctx, userID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}
	products, err := s.product.SearchProduct(ctx, message)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Products []domain.Product
		Customer domain.Customer
	}

	data := Data{
		Products: products,
		Customer: *customer,
	}
	t := template.Must(template.New("metaprompt").Parse(string(myTemplate.Metaprompt)))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	chat := domain.NewChat(userID)
	chat.AddMessage(domain.RoleSystem, buf.String())
	chat.AddMessage(domain.RoleUser, message)

	ansMessage, err := s.openaiRepo.Chat(ctx, chat.Messages)
	if err != nil {
		return nil, err
	}
	if ansMessage == nil {
		return nil, nil
	}
	chat.AddMessage(domain.RoleAssistant, ansMessage.Content)

	err = s.chat.SaveChatHistory(ctx, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *Service) GetChatHistory(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	return s.chat.GetChatHistory(ctx, id)
}

func (s *Service) Chat(ctx context.Context, chatID uuid.UUID, message string) (*domain.Message, error) {
	chat, err := s.chat.GetChatHistory(ctx, chatID)
	if err != nil {
		return nil, err
	}

	chat.AddMessage(domain.RoleUser, message)
	ansMessage, err := s.openaiRepo.Chat(ctx, chat.Messages)
	if err != nil {
		return nil, err
	}
	if ansMessage == nil {
		return nil, nil
	}
	chat.AddMessage(domain.RoleAssistant, ansMessage.Content)

	if err := s.chat.SaveChatHistory(ctx, chat); err != nil {
		return nil, err
	}
	return ansMessage, nil
}
