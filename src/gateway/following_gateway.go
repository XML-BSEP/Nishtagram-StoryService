package gateway

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	logger "github.com/jelena-vlajkov/logger/logger"
	"os"
	"story-service/domain"
)

func GetAllUserFollowing(ctx context.Context, userId string, logger *logger.Logger) ([]domain.Profile, error) {
	client := resty.New()
	domain_usr := os.Getenv("FOLLOW_DOMAIN")
	if domain_usr == "" {
		domain_usr = "127.0.0.1"
	}

	userDto := domain.Profile{Id: userId}
	resp, _ := client.R().
		SetBody(userDto).
		EnableTrace().
		Post("https://" + domain_usr + ":8089/usersFollowings")

	var responseDTO FollowingResponseDTO
	err := json.Unmarshal(resp.Body(), &responseDTO)
	if err != nil {
		logger.Logger.Errorf("error while getting followings for user %v\n", userId)
		return nil, err
	}

	return responseDTO.Data, nil
}
