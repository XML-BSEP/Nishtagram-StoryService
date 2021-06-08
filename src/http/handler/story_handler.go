package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"story-service/dto"
	"story-service/http/middleware"
	"story-service/usecase"
)

type StoryHandler interface {
	AddStory(ctx *gin.Context)
	RemoveStory(ctx *gin.Context)
	GetStoriesForUser(ctx *gin.Context)
	GetStoriesInUserProfile(ctx *gin.Context)
}

type storyHandler struct {
	storyUseCase usecase.StoryUseCase
}

func (s storyHandler) GetStoriesInUserProfile(ctx *gin.Context) {
	var req dto.UserDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	userRequested, _ := middleware.ExtractUserId(ctx.Request)
	var err error
	var stories []dto.StoryDTO
	stories, err = s.storyUseCase.GetAllStoriesByUser(req.UserId, userRequested, context.Background())
	/*
	if userRequested == req.UserId {
		stories, err = s.storyUseCase.GetAllStoriesByUser(req.UserId, context.Background())
	} else {
		stories, err = s.storyUseCase.GetActiveUsersStories(req.UserId, context.Background())
	}*/

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, stories)
}

func (s storyHandler) AddStory(ctx *gin.Context) {
	var req dto.StoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}
	req.UserId, _ = middleware.ExtractUserId(ctx.Request)
	err := s.storyUseCase.AddStory(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "story successfully added"})
}

func (s storyHandler) RemoveStory(ctx *gin.Context) {
	var req dto.RemoveStoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	err := s.storyUseCase.RemoveStory(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "story successfully removed"})
}

func (s storyHandler) GetStoriesForUser(ctx *gin.Context) {
	/*var req dto.RemoveStoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}*/

	userId, _ := middleware.ExtractUserId(ctx.Request)

	stories, err := s.storyUseCase.GetAllStoriesForOneUser(context.Background(), userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, stories)
}

func NewStoryHandler(storyUseCase usecase.StoryUseCase) StoryHandler {
	return &storyHandler{storyUseCase: storyUseCase}
}
