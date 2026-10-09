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
	repo Repository
	now  func() time.Time
}

// NewHandler creates a Content Management HTTP handler.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo, now: time.Now}
}

// Register mounts the Content Management routes on the given router.
func (h *Handler) Register(router gin.IRouter) {
	router.POST("", h.create)
	router.GET("", h.list)
	router.GET("/:id", h.get)
	router.PUT("/:id", h.update)
	router.DELETE("/:id", h.delete)
}

type contentRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body"`
}

func (h *Handler) create(c *gin.Context) {
	var req contentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := h.now().UTC()
	item := &Content{
		ID:        uuid.New(),
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.repo.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create content"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list content"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	item, err := h.repo.Get(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read content"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req contentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.repo.Get(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read content"})
		return
	}

	item.Title = req.Title
	item.Body = req.Body
	item.UpdatedAt = h.now().UTC()

	if err := h.repo.Update(c.Request.Context(), item); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update content"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete content"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content id"})
		return uuid.Nil, false
	}
	return id, true
}
