package rest

import (
	"context"
	"net/http"

	"github.com/Bass-Peerapon/openai-demo/domain"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type OpenaiService interface {
	NewChat(ctx context.Context, userID uuid.UUID, message string) (*domain.Chat, error)
	Chat(ctx context.Context, chatID uuid.UUID, message string) (*domain.Message, error)
	GetChatHistory(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
}

type OpenaiHandler struct {
	service OpenaiService
}

func NewOpenaiHandler(e *echo.Echo, service OpenaiService) {
	handler := &OpenaiHandler{
		service: service,
	}

	e.POST("/api/chat", handler.NewChat)
	e.POST("/api/chat/:chat_id", handler.Chat)
	e.GET("/api/chat/:chat_id", handler.GetChatHistory)
}

// Request struct for incoming questions
type questionRequest struct {
	Question string `json:"question" validate:"required"`
}

type chatResponse struct {
	Chat *domain.Chat `json:"chat"`
}

type messageResponse struct {
	Message *domain.Message `json:"message"`
}

func (h *OpenaiHandler) Chat(c echo.Context) error {
	chatID := uuid.FromStringOrNil(c.Param("chat_id"))
	req := new(questionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if req.Question == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Question is required")
	}

	message, err := h.service.Chat(c.Request().Context(), chatID, req.Question)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, messageResponse{message})
}

func (h *OpenaiHandler) NewChat(c echo.Context) error {
	req := new(questionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if req.Question == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Question is required")
	}
	userID := uuid.FromStringOrNil(c.Request().Header.Get("X-User-ID"))
	chat, err := h.service.NewChat(c.Request().Context(), userID, req.Question)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, chatResponse{chat})
}

func (h *OpenaiHandler) GetChatHistory(c echo.Context) error {
	chatID := uuid.FromStringOrNil(c.Param("chat_id"))
	chat, err := h.service.GetChatHistory(c.Request().Context(), chatID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if chat == nil {
		return echo.NewHTTPError(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, chatResponse{chat})
}
