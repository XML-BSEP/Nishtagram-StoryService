package redis_config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func init_viper() {
	viper.SetConfigFile(`config/redis_config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
func NewReddisConn() *redis.Client {
	init_viper()
	address := viper.GetString(`server.address`)
	port := viper.GetString(`server.port`)

	return redis.NewClient(&redis.Options{
		Addr: address + ":" + port,
		Password: "",
		DB: 0,
	})
}