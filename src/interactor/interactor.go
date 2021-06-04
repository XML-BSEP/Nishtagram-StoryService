package interactor

import (
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"story-service/http/handler"
	"story-service/infrastructure/data_seeder"
	"story-service/repository"
	"story-service/usecase"
)

type appHandler struct {
	handler.HighlightHandler
	handler.StoryHandler
}


type Interactor interface {
	NewStoryRepository() repository.StoryRepo
	NewHighlightRepository() repository.HighlightRepo

	NewStoryUseCase() usecase.StoryUseCase
	NewHighlightUseCase() usecase.HighlightUseCase
	NewRedisUseCase() usecase.RedisUseCase

	NewAppHandler() handler.AppHandler
	NewStoryHandler() handler.StoryHandler
	NewHighlightHandler() handler.HighlightHandler
}

type interactor struct {
	cassandraClient *gocql.Session
	redisClient *redis.Client
}

func (i interactor) NewAppHandler() handler.AppHandler {
	appHandler := appHandler{}
	appHandler.StoryHandler = i.NewStoryHandler()
	appHandler.HighlightHandler = i.NewHighlightHandler()

	data_seeder.SeedData()
	return appHandler
}

func (i interactor) NewRedisUseCase() usecase.RedisUseCase {
	return usecase.NewRedisUsecase(i.redisClient)
}

func (i interactor) NewStoryRepository() repository.StoryRepo {
	return repository.NewStoryRepo(i.cassandraClient)
}

func (i interactor) NewHighlightRepository() repository.HighlightRepo {
	return repository.NewHighlightRepo(i.cassandraClient)
}

func (i interactor) NewStoryUseCase() usecase.StoryUseCase {
	return usecase.NewStoryUseCase(i.NewStoryRepository(), i.NewRedisUseCase())
}

func (i interactor) NewHighlightUseCase() usecase.HighlightUseCase {
	return usecase.NewHighlightUseCase(i.NewHighlightRepository(), i.NewStoryRepository())
}

func (i interactor) NewStoryHandler() handler.StoryHandler {
	return handler.NewStoryHandler(i.NewStoryUseCase())
}

func (i interactor) NewHighlightHandler() handler.HighlightHandler {
	return handler.NewHighlightHandler(i.NewHighlightUseCase())
}

func NewInteractor(cassandraClient *gocql.Session, redisClient *redis.Client) Interactor {
	return &interactor{cassandraClient: cassandraClient, redisClient: redisClient}
}
