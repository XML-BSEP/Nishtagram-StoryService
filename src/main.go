package main

import (
	"context"
	logger "github.com/jelena-vlajkov/logger/logger"
	"log"
	"os"
	router2 "story-service/http/router"
	"story-service/infrastructure/cassandra_config"
	"story-service/infrastructure/redis_config"
	"story-service/interactor"
)

func main() {
	logger := logger.InitializeLogger("story-service", context.Background())
	cassandraSession, err := cassandra_config.NewCassandraSession(logger)
	if err != nil {
		log.Println(err)
	}

	redisClient := redis_config.NewReddisConn(logger)

	i := interactor.NewInteractor(cassandraSession, redisClient, logger)
	appHandler := i.NewAppHandler()

	router := router2.NewRouter(appHandler, logger)

	logger.Logger.Infof("server listening on port %v\n", "8084")
	
	if os.Getenv("DOCKER_ENV") == "" {

		err := router.RunTLS("localhost:8084", "src/certificate/cert.pem", "src/certificate/key.pem")
		if err != nil {
			return 
		}
	} else {
		err := router.Run(":8084")
		if err != nil {
			return 
		}
	}


}

