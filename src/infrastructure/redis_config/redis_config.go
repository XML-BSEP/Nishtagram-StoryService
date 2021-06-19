package redis_config

import (
	"github.com/go-redis/redis/v8"
	logger "github.com/jelena-vlajkov/logger/logger"
	"github.com/spf13/viper"
)

func init_viper(logger *logger.Logger) {
	viper.SetConfigFile(`config/redis_config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Logger.Errorf("error while connecting to redis, error: %v\n", err)
	}

}
func NewReddisConn(logger *logger.Logger) *redis.Client {
	init_viper(logger)
	var domain string
	var port string
	if viper.GetBool(`docker`) {
		domain = viper.GetString(`server.docker_address`)
		port = viper.GetString(`server.port_docker`)

	} else {
		domain = viper.GetString(`server.address`)
		port = viper.GetString(`server.port_localhost`)

	}


	address := domain

	return redis.NewClient(&redis.Options{
		Addr: address + ":" + port,
		Password: "",
		DB: 0,
	})
}