package router

import (
	"github.com/gin-gonic/gin"
	"story-service/http/handler"
	"story-service/http/middleware"
)

func NewRouter(handler handler.AppHandler) *gin.Engine {
	router := gin.Default()

	g := router.Group("/story")

	g.Use(middleware.AuthMiddleware())
	g.Use(gin.Logger())

	g.POST("/addStory", handler.AddStory)
	g.POST("/removeStory", handler.RemoveStory)
	g.GET("/getStories", handler.GetStoriesForUser)
	g.POST("/getHighlights", handler.GetHighlightsByUser)
	g.POST("/getStoriesHighlight", handler.GetStoriesInHighlight)
	g.POST("/addToHighlight", handler.AddStoryToHighlight)
	g.POST("/removeFromHighlight", handler.RemoveStoryFromHighlight)

	return router
}
