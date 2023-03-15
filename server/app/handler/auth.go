package handler

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"server/app/model"
	"server/app/repo"
	"server/pkg/hash"
	"server/pkg/server"
	"server/pkg/session"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB, a fiber.Router) *AuthHandler {
	handler := AuthHandler{db: db}
	handler.setup(a)
	return &handler
}

func (h *AuthHandler) setup(a fiber.Router) {
	a.Group("/auth").
		Post("login", server.WithJSONBody(h.login)).
		Post("register", server.WithJSONBody(h.register)).
		Post("logout", h.logout)
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type userVo struct {
	Id       hash.ID `json:"id"`
	Nickname string  `json:"nickname"`
}

func (h *AuthHandler) login(c *fiber.Ctx, input *loginInput) (*userVo, error) {
	user, err := repo.GetUserByUsername(c.Context(), h.db, input.Username)
	if errors.Is(err, sql.ErrNoRows) {
		inputErr := server.NewInputErr("用户名不存在")
		return nil, inputErr
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, server.NewInputErr("密码错误")
	}
	sess, err := session.GetSession(c)
	sess.Set("uid", user.Id)
	if err = sess.Save(); err != nil {
		return nil, err
	}
	return &userVo{
		Id:       hash.ID(user.Id),
		Nickname: user.Nickname,
	}, nil
}

func (h *AuthHandler) register(c *fiber.Ctx, input *loginInput) (*userVo, error) {
	user, err := repo.GetUserByUsername(c.Context(), h.db, input.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if user != nil {
		return nil, server.NewInputErr("用户名已存在")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user = &model.User{
		Nickname: input.Username,
		Username: input.Username,
		Password: string(hashedPassword),
	}
	err = repo.InsertUser(c.Context(), h.db, user)
	if err != nil {
		return nil, err
	}
	return &userVo{
		Id:       hash.ID(user.Id),
		Nickname: user.Nickname,
	}, nil
}

func (h *AuthHandler) logout(c *fiber.Ctx) error {
	sess, err := session.GetSession(c)
	if err != nil {
		return err
	}
	if err = sess.Destroy(); err != nil {
		return err
	}
	return nil
}

func (h *AuthHandler) changePassword(c context.Context, input *loginInput) error {
	// todo
	return nil
}

func (h *AuthHandler) changeNickname(c context.Context, input *loginInput) error {
	// todo
	return nil
}
