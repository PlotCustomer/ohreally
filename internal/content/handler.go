package content

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler exposes the Content Management function over HTTP.
type Handler struct {
	repo  Repository
	newID func() string
}

// NewHandler returns a Handler backed by repo.
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo:  repo,
		newID: func() string { return uuid.NewString() },
	}
}

// Register attaches the Content Management routes to the router.
func (h *Handler) Register(router *gin.RouterGroup) {
	router.GET("/content", h.list)
	router.POST("/content", h.create)
	router.GET("/content/:id", h.get)
	router.PUT("/content/:id", h.update)
	router.DELETE("/content/:id", h.delete)
}

type contentRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *Handler) create(c *gin.Context) {
	var req contentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item := Content{Title: req.Title, Body: req.Body}
	if err := item.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	item.ID = h.newID()
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	if err := h.repo.Create(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create content"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *Handler) get(c *gin.Context) {
	item, err := h.repo.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeRepoError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.repo.List(c.Request.Context())
	if err != nil {
		writeRepoError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) update(c *gin.Context) {
	var req contentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item := Content{
		ID:        c.Param("id"),
		Title:     req.Title,
		Body:      req.Body,
		UpdatedAt: time.Now().UTC(),
	}
	if err := item.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), &item); err != nil {
		writeRepoError(c, err)
		return
	}

	updated, err := h.repo.Get(c.Request.Context(), item.ID)
	if err != nil {
		writeRepoError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		writeRepoError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeRepoError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
