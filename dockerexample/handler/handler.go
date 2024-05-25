package handler

import (
	"log"
	"net/http"

	"github.com/gendutski/my-playground/dockerexample/core/entity"
	"github.com/gendutski/my-playground/dockerexample/core/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo repository.Repository
}

func Init(repo repository.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) Get(c echo.Context) error {
	p := new(entity.GetRequest)
	if err := c.Bind(p); err != nil {
		log.Println("Error binding:", err)
		return err
	}

	result, err := h.repo.Get(p.Offset, p.Limit)
	if err != nil {
		log.Println("Error get:", err)
		return err
	}
	return c.JSON(http.StatusOK, entity.ToGetResponse(result))
}

func (h *Handler) Post(c echo.Context) error {
	p := new(entity.User)
	if err := c.Bind(p); err != nil {
		log.Println("Error binding:", err)
		return err
	}

	err := h.repo.Set(p)
	if err != nil {
		log.Println("Error set:", err)
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}
