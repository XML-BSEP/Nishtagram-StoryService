package main

import (
	"log"
	router2 "story-service/http/router"
	"story-service/infrastructure/cassandra_config"
	"story-service/infrastructure/redis_config"
	"story-service/interactor"
)

func main() {
	cassandraSession, err := cassandra_config.NewCassandraSession()
	if err != nil {
		log.Println(err)
	}

	redisClient := redis_config.NewReddisConn()

	i := interactor.NewInteractor(cassandraSession, redisClient)
	appHandler := i.NewAppHandler()

	router := router2.NewRouter(appHandler)

	router.RunTLS("localhost:8084", "certificate/cert.pem", "certificate/key.pem")
}

