package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	logger "github.com/jelena-vlajkov/logger/logger"
	"github.com/microcosm-cc/bluemonday"
	"story-service/dto"
	"story-service/http/middleware"
	"story-service/usecase"
	"strings"
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
	logger *logger.Logger
}

func (h highlightHandler) SaveHighlight(ctx *gin.Context) {
	h.logger.Logger.Println("Handling SAVING HIGHLIGHT")
	var req dto.NewHighlight

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		h.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}
	policy := bluemonday.UGCPolicy()
	req.Id = strings.TrimSpace(policy.Sanitize(req.Id))
	req.HighlightName = strings.TrimSpace(policy.Sanitize(req.HighlightName))

	for _,story := range req.Stories {
		story = strings.TrimSpace(policy.Sanitize(story))
		if story == "" {
			h.logger.Logger.Errorf("fields are empty or xss attack happened")
			ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
			return
		}
	}


	if req.Id == "" || req.HighlightName == "" {
		h.logger.Logger.Errorf("fields are empty or xss attack happened")
		ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
		return
	}



	req.UserId, _ = middleware.ExtractUserId(ctx.Request, h.logger)

	err := h.highlightUseCase.UpdateHighlights(req, context.Background())

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "successfully removed story"})
}

func (h highlightHandler) AddStoryToHighlight(ctx *gin.Context) {
	h.logger.Logger.Println("Handling ADDING STORY TO HIGHLIGHT")
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		h.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	policy := bluemonday.UGCPolicy()
	req.HighlightName = strings.TrimSpace(policy.Sanitize(req.HighlightName))

	for _,story := range req.StoryId {
		story = strings.TrimSpace(policy.Sanitize(story))
		if story == "" {
			h.logger.Logger.Errorf("fields are empty or xss attack happened")
			ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
			return
		}
	}

	if req.HighlightName == "" {
		h.logger.Logger.Errorf("fields are empty or xss attack happened")
		ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
		return
	}

	req.UserId, _ = middleware.ExtractUserId(ctx.Request, h.logger)
	err := h.highlightUseCase.AddStoryToHighlight(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "successfully added story"})

}

func (h highlightHandler) RemoveStoryFromHighlight(ctx *gin.Context) {
	h.logger.Logger.Println("Handling REMOVING STORY FROM HIGHLIGHT")
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		h.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	policy := bluemonday.UGCPolicy()
	req.HighlightName = strings.TrimSpace(policy.Sanitize(req.HighlightName))

	for _,story := range req.StoryId {
		story = strings.TrimSpace(policy.Sanitize(story))
		if story == "" {
			h.logger.Logger.Errorf("fields are empty or xss attack happened")
			ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
			return
		}
	}

	if req.HighlightName == "" {
		h.logger.Logger.Errorf("fields are empty or xss attack happened")
		ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
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
	h.logger.Logger.Println("Handling GETTING HIGHLIGHTS BY USER")
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		h.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}



	req.UserId, _ = middleware.ExtractUserId(ctx.Request, h.logger)
	highlights, err := h.highlightUseCase.GetHighlights(context.Background(), req.UserId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, highlights)
}

func (h highlightHandler) GetStoriesInHighlight(ctx *gin.Context) {
	h.logger.Logger.Println("Handling GETTING STORIES IN HIGHLIGHTS")
	var req dto.HighlightDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		h.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	policy := bluemonday.UGCPolicy()
	req.UserId = strings.TrimSpace(policy.Sanitize(req.UserId))
	req.HighlightName = strings.TrimSpace(policy.Sanitize(req.HighlightName))


	if req.UserId == "" || req.HighlightName == "" {
		h.logger.Logger.Errorf("fields are empty or xss attack happened")
		ctx.JSON(400, gin.H{"message" : "Fields are empty or xss attack happened"})
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

func NewHighlightHandler(highlightUseCase usecase.HighlightUseCase, logger *logger.Logger) HighlightHandler {
	return &highlightHandler{highlightUseCase: highlightUseCase, logger: logger}
}
