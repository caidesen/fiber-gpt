package handler

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
	"server/app/model"
	"server/app/repo"
	"server/pkg/config"
	"server/pkg/hash"
	"server/pkg/server"
	"server/pkg/session"
	"time"
)

type ChatHandler struct {
	db *sql.DB
}

func NewChatHandler(db *sql.DB, a fiber.Router) *ChatHandler {
	handler := ChatHandler{db: db}
	handler.setup(a)
	return &handler
}

func (h *ChatHandler) setup(a fiber.Router) {
	a.Use(session.AuthRequired)
	a.Post("getChats", h.getChats)
	a.Post("createChat", h.createChat)
	a.Post("getQuestions", server.WithJSONBody(h.getQuestionsByChatId))
	a.Post("createQuestion", server.WithJSONBody(h.createQuestion))
	a.Post("readAnswer", server.WithJSONBody(h.readAnswer))
}

type questionVo struct {
	Id        hash.ID   `json:"id"`
	Q         string    `json:"q"`
	A         string    `json:"a"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type chatVo struct {
	Id        hash.ID      `json:"id"`
	Questions []questionVo `json:"questions"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (h *ChatHandler) getChats(c *fiber.Ctx) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	chats, err := repo.FindChatsWithLatestQuestionByUserId(c.Context(), h.db, userId)
	if err != nil {
		return err
	}
	chatVos := make([]chatVo, 0)
	for _, chat := range chats {
		questionVos := make([]questionVo, 0)
		for _, question := range chat.Questions {
			questionVos = append(questionVos, questionVo{
				Id:        hash.ID(question.Id),
				Q:         question.Q,
				A:         question.A,
				CreatedAt: question.CreatedAt,
				UpdatedAt: question.UpdatedAt,
			})
		}
		chatVos = append(chatVos, chatVo{
			Id:        hash.ID(chat.Id),
			Questions: questionVos,
			CreatedAt: chat.CreatedAt,
			UpdatedAt: chat.UpdatedAt,
		})
	}
	return c.JSON(chatVos)
}
func (h *ChatHandler) createChat(c *fiber.Ctx) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	settings, err := GetUserSettings(c.Context(), h.db, userId)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	chat := model.Chat{
		UserId: userId,
	}
	chat.Settings.Temperature = settings.Temperature
	chat.Settings.MaxToken = settings.MaxToken
	if err = repo.InsertChat(c.Context(), h.db, &chat); err != nil {
		return err
	}
	return c.JSON(chatVo{
		Id:        hash.ID(chat.Id),
		Questions: []questionVo{},
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	})
}

type createQuestionInput struct {
	ChatId hash.ID `json:"chatId"`
	Q      string  `json:"q"`
}

func (h *ChatHandler) createQuestion(c *fiber.Ctx, input createQuestionInput) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	chat, err := repo.FindOneChatById(c.Context(), h.db, input.ChatId.Uint())
	if err != nil {
		return err
	}
	if chat.UserId != userId {
		return server.NewInputErr("chat is not yours")
	}
	question := model.Question{
		ChatId: chat.Id,
		Q:      input.Q,
		Status: "pending",
	}
	if err = repo.InsertQuestion(c.Context(), h.db, &question); err != nil {
		return err
	}
	return c.JSON(questionVo{
		Id:        hash.ID(question.Id),
		Q:         question.Q,
		A:         question.A,
		Status:    question.Status,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
	})
}

type chatIdInput struct {
	ChatId hash.ID `json:"chatId"`
}

func (h *ChatHandler) getQuestionsByChatId(c *fiber.Ctx, input *chatIdInput) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	if chat, err := repo.FindOneChatById(c.Context(), h.db, input.ChatId.Uint()); err != nil {
		return err
	} else {
		if chat.UserId != userId {
			return server.NewInputErr("chat is not yours")
		}
	}
	questions, err := repo.FindQuestionsByChatId(c.Context(), h.db, input.ChatId.Uint())
	if err != nil {
		return err
	}
	questionVos := make([]questionVo, 0)
	for _, question := range questions {
		questionVos = append(questionVos, questionVo{
			Id:        hash.ID(question.Id),
			Q:         question.Q,
			A:         question.A,
			Status:    question.Status,
			CreatedAt: question.CreatedAt,
			UpdatedAt: question.UpdatedAt,
		})
	}
	return c.JSON(questionVos)
}

type readAnswerInput struct {
	QId hash.ID `json:"qid"`
}

func createOpenaiClient(apiKey string) *openai.Client {
	proxyUrl := config.GetConfig().ProxyUrl
	if proxyUrl == "" {
		return openai.NewClient(apiKey)
	}
	clientConfig := openai.DefaultConfig(apiKey)
	url, err := url.Parse(proxyUrl)
	if err != nil {
		panic(fmt.Errorf("failed to parse proxy url: %w", err))
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(url),
	}
	clientConfig.HTTPClient = &http.Client{
		Transport: transport,
	}
	return openai.NewClientWithConfig(clientConfig)
}
func (h *ChatHandler) readAnswer(c *fiber.Ctx, input *readAnswerInput) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	setting, err := GetUserSettings(c.Context(), h.db, userId)
	if err != nil {
		return err
	}
	client := createOpenaiClient(setting.ApiKey)
	question, err := repo.FindOneQuestionById(c.Context(), h.db, input.QId.Uint())
	if err != nil {
		return err
	}
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   setting.MaxToken,
		Temperature: setting.Temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question.Q,
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(c.Context(), req)
	if err != nil {
		question.Status = "failed"
		question.A = err.Error()
		repo.UpdateQuestion(context.Background(), h.db, question)
		return err
	}
	c.Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer stream.Close()
		var answer string
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				w.WriteString(fmt.Sprintf("\n[ERROR]\n%v", err.Error()))
				w.Flush()
				break
			}
			content := response.Choices[0].Delta.Content
			answer += content
			w.WriteString(content)
			w.Flush()
		}
		question.A = answer
		if err != nil {
			question.Status = "failed"
		} else {
			question.Status = "success"
		}
		repo.UpdateQuestion(context.Background(), h.db, question)
	})
	return nil
}
