package interactor

import (
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/http/handler"
	"story-service/infrastructure/data_seeder"
	"story-service/repository"
	"story-service/usecase"
)

type appHandler struct {
	handler.HighlightHandler
	handler.StoryHandler
	handler.ReportHandler
}


type Interactor interface {
	NewStoryRepository() repository.StoryRepo
	NewHighlightRepository() repository.HighlightRepo
	NewReportRepository() repository.ReportRepository

	NewStoryUseCase() usecase.StoryUseCase
	NewHighlightUseCase() usecase.HighlightUseCase
	NewRedisUseCase() usecase.RedisUseCase
	NewReportUseCase() usecase.ReportUseCase

	NewAppHandler() handler.AppHandler
	NewStoryHandler() handler.StoryHandler
	NewHighlightHandler() handler.HighlightHandler
	NewReportHandler() handler.ReportHandler
}

type interactor struct {
	cassandraClient *gocql.Session
	redisClient *redis.Client
	logger *logger.Logger
}

func (i interactor) NewReportRepository() repository.ReportRepository {
	return repository.NewReportRepository(i.cassandraClient)
}

func (i interactor) NewReportUseCase() usecase.ReportUseCase {
	return usecase.NewReportUseCase(i.NewReportRepository(), i.NewStoryRepository())
}

func (i interactor) NewReportHandler() handler.ReportHandler {
	return handler.NewReportHandler(i.NewReportUseCase(), i.logger)
}

func (i interactor) NewAppHandler() handler.AppHandler {
	appHandler := appHandler{}
	appHandler.StoryHandler = i.NewStoryHandler()
	appHandler.HighlightHandler = i.NewHighlightHandler()
	appHandler.ReportHandler = i.NewReportHandler()

	data_seeder.SeedData(i.cassandraClient, i.redisClient)
	return appHandler
}

func (i interactor) NewRedisUseCase() usecase.RedisUseCase {
	return usecase.NewRedisUsecase(i.redisClient, i.logger)
}

func (i interactor) NewStoryRepository() repository.StoryRepo {
	return repository.NewStoryRepo(i.cassandraClient, i.logger)
}

func (i interactor) NewHighlightRepository() repository.HighlightRepo {
	return repository.NewHighlightRepo(i.cassandraClient, i.logger)
}

func (i interactor) NewStoryUseCase() usecase.StoryUseCase {
	return usecase.NewStoryUseCase(i.NewStoryRepository(), i.NewRedisUseCase(), i.logger)
}

func (i interactor) NewHighlightUseCase() usecase.HighlightUseCase {
	return usecase.NewHighlightUseCase(i.NewHighlightRepository(), i.NewStoryRepository(), i.logger)
}

func (i interactor) NewStoryHandler() handler.StoryHandler {
	return handler.NewStoryHandler(i.NewStoryUseCase(), i.logger)
}

func (i interactor) NewHighlightHandler() handler.HighlightHandler {
	return handler.NewHighlightHandler(i.NewHighlightUseCase(), i.logger)
}

func NewInteractor(cassandraClient *gocql.Session, redisClient *redis.Client, logger *logger.Logger) Interactor {
	return &interactor{cassandraClient: cassandraClient, redisClient: redisClient, logger: logger}
}
