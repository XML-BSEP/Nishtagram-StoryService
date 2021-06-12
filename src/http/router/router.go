package router

import (
	"github.com/gin-gonic/gin"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/http/handler"
	"story-service/http/middleware"
)

func NewRouter(handler handler.AppHandler, logger *logger.Logger) *gin.Engine {
	router := gin.Default()

	g := router.Group("/story")

	g.Use(middleware.AuthMiddleware(logger))
	g.Use(gin.Logger())

	g.POST("/addStory", handler.AddStory)
	g.POST("/removeStory", handler.RemoveStory)
	g.GET("/getStories", handler.GetStoriesForUser)
	g.POST("/getHighlights", handler.GetHighlightsByUser)
	g.POST("/getStoriesHighlight", handler.GetStoriesInHighlight)
	g.POST("/addToHighlight", handler.AddStoryToHighlight)
	g.POST("/removeFromHighlight", handler.RemoveStoryFromHighlight)
	g.POST("/getAllStoriesOnUserProfile", handler.GetStoriesInUserProfile)
	g.POST("/saveHighlight", handler.SaveHighlight)

	return router
}
