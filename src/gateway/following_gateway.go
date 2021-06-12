package gateway

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/domain"
)

func GetAllUserFollowing(ctx context.Context, userId string, logger *logger.Logger) ([]domain.Profile, error) {
	client := resty.New()
	userDto := domain.Profile{Id: userId}
	resp, _ := client.R().
		SetBody(userDto).
		EnableTrace().
		Post("https://127.0.0.1:8089/usersFollowings")

	var responseDTO FollowingResponseDTO
	err := json.Unmarshal(resp.Body(), &responseDTO)
	if err != nil {
		logger.Logger.Errorf("error while getting followings for user %v\n", userId)
		return nil, err
	}

	return responseDTO.Data, nil
}
