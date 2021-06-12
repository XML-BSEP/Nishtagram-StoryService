module story-service

go 1.16

replace github.com/jelena-vlajkov/logger/logger => ../../Nishtagram-Logger/

require (
	github.com/casbin/casbin/v2 v2.31.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-resty/resty/v2 v2.6.0
	github.com/gocql/gocql v0.0.0-20210515062232-b7ef815b4556
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/jelena-vlajkov/logger/logger v1.0.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/microcosm-cc/bluemonday v1.0.10
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
