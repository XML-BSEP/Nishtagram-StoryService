package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/domain"
	"story-service/dto"
	"story-service/http/middleware"
	"story-service/usecase"
	"strings"
)

type ReportHandler interface {
	ReportStory(ctx *gin.Context)
	GetAllReportTypes(ctx *gin.Context)
}

type reportHandler struct {
	reportUseCase usecase.ReportUseCase
	logger *logger.Logger
}

func (r reportHandler) ReportStory(ctx *gin.Context) {
	var req dto.ReportStory

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&req); err != nil {
		r.logger.Logger.Errorf("error while decoding json, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "invalid request"})
		ctx.Abort()
		return
	}
	userId, _ := middleware.ExtractUserId(ctx.Request, r.logger)
	reportStory := domain.StoryReport{ReportedStoryBy: domain.Profile{Id: req.StoryBy}, StoryId: req.StoryId,
		ReportedBy: domain.Profile{Id: userId}, ReportType: domain.ReportType{Type: strings.ToUpper(req.ReportType)}}

	err := r.reportUseCase.ReportStory(reportStory, ctx)

	if err != nil {
		r.logger.Logger.Errorf("error while adding report, error: %v\n", err)
		ctx.JSON(400, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "Successfully added report"})

}


func (r reportHandler) GetAllReportTypes(context *gin.Context) {
	types, err := r.reportUseCase.GetAllReportType(context)

	if err != nil {
		context.JSON(500, gin.H{"message" : "server error"})
		context.Abort()
		return
	}

	context.JSON(200, types)
}

func NewReportHandler(reportUseCase usecase.ReportUseCase, logger *logger.Logger) ReportHandler {
	return &reportHandler{reportUseCase: reportUseCase, logger: logger}
}
