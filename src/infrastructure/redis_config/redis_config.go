package redis_config

import (
	"github.com/go-redis/redis/v8"
	logger "github.com/jelena-vlajkov/logger/logger"
	"github.com/spf13/viper"
)

func init_viper(logger *logger.Logger) {
	viper.SetConfigFile(`src/config/redis_config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Logger.Errorf("error while connecting to redis, error: %v\n", err)
	}

}
func NewReddisConn(logger *logger.Logger) *redis.Client {
	init_viper(logger)
	address := viper.GetString(`server.address`)
	port := viper.GetString(`server.port`)

	return redis.NewClient(&redis.Options{
		Addr: address + ":" + port,
		Password: "",
		DB: 0,
	})
}