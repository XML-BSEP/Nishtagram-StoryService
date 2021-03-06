package cassandra_config

import (
	"fmt"
	"github.com/gocql/gocql"
	logger "github.com/jelena-vlajkov/logger/logger"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"time"
)

func init_viper() {
	if os.Getenv("DOCKER_ENV") != "" {
		viper.SetConfigFile(`src/config/cassandra_config.json`)
	} else {
		viper.SetConfigFile(`config/cassandra_config.json`)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}
}

const (
	CreateKeyspace ="CREATE KEYSPACE if not exists story_keyspace WITH replication = { 'class': 'SimpleStrategy', 'replication_factor': '1' };"
)

func NewCassandraSession(logger *logger.Logger) (*gocql.Session, error) {
	init_viper()
	var domain string
	if os.Getenv("DOCKER_ENV") != "" {
		domain = viper.GetString(`server.docker_domain`)+ ":" + viper.GetString(`server.port_docker`)
	} else {
		domain = viper.GetString(`server.domain`)+ ":" + viper.GetString(`server.port_localhost`)
	}
	fmt.Println(domain)
	cluster := gocql.NewCluster(domain)
	cluster.ProtoVersion, _ = strconv.Atoi(viper.GetString(`proto_version`))
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = time.Second * 1000
	//cluster.Keyspace = "post_keyspace"
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}
	cluster.DisableInitialHostLookup = true

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Logger.Fatalf("failed to connect to Cassandra Story DB, error: %v\n", err)
		return nil, err
	}

	err = session.Query(CreateKeyspace).Exec()

	if err != nil {
		logger.Logger.Fatalf("cannot create keyspace in Cassandra Story DB, error: %v\n", err)
		return nil, err
	}
	return session, err
}

