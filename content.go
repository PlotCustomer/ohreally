package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrContentNotFound is returned when a Content record does not exist.
var ErrContentNotFound = errors.New("content not found")

// Content is the Data Object managed by the Content Management function of the
// Epileptic application component.
type Content struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// contentInput is the payload accepted when creating or updating Content.
type contentInput struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

func (in contentInput) validate() error {
	if strings.TrimSpace(in.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

func listContent(store ContentStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := store.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list content"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"content": items})
	}
}

func createContent(store ContentStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in contentInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content payload"})
			return
		}
		if err := in.validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now().UTC()
		item, err := store.Create(c.Request.Context(), Content{
			ID:        newContentID(),
			Title:     in.Title,
			Body:      in.Body,
			Author:    in.Author,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create content"})
			return
		}
		c.JSON(http.StatusCreated, item)
	}
}

func getContent(store ContentStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		item, err := store.Get(c.Request.Context(), c.Param("id"))
		if err != nil {
			writeContentError(c, err)
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func updateContent(store ContentStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in contentInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content payload"})
			return
		}
		if err := in.validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := store.Get(c.Request.Context(), c.Param("id"))
		if err != nil {
			writeContentError(c, err)
			return
		}
		item.Title = in.Title
		item.Body = in.Body
		item.Author = in.Author
		item.UpdatedAt = time.Now().UTC()

		item, err = store.Update(c.Request.Context(), item)
		if err != nil {
			writeContentError(c, err)
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func deleteContent(store ContentStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := store.Delete(c.Request.Context(), c.Param("id")); err != nil {
			writeContentError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func writeContentError(c *gin.Context, err error) {
	if errors.Is(err, ErrContentNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "content management failed"})
}
