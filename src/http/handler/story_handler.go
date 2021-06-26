package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/dto"
	"story-service/http/middleware"
	"story-service/usecase"
)

type StoryHandler interface {
	AddStory(ctx *gin.Context)
	RemoveStory(ctx *gin.Context)
	GetStoriesForUser(ctx *gin.Context)
	GetStoriesInUserProfile(ctx *gin.Context)
	GetStoryByIdForAdmin(ctx *gin.Context)
}

type storyHandler struct {
	storyUseCase usecase.StoryUseCase
	logger *logger.Logger
}

func (s storyHandler) GetStoryByIdForAdmin(ctx *gin.Context) {
	s.logger.Logger.Println("Handling GETTING STORY FOR ADMIN")
	var req dto.GetStoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		s.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	story, err := s.storyUseCase.GetStoryByIdForAdmin(req.Id, req.StoryBy, context.Background())

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, story)
}

func (s storyHandler) GetStoriesInUserProfile(ctx *gin.Context) {
	s.logger.Logger.Println("Handling GETTING STORIES ON PROFILE")
	var req dto.UserDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		s.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}

	userRequested, _ := middleware.ExtractUserId(ctx.Request, s.logger)
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
	s.logger.Logger.Println("Handling ADDING STORIES")
	var req dto.StoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		s.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}
	req.UserId, _ = middleware.ExtractUserId(ctx.Request, s.logger)
	err := s.storyUseCase.AddStory(context.Background(), req)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "story successfully added"})
}

func (s storyHandler) RemoveStory(ctx *gin.Context) {
	s.logger.Logger.Println("Handling REMOVING STORY")
	var req dto.RemoveStoryDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		s.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
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
	s.logger.Logger.Println("Handling GETTING STORIES ON FEED")
	userId, _ := middleware.ExtractUserId(ctx.Request, s.logger)

	stories, err := s.storyUseCase.GetAllStoriesForOneUser(context.Background(), userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, stories)
}

func NewStoryHandler(storyUseCase usecase.StoryUseCase, logger *logger.Logger) StoryHandler {
	return &storyHandler{storyUseCase: storyUseCase, logger: logger}
}
