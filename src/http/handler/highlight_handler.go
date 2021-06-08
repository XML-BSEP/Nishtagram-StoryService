package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"story-service/dto"
	"story-service/http/middleware"
	"story-service/usecase"
)

type HighlightHandler interface {
	AddStoryToHighlight(context *gin.Context)
	RemoveStoryFromHighlight(context *gin.Context)
	GetHighlightsByUser(context *gin.Context)
	GetStoriesInHighlight(context *gin.Context)
	SaveHighlight(context *gin.Context)
}

type highlightHandler struct {
	highlightUseCase usecase.HighlightUseCase
}

func (h highlightHandler) SaveHighlight(ctx *gin.Context) {
	var req dto.NewHighlight

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	req.UserId, _ = middleware.ExtractUserId(ctx.Request)

	err := h.highlightUseCase.UpdateHighlights(req, context.Background())

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "successfully removed story"})
}

func (h highlightHandler) AddStoryToHighlight(ctx *gin.Context) {
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}
	req.UserId, _ = middleware.ExtractUserId(ctx.Request)
	err := h.highlightUseCase.AddStoryToHighlight(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "successfully added story"})

}

func (h highlightHandler) RemoveStoryFromHighlight(ctx *gin.Context) {
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	err := h.highlightUseCase.RemoveStoryFrom(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "successfully removed story"})

}

func (h highlightHandler) GetHighlightsByUser(ctx *gin.Context) {
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	req.UserId, _ = middleware.ExtractUserId(ctx.Request)
	highlights, err := h.highlightUseCase.GetHighlights(context.Background(), req.UserId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, highlights)
}

func (h highlightHandler) GetStoriesInHighlight(ctx *gin.Context) {
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	highlights, err := h.highlightUseCase.GetHighlightByName(context.Background(), req.UserId, req.HighlightName)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, highlights)
}

func NewHighlightHandler(highlightUseCase usecase.HighlightUseCase) HighlightHandler {
	return &highlightHandler{highlightUseCase: highlightUseCase}
}
